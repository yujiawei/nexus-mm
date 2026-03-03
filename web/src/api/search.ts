import client from './client';
import type { Message } from './types';

export async function searchMessages(query: string): Promise<Message[]> {
  const res = await client.get<Message[]>('/search', { params: { q: query } });
  return res.data;
}
