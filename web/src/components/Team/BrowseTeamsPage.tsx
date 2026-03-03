import { useEffect, useState } from 'react';
import { useTeamStore } from '../../store/team';
import Button from '../common/Button';
import Spinner from '../common/Spinner';

export default function BrowseTeamsPage() {
  const { teams, allTeams, loadAllTeams, joinTeam, selectTeam } = useTeamStore();
  const [joining, setJoining] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadAllTeams().finally(() => setLoading(false));
  }, [loadAllTeams]);

  const myTeamIds = new Set(teams.map((t) => t.id));
  const joinableTeams = allTeams.filter((t) => !myTeamIds.has(t.id));

  const handleJoin = async (teamId: string) => {
    setJoining(teamId);
    try {
      await joinTeam(teamId);
      const joined = allTeams.find((t) => t.id === teamId);
      if (joined) await selectTeam(joined);
    } finally {
      setJoining(null);
    }
  };

  if (loading) {
    return (
      <div className="flex-1 flex items-center justify-center bg-gray-50">
        <Spinner size="lg" />
      </div>
    );
  }

  return (
    <div className="flex-1 flex items-center justify-center bg-gray-50">
      <div className="w-full max-w-lg mx-4">
        <div className="text-center mb-8">
          <div className="w-16 h-16 bg-blue-100 rounded-full flex items-center justify-center mx-auto mb-4">
            <svg className="w-8 h-8 text-blue-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z" />
            </svg>
          </div>
          <h1 className="text-2xl font-bold text-gray-900 mb-2">Browse Public Teams</h1>
          <p className="text-gray-500">Join a team to start collaborating with others</p>
        </div>

        <div className="bg-white rounded-lg shadow-sm border border-gray-200 divide-y divide-gray-100">
          {joinableTeams.length === 0 && allTeams.length === 0 && (
            <div className="p-8 text-center">
              <p className="text-gray-500 mb-2">No teams available yet</p>
              <p className="text-sm text-gray-400">Ask an administrator to create one, or check back later</p>
            </div>
          )}
          {joinableTeams.length === 0 && allTeams.length > 0 && (
            <div className="p-8 text-center">
              <p className="text-gray-500">You've already joined all available teams</p>
            </div>
          )}
          {joinableTeams.map((team) => (
            <div key={team.id} className="flex items-center justify-between p-4 hover:bg-gray-50 transition-colors">
              <div className="flex items-center gap-3 min-w-0">
                <div className="w-10 h-10 rounded-lg bg-blue-600 text-white flex items-center justify-center font-bold text-sm flex-shrink-0">
                  {team.display_name.charAt(0).toUpperCase()}
                </div>
                <div className="min-w-0">
                  <p className="font-medium text-gray-900 truncate">{team.display_name}</p>
                  <p className="text-xs text-gray-500 truncate">{team.description || team.name}</p>
                </div>
              </div>
              <Button size="sm" onClick={() => handleJoin(team.id)} loading={joining === team.id}>
                Join
              </Button>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
