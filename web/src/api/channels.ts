import client from './client';
import type {
  Channel,
  ChannelMember,
  CreateChannelRequest,
  ChannelCategory,
  CreateCategoryRequest,
  PinnedMessage,
} from './types';

export async function listChannels(teamId: string): Promise<Channel[]> {
  const res = await client.get<Channel[]>(`/teams/${teamId}/channels`);
  return res.data;
}

export async function getChannel(id: string): Promise<Channel> {
  const res = await client.get<Channel>(`/channels/${id}`);
  return res.data;
}

export async function createChannel(teamId: string, data: CreateChannelRequest): Promise<Channel> {
  const res = await client.post<Channel>(`/teams/${teamId}/channels`, data);
  return res.data;
}

export async function listCategories(teamId: string): Promise<ChannelCategory[]> {
  const res = await client.get<ChannelCategory[]>(`/teams/${teamId}/categories`);
  return res.data;
}

export async function createCategory(teamId: string, data: CreateCategoryRequest): Promise<ChannelCategory> {
  const res = await client.post<ChannelCategory>(`/teams/${teamId}/categories`, data);
  return res.data;
}

export async function getCategoryChannels(categoryId: string): Promise<Channel[]> {
  const res = await client.get<Channel[]>(`/categories/${categoryId}/channels`);
  return res.data;
}

export async function addChannelToCategory(categoryId: string, channelId: string): Promise<void> {
  await client.post(`/categories/${categoryId}/channels`, { channel_id: channelId });
}

export async function getPinnedMessages(channelId: string): Promise<PinnedMessage[]> {
  const res = await client.get<PinnedMessage[]>(`/channels/${channelId}/pinned`);
  return res.data;
}

export async function pinMessage(channelId: string, messageId: string): Promise<void> {
  await client.post(`/channels/${channelId}/messages/${messageId}/pin`);
}

export async function unpinMessage(channelId: string, messageId: string): Promise<void> {
  await client.delete(`/channels/${channelId}/messages/${messageId}/pin`);
}

export async function joinChannel(channelId: string): Promise<void> {
  await client.post(`/channels/${channelId}/join`);
}

export async function listChannelMembers(channelId: string): Promise<ChannelMember[]> {
  const res = await client.get<ChannelMember[]>(`/channels/${channelId}/members`);
  return res.data;
}

export async function removeChannelMember(channelId: string, userId: string): Promise<void> {
  await client.delete(`/channels/${channelId}/members/${userId}`);
}
