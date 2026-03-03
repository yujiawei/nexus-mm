import {
  buildChannelConfigSchema,
  DEFAULT_ACCOUNT_ID,
  type ChannelPlugin,
} from "openclaw/plugin-sdk";
import type { OpenClawConfig } from "openclaw/plugin-sdk";
import { NexusmmConfigSchema } from "./config-schema.js";
import {
  listNexusmmAccountIds,
  resolveDefaultNexusmmAccountId,
  resolveNexusmmAccount,
  type ResolvedNexusmmAccount,
} from "./accounts.js";
import { sendMessage } from "./api.js";
import { Poller } from "./polling.js";
import { handleInboundUpdate, type NexusmmStatusSink } from "./inbound.js";

const meta = {
  id: "nexusmm",
  label: "Nexus-MM",
  selectionLabel: "Nexus-MM (Bot API)",
  docsPath: "/channels/nexusmm",
  docsLabel: "nexusmm",
  blurb: "Bot API gateway for Nexus-MM collaboration platform",
  order: 95,
};

export const nexusmmPlugin: ChannelPlugin<ResolvedNexusmmAccount> = {
  id: "nexusmm",
  meta,
  capabilities: {
    chatTypes: ["group"],
    media: false,
    reactions: false,
    threads: false,
  },
  reload: { configPrefixes: ["channels.nexusmm"] },
  configSchema: buildChannelConfigSchema(NexusmmConfigSchema),
  config: {
    listAccountIds: (cfg) => listNexusmmAccountIds(cfg),
    resolveAccount: (cfg, accountId) => resolveNexusmmAccount({ cfg, accountId }),
    defaultAccountId: (cfg) => resolveDefaultNexusmmAccountId(cfg),
    isEnabled: (account) => account.enabled,
    isConfigured: (account) => account.configured,
    describeAccount: (account) => ({
      accountId: account.accountId,
      name: account.name,
      enabled: account.enabled,
      configured: account.configured,
      apiUrl: account.config.apiUrl,
      botToken: account.config.botToken ? "[set]" : "[missing]",
    }),
  },
  messaging: {
    normalizeTarget: (target) => target.trim(),
    targetResolver: {
      looksLikeId: (input) => Boolean(input.trim()),
      hint: "<channelId>",
    },
  },
  outbound: {
    deliveryMode: "direct",
    sendText: async (ctx) => {
      const account = resolveNexusmmAccount({
        cfg: ctx.cfg as OpenClawConfig,
        accountId: ctx.accountId ?? DEFAULT_ACCOUNT_ID,
      });
      if (!account.config.botToken) {
        throw new Error("Nexus-MM botToken is not configured");
      }
      const content = ctx.text?.trim();
      if (!content) {
        return { channel: "nexusmm", to: ctx.to, messageId: "" };
      }

      await sendMessage({
        apiUrl: account.config.apiUrl,
        botToken: account.config.botToken,
        channelId: ctx.to,
        content,
      });

      return { channel: "nexusmm", to: ctx.to, messageId: "" };
    },
  },
  status: {
    defaultRuntime: {
      accountId: DEFAULT_ACCOUNT_ID,
      running: false,
      lastStartAt: null,
      lastStopAt: null,
      lastError: null,
    },
    buildAccountSnapshot: ({ account, runtime }) => ({
      accountId: account.accountId,
      name: account.name,
      enabled: account.enabled,
      configured: account.configured,
      apiUrl: account.config.apiUrl,
      running: runtime?.running ?? false,
      lastStartAt: runtime?.lastStartAt ?? null,
      lastStopAt: runtime?.lastStopAt ?? null,
      lastError: runtime?.lastError ?? null,
      lastInboundAt: runtime?.lastInboundAt ?? null,
      lastOutboundAt: runtime?.lastOutboundAt ?? null,
    }),
  },
  gateway: {
    startAccount: async (ctx) => {
      const account = ctx.account;
      if (!account.configured || !account.config.botToken) {
        throw new Error(
          `Nexus-MM not configured for account "${account.accountId}" (missing botToken)`,
        );
      }

      const log = ctx.log;
      const statusSink: NexusmmStatusSink = (patch) =>
        ctx.setStatus({ accountId: account.accountId, ...patch });

      log?.info?.(`[${account.accountId}] starting Nexus-MM polling...`);

      ctx.setStatus({
        accountId: account.accountId,
        running: true,
        lastStartAt: Date.now(),
        lastError: null,
      });

      const abortController = new AbortController();

      const poller = new Poller({
        apiUrl: account.config.apiUrl,
        botToken: account.config.botToken,
        pollIntervalMs: account.config.pollIntervalMs,
        log,
        signal: abortController.signal,
        onUpdate: (update) => {
          handleInboundUpdate({
            account,
            update,
            log,
            statusSink,
          }).catch((err) => {
            log?.error?.(`nexusmm: inbound handler failed: ${String(err)}`);
          });
        },
      });

      poller.start();

      const onAbort = () => {
        poller.stop();
        abortController.abort();
      };

      if (ctx.abortSignal.aborted) {
        onAbort();
      } else {
        ctx.abortSignal.addEventListener("abort", onAbort, { once: true });
      }

      return {
        stop: () => {
          poller.stop();
          abortController.abort();
          ctx.abortSignal.removeEventListener("abort", onAbort);
          ctx.setStatus({
            accountId: account.accountId,
            running: false,
            lastStopAt: Date.now(),
          });
        },
      };
    },
  },
};
