import client from './client';
import type { Team, CreateTeamRequest } from './types';

export async function listTeams(): Promise<Team[]> {
  const res = await client.get<Team[]>('/teams');
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
