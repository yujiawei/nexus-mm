import type { ChannelLogSink, OpenClawConfig } from "openclaw/plugin-sdk";
import { sendMessage } from "./api.js";
import type { ResolvedNexusmmAccount } from "./accounts.js";
import type { NexusUpdate } from "./types.js";
import { getNexusmmRuntime } from "./runtime.js";

export type NexusmmStatusSink = (patch: {
  lastInboundAt?: number;
  lastOutboundAt?: number;
  lastError?: string | null;
}) => void;

export async function handleInboundUpdate(params: {
  account: ResolvedNexusmmAccount;
  update: NexusUpdate;
  log?: ChannelLogSink;
  statusSink?: NexusmmStatusSink;
}) {
  const { account, update, log, statusSink } = params;

  const rawBody = update.content?.trim();
  if (!rawBody) {
    log?.info?.(
      `nexusmm: inbound dropped update_id=${update.update_id} reason=empty-content`,
    );
    return;
  }

  const sessionId = update.channel_id || update.user_id || "";
  if (!sessionId) {
    log?.info?.("nexusmm: inbound dropped — no channel_id or user_id");
    return;
  }

  const senderId = update.user_id || "unknown";
  const fromLabel = `user:${senderId}`;

  log?.info?.(
    `nexusmm: recv message from=${senderId} channel=${sessionId}`,
  );

  const core = getNexusmmRuntime();
  const config = core.config.loadConfig() as OpenClawConfig;

  const route = core.channel.routing.resolveAgentRoute({
    cfg: config,
    channel: "nexusmm",
    accountId: account.accountId,
    peer: {
      kind: "group",
      id: sessionId,
    },
  });

  const storePath = core.channel.session.resolveStorePath(config.session?.store, {
    agentId: route.agentId,
  });

  const envelopeOptions = core.channel.reply.resolveEnvelopeFormatOptions(config);
  const previousTimestamp = core.channel.session.readSessionUpdatedAt({
    storePath,
    sessionKey: route.sessionKey,
  });

  const messageTimestamp = update.created_at
    ? new Date(update.created_at).getTime()
    : undefined;

  const body = core.channel.reply.formatAgentEnvelope({
    channel: "Nexus-MM",
    from: fromLabel,
    timestamp: messageTimestamp,
    previousTimestamp,
    envelope: envelopeOptions,
    body: rawBody,
  });

  const ctxPayload = core.channel.reply.finalizeInboundContext({
    Body: body,
    RawBody: rawBody,
    CommandBody: rawBody,
    From: `nexusmm:${senderId}`,
    To: `nexusmm:${sessionId}`,
    SessionKey: route.sessionKey,
    AccountId: route.accountId,
    ChatType: "group",
    ConversationLabel: fromLabel,
    SenderId: senderId,
    MessageSid: update.message_id || String(update.update_id),
    Timestamp: messageTimestamp,
    Provider: "nexusmm",
    Surface: "nexusmm",
    OriginatingChannel: "nexusmm",
    OriginatingTo: `nexusmm:${sessionId}`,
  });

  await core.channel.session.recordInboundSession({
    storePath,
    sessionKey: ctxPayload.SessionKey ?? route.sessionKey,
    ctx: ctxPayload,
    onRecordError: (err) => {
      log?.error?.(`nexusmm: failed updating session meta: ${String(err)}`);
    },
  });

  statusSink?.({ lastInboundAt: Date.now(), lastError: null });

  await core.channel.reply.dispatchReplyWithBufferedBlockDispatcher({
    ctx: ctxPayload,
    cfg: config,
    dispatcherOptions: {
      deliver: async (payload: {
        text?: string;
        mediaUrls?: string[];
        mediaUrl?: string;
      }) => {
        const contentParts: string[] = [];
        if (payload.text) contentParts.push(payload.text);
        const mediaUrls = [
          ...(payload.mediaUrls ?? []),
          ...(payload.mediaUrl ? [payload.mediaUrl] : []),
        ].filter(Boolean);
        if (mediaUrls.length > 0) contentParts.push(...mediaUrls);
        const content = contentParts.join("\n").trim();
        if (!content) return;

        await sendMessage({
          apiUrl: account.config.apiUrl,
          botToken: account.config.botToken ?? "",
          channelId: sessionId,
          content,
        });

        statusSink?.({ lastOutboundAt: Date.now(), lastError: null });
      },
      onError: (err, info) => {
        log?.error?.(`nexusmm ${info.kind} reply failed: ${String(err)}`);
      },
    },
  });
}
