/**
 * Nexus-MM Bot API helpers.
 * Telegram-compatible HTTP endpoints: GET getUpdates, POST sendMessage.
 */

import type { GetUpdatesResponse, SendMessageResponse } from "./types.js";

export async function getUpdates(params: {
  apiUrl: string;
  botToken: string;
  offset: number;
  limit?: number;
  signal?: AbortSignal;
}): Promise<GetUpdatesResponse> {
  const base = params.apiUrl.replace(/\/+$/, "");
  const url = `${base}/bot/${params.botToken}/getUpdates?offset=${params.offset}&limit=${params.limit ?? 100}`;

  const response = await fetch(url, { signal: params.signal });

  if (!response.ok) {
    const text = await response.text().catch(() => "");
    throw new Error(`getUpdates failed (${response.status}): ${text || response.statusText}`);
  }

  return (await response.json()) as GetUpdatesResponse;
}

export async function sendMessage(params: {
  apiUrl: string;
  botToken: string;
  channelId: string;
  content: string;
  signal?: AbortSignal;
}): Promise<SendMessageResponse> {
  const base = params.apiUrl.replace(/\/+$/, "");
  const url = `${base}/bot/${params.botToken}/sendMessage`;

  const response = await fetch(url, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      channel_id: params.channelId,
      content: params.content,
    }),
    signal: params.signal,
  });

  if (!response.ok) {
    const text = await response.text().catch(() => "");
    throw new Error(`sendMessage failed (${response.status}): ${text || response.statusText}`);
  }

  return (await response.json()) as SendMessageResponse;
}
