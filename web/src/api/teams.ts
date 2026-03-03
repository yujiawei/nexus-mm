import client from './client';
import type { Team, TeamMember, CreateTeamRequest } from './types';

export async function listTeams(): Promise<Team[]> {
  const res = await client.get<Team[]>('/teams');
  return res.data;
}

export async function listAllTeams(): Promise<Team[]> {
  const res = await client.get<Team[]>('/teams/all');
  return res.data;
}

export async function getTeam(id: string): Promise<Team> {
  const res = await client.get<Team>(`/teams/${id}`);
  return res.data;
}

export async function createTeam(data: CreateTeamRequest): Promise<Team> {
  const res = await client.post<Team>('/teams', data);
  return res.data;
}

export async function joinTeam(teamId: string): Promise<void> {
  await client.post(`/teams/${teamId}/join`);
}

export async function listTeamMembers(teamId: string): Promise<TeamMember[]> {
  const res = await client.get<TeamMember[]>(`/teams/${teamId}/members`);
  return res.data;
}

export async function removeTeamMember(teamId: string, userId: string): Promise<void> {
  await client.delete(`/teams/${teamId}/members/${userId}`);
}

export async function createInviteLink(teamId: string, maxUses = 0, expireDays = 0): Promise<{ code: string }> {
  const res = await client.post(`/teams/${teamId}/invite-link`, {
    max_uses: maxUses,
    expire_days: expireDays,
  });
  return res.data;
}

export async function joinByCode(code: string): Promise<Team> {
  const res = await client.post<Team>(`/teams/join-by-code/${code}`);
  return res.data;
}
