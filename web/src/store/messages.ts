import { create } from 'zustand';
import type { Message, ReactionGroup } from '../api/types';
import * as messagesApi from '../api/messages';
import * as searchApi from '../api/search';

interface MessagesState {
  messages: Message[];
  threadMessages: Message[];
  threadRoot: Message | null;
  searchResults: Message[];
  searchQuery: string;
  searching: boolean;
  loading: boolean;
  hasMore: boolean;
  loadMessages: (channelId: string) => Promise<void>;
  loadMore: (channelId: string) => Promise<void>;
  sendMessage: (channelId: string, content: string, rootId?: string) => Promise<void>;
  addIncomingMessage: (msg: Message) => void;
  openThread: (channelId: string, message: Message) => Promise<void>;
  closeThread: () => void;
  toggleReaction: (channelId: string, messageId: string, emoji: string, userId: string) => Promise<void>;
  search: (query: string) => Promise<void>;
  clearSearch: () => void;
}

export function groupReactions(reactions: { emoji_name: string; user_id: string }[], currentUserId: string): ReactionGroup[] {
  const map = new Map<string, { count: number; user_ids: string[] }>();
  for (const r of reactions) {
    const existing = map.get(r.emoji_name);
    if (existing) {
      existing.count++;
      existing.user_ids.push(r.user_id);
    } else {
      map.set(r.emoji_name, { count: 1, user_ids: [r.user_id] });
    }
  }
  return Array.from(map.entries()).map(([emoji_name, data]) => ({
    emoji_name,
    count: data.count,
    user_ids: data.user_ids,
    includes_me: data.user_ids.includes(currentUserId),
  }));
}

export const useMessagesStore = create<MessagesState>((set, get) => ({
  messages: [],
  threadMessages: [],
  threadRoot: null,
  searchResults: [],
  searchQuery: '',
  searching: false,
  loading: false,
  hasMore: true,

  loadMessages: async (channelId) => {
    set({ loading: true, hasMore: true });
    try {
      const messages = await messagesApi.listMessages(channelId, { limit: 50 });
      set({ messages: (messages || []).reverse(), loading: false, hasMore: (messages || []).length >= 50 });
    } catch {
      set({ loading: false, messages: [] });
    }
  },

  loadMore: async (channelId) => {
    const { messages, hasMore } = get();
    if (!hasMore || messages.length === 0) return;
    const oldest = messages[0];
    try {
      const older = await messagesApi.listMessages(channelId, { before: oldest.id, limit: 50 });
      const reversed = (older || []).reverse();
      set({
        messages: [...reversed, ...messages],
        hasMore: (older || []).length >= 50,
      });
    } catch {
      /* ignore */
    }
  },

  sendMessage: async (channelId, content, rootId) => {
    const msg = await messagesApi.sendMessage(channelId, { content, root_id: rootId });
    if (rootId) {
      set((state) => ({ threadMessages: [...state.threadMessages, msg] }));
      set((state) => ({
        messages: state.messages.map((m) =>
          m.id === rootId ? { ...m, reply_count: m.reply_count + 1 } : m
        ),
      }));
    } else {
      set((state) => ({ messages: [...state.messages, msg] }));
    }
  },

  addIncomingMessage: (msg) => {
    set((state) => {
      if (state.messages.some((m) => m.id === msg.id)) return state;
      if (msg.root_id) {
        const newThread = state.threadRoot && state.threadRoot.id === msg.root_id
          ? [...state.threadMessages, msg]
          : state.threadMessages;
        return {
          threadMessages: newThread,
          messages: state.messages.map((m) =>
            m.id === msg.root_id ? { ...m, reply_count: m.reply_count + 1 } : m
          ),
        };
      }
      return { messages: [...state.messages, msg] };
    });
  },

  openThread: async (channelId, message) => {
    set({ threadRoot: message, threadMessages: [] });
    try {
      const thread = await messagesApi.getThread(channelId, message.id);
      set({ threadMessages: thread || [] });
    } catch {
      /* ignore */
    }
  },

  closeThread: () => {
    set({ threadRoot: null, threadMessages: [] });
  },

  toggleReaction: async (channelId, messageId, emoji, userId) => {
    const { messages } = get();
    const msg = messages.find((m) => m.id === messageId);
    const existing = msg?.reactions?.find((r) => r.emoji_name === emoji && r.user_id === userId);

    if (existing) {
      await messagesApi.removeReaction(channelId, messageId, emoji);
      set((state) => ({
        messages: state.messages.map((m) =>
          m.id === messageId
            ? { ...m, reactions: (m.reactions || []).filter((r) => !(r.emoji_name === emoji && r.user_id === userId)) }
            : m
        ),
      }));
    } else {
      await messagesApi.addReaction(channelId, messageId, emoji);
      set((state) => ({
        messages: state.messages.map((m) =>
          m.id === messageId
            ? {
                ...m,
                reactions: [
                  ...(m.reactions || []),
                  { id: '', message_id: messageId, user_id: userId, emoji_name: emoji, created_at: new Date().toISOString() },
                ],
              }
            : m
        ),
      }));
    }
  },

  search: async (query) => {
    if (!query.trim()) return;
    set({ searching: true, searchQuery: query });
    try {
      const results = await searchApi.searchMessages(query);
      set({ searchResults: results || [], searching: false });
    } catch {
      set({ searchResults: [], searching: false });
    }
  },

  clearSearch: () => {
    set({ searchResults: [], searchQuery: '', searching: false });
  },
}));
