/** Nexus-MM Bot API types (Telegram-compatible) */

export interface NexusUpdate {
  update_id: number;
  message_id?: string;
  channel_id?: string;
  channel_type?: number;
  content?: string;
  user_id?: string;
  created_at?: string;
}

export interface GetUpdatesResponse {
  ok: boolean;
  result?: NexusUpdate[];
  error?: string;
}

export interface SendMessageResponse {
  ok: boolean;
  error?: string;
  message_id?: string;
}
