export interface User {
  id: string;
  username: string;
  email: string;
  nickname: string;
  avatar_url?: string;
  role: string;
  created_at: string;
  updated_at: string;
}

export interface LoginRequest {
  username: string;
  password: string;
}

export interface RegisterRequest {
  username: string;
  email: string;
  password: string;
  nickname: string;
}

export interface LoginResponse {
  token: string;
  user: User;
}

export interface Team {
  id: string;
  name: string;
  display_name: string;
  description?: string;
  creator_id: string;
  retention_days: number;
  created_at: string;
  updated_at: string;
}

export interface CreateTeamRequest {
  name: string;
  display_name: string;
  description?: string;
}

export interface Channel {
  id: string;
  team_id: string;
  name: string;
  display_name: string;
  type: 'open' | 'private' | 'direct';
  purpose?: string;
  creator_id: string;
  retention_days: number;
  created_at: string;
  updated_at: string;
}

export interface CreateChannelRequest {
  name: string;
  display_name: string;
  type: 'open' | 'private';
  purpose?: string;
}

export interface Message {
  id: string;
  channel_id: string;
  user_id: string;
  content: string;
  type: string;
  root_id?: string;
  reply_count: number;
  created_at: string;
  updated_at: string;
  user?: User;
  reactions?: Reaction[];
  is_pinned?: boolean;
}

export interface SendMessageRequest {
  content: string;
  root_id?: string;
}

export interface Reaction {
  id: string;
  message_id: string;
  user_id: string;
  emoji_name: string;
  created_at: string;
}

export interface PinnedMessage {
  id: string;
  channel_id: string;
  message_id: string;
  user_id: string;
  created_at: string;
}

export interface ChannelCategory {
  id: string;
  user_id: string;
  team_id: string;
  display_name: string;
  sort_order: number;
  created_at: string;
  updated_at: string;
  channels?: Channel[];
}

export interface CreateCategoryRequest {
  display_name: string;
  sort_order?: number;
}

export interface SlashCommand {
  id: string;
  team_id: string;
  trigger: string;
  url: string;
  method: string;
  creator_id: string;
  created_at: string;
  updated_at: string;
}

export interface TeamMember {
  team_id: string;
  user_id: string;
  role: string;
  created_at: string;
}

export interface ChannelMember {
  channel_id: string;
  user_id: string;
  role: string;
  created_at: string;
}

export interface SearchResult {
  messages: Message[];
}

export interface ReactionGroup {
  emoji_name: string;
  count: number;
  user_ids: string[];
  includes_me: boolean;
}

export interface BotInfo {
  id: string;
  username: string;
  nickname: string;
  description: string;
  webhook_url: string;
  token: string;
  created_at: string;
}

export interface CreateBotRequest {
  name: string;
  description: string;
}
