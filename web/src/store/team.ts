import { create } from 'zustand';
import type { Team, Channel, ChannelCategory } from '../api/types';
import * as teamsApi from '../api/teams';
import * as channelsApi from '../api/channels';

interface TeamState {
  teams: Team[];
  allTeams: Team[];
  currentTeam: Team | null;
  channels: Channel[];
  categories: ChannelCategory[];
  currentChannel: Channel | null;
  loading: boolean;
  loadTeams: () => Promise<void>;
  loadAllTeams: () => Promise<void>;
  joinTeam: (teamId: string) => Promise<void>;
  joinChannel: (channelId: string) => Promise<void>;
  joinByCode: (code: string) => Promise<Team>;
  selectTeam: (team: Team) => Promise<void>;
  selectChannel: (channel: Channel) => void;
  loadChannels: (teamId: string) => Promise<void>;
  loadCategories: (teamId: string) => Promise<void>;
  createTeam: (name: string, displayName: string, description?: string) => Promise<Team>;
  createChannel: (teamId: string, name: string, displayName: string, type: 'open' | 'private', purpose?: string) => Promise<Channel>;
}

export const useTeamStore = create<TeamState>((set, get) => ({
  teams: [],
  allTeams: [],
  currentTeam: null,
  channels: [],
  categories: [],
  currentChannel: null,
  loading: false,

  loadAllTeams: async () => {
    try {
      const allTeams = await teamsApi.listAllTeams();
      set({ allTeams: allTeams || [] });
    } catch {
      set({ allTeams: [] });
    }
  },

  joinTeam: async (teamId) => {
    await teamsApi.joinTeam(teamId);
    // Reload my teams.
    const teams = await teamsApi.listTeams();
    set({ teams: teams || [] });
  },

  joinChannel: async (channelId) => {
    await channelsApi.joinChannel(channelId);
    // Reload channels for current team.
    const currentTeam = get().currentTeam;
    if (currentTeam) {
      await get().loadChannels(currentTeam.id);
    }
  },

  joinByCode: async (code) => {
    const team = await teamsApi.joinByCode(code);
    const teams = await teamsApi.listTeams();
    set({ teams: teams || [] });
    return team;
  },

  loadTeams: async () => {
    set({ loading: true });
    try {
      const teams = await teamsApi.listTeams();
      set({ teams: teams || [], loading: false });
      if ((teams || []).length > 0 && !get().currentTeam) {
        await get().selectTeam(teams[0]);
      }
    } catch {
      set({ loading: false });
    }
  },

  selectTeam: async (team) => {
    set({ currentTeam: team, currentChannel: null });
    await Promise.all([get().loadChannels(team.id), get().loadCategories(team.id)]);
    const channels = get().channels;
    if (channels.length > 0) {
      set({ currentChannel: channels[0] });
    }
  },

  selectChannel: (channel) => {
    set({ currentChannel: channel });
  },

  loadChannels: async (teamId) => {
    try {
      const channels = await channelsApi.listChannels(teamId);
      set({ channels: channels || [] });
    } catch {
      set({ channels: [] });
    }
  },

  loadCategories: async (teamId) => {
    try {
      const categories = await channelsApi.listCategories(teamId);
      set({ categories: categories || [] });
    } catch {
      set({ categories: [] });
    }
  },

  createTeam: async (name, displayName, description) => {
    const team = await teamsApi.createTeam({ name, display_name: displayName, description });
    set((state) => ({ teams: [...state.teams, team] }));
    return team;
  },

  createChannel: async (teamId, name, displayName, type, purpose) => {
    const channel = await channelsApi.createChannel(teamId, {
      name,
      display_name: displayName,
      type,
      purpose,
    });
    set((state) => ({ channels: [...state.channels, channel] }));
    return channel;
  },
}));
