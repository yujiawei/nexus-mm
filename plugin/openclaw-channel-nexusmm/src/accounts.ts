import { DEFAULT_ACCOUNT_ID } from "openclaw/plugin-sdk";
import type { OpenClawConfig } from "openclaw/plugin-sdk";
import type { NexusmmConfig } from "./config-schema.js";

export type NexusmmAccountConfig = NexusmmConfig & {
  accounts?: Record<string, NexusmmConfig | undefined>;
};

export type ResolvedNexusmmAccount = {
  accountId: string;
  name?: string;
  enabled: boolean;
  configured: boolean;
  config: {
    botToken?: string;
    apiUrl: string;
    pollIntervalMs: number;
  };
};

const DEFAULT_API_URL = "http://localhost:9876";
const DEFAULT_POLL_INTERVAL_MS = 3000;

export function listNexusmmAccountIds(cfg: OpenClawConfig): string[] {
  const channel = (cfg.channels?.nexusmm ?? {}) as NexusmmAccountConfig;
  const accountIds = Object.keys(channel.accounts ?? {});
  if (accountIds.length > 0) {
    return accountIds;
  }
  return [DEFAULT_ACCOUNT_ID];
}

export function resolveDefaultNexusmmAccountId(_cfg: OpenClawConfig): string {
  return DEFAULT_ACCOUNT_ID;
}

export function resolveNexusmmAccount(params: {
  cfg: OpenClawConfig;
  accountId?: string | null;
}): ResolvedNexusmmAccount {
  const accountId = params.accountId ?? DEFAULT_ACCOUNT_ID;
  const channel = (params.cfg.channels?.nexusmm ?? {}) as NexusmmAccountConfig;
  const accountConfig = channel.accounts?.[accountId] ?? channel;

  const botToken = accountConfig.botToken ?? channel.botToken;
  const apiUrl = accountConfig.apiUrl ?? channel.apiUrl ?? DEFAULT_API_URL;
  const pollIntervalMs =
    accountConfig.pollIntervalMs ??
    channel.pollIntervalMs ??
    DEFAULT_POLL_INTERVAL_MS;

  const enabled = accountConfig.enabled ?? channel.enabled ?? true;
  const configured = Boolean(botToken?.trim());

  return {
    accountId,
    name: accountConfig.name ?? channel.name,
    enabled,
    configured,
    config: {
      botToken,
      apiUrl,
      pollIntervalMs,
    },
  };
}
