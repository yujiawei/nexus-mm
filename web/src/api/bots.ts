import client from './client';
import type { BotInfo, CreateBotRequest } from './types';

export async function createBot(req: CreateBotRequest): Promise<BotInfo> {
  const { data } = await client.post('/bots', req);
  return data;
}

export async function listBots(): Promise<BotInfo[]> {
  const { data } = await client.get('/bots');
  return data;
}

export async function regenerateToken(botId: string): Promise<{ bot_token: string }> {
  const { data } = await client.post(`/bots/${botId}/regenerate-token`);
  return data;
}

export async function updateWebhook(botId: string, url: string): Promise<void> {
  await client.put(`/bots/${botId}/webhook`, { url });
}
