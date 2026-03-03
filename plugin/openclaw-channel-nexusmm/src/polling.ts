/**
 * HTTP polling loop for Nexus-MM getUpdates.
 */

import type { ChannelLogSink } from "openclaw/plugin-sdk";
import { getUpdates } from "./api.js";
import type { NexusUpdate } from "./types.js";

const MAX_RETRIES = 5;
const RETRY_BACKOFF_MS = 5000;

export interface PollerOptions {
  apiUrl: string;
  botToken: string;
  pollIntervalMs: number;
  onUpdate: (update: NexusUpdate) => void;
  log?: ChannelLogSink;
  signal?: AbortSignal;
}

export class Poller {
  private polling = false;
  private timer: ReturnType<typeof setTimeout> | null = null;
  private offset = 0;
  private retryCount = 0;

  constructor(private opts: PollerOptions) {}

  start(): void {
    if (this.polling) return;
    this.polling = true;
    this.offset = 0;
    this.retryCount = 0;
    this.opts.log?.info?.(
      `nexusmm: starting poll loop -> ${this.opts.apiUrl}/bot/***`,
    );
    this.poll();
  }

  stop(): void {
    this.polling = false;
    if (this.timer) {
      clearTimeout(this.timer);
      this.timer = null;
    }
  }

  private async poll(): Promise<void> {
    if (!this.polling) return;

    try {
      const data = await getUpdates({
        apiUrl: this.opts.apiUrl,
        botToken: this.opts.botToken,
        offset: this.offset,
        signal: this.opts.signal,
      });

      if (data.ok && Array.isArray(data.result)) {
        for (const update of data.result) {
          if (update.update_id >= this.offset) {
            this.offset = update.update_id + 1;
          }
          this.opts.onUpdate(update);
        }
      }

      this.retryCount = 0;
    } catch (err) {
      if (this.opts.signal?.aborted) return;

      this.retryCount++;
      const msg = err instanceof Error ? err.message : String(err);
      this.opts.log?.error?.(
        `nexusmm: poll error (attempt ${this.retryCount}): ${msg}`,
      );

      if (this.retryCount >= MAX_RETRIES) {
        this.opts.log?.error?.(
          "nexusmm: max retries reached, backing off...",
        );
        this.retryCount = 0;
        await new Promise((r) => setTimeout(r, RETRY_BACKOFF_MS));
      }
    }

    if (this.polling) {
      this.timer = setTimeout(() => this.poll(), this.opts.pollIntervalMs);
    }
  }
}
