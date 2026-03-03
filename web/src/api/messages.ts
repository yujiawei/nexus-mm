import client from './client';
import type { Message, SendMessageRequest } from './types';

export async function listMessages(
  channelId: string,
  params?: { before?: string; limit?: number }
): Promise<Message[]> {
  const res = await client.get<Message[]>(`/channels/${channelId}/messages`, { params });
  return res.data;
}

export async function sendMessage(channelId: string, data: SendMessageRequest): Promise<Message> {
  const res = await client.post<Message>(`/channels/${channelId}/messages`, data);
  return res.data;
}

export async function getThread(channelId: string, messageId: string): Promise<Message[]> {
  const res = await client.get<Message[]>(`/channels/${channelId}/messages/${messageId}/thread`);
  return res.data;
}

export async function addReaction(channelId: string, messageId: string, emojiName: string): Promise<void> {
  await client.post(`/channels/${channelId}/messages/${messageId}/reactions`, { emoji_name: emojiName });
}

export async function removeReaction(channelId: string, messageId: string, emojiName: string): Promise<void> {
  await client.delete(`/channels/${channelId}/messages/${messageId}/reactions/${emojiName}`);
}
