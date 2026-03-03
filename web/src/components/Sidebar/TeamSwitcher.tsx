import { useState, useEffect } from 'react';
import { useTeamStore } from '../../store/team';
import Modal from '../common/Modal';
import Button from '../common/Button';

export default function TeamSwitcher() {
  const { teams, allTeams, currentTeam, selectTeam, createTeam, loadAllTeams, joinTeam } = useTeamStore();
  const [showCreate, setShowCreate] = useState(false);
  const [showBrowse, setShowBrowse] = useState(false);
  const [name, setName] = useState('');
  const [displayName, setDisplayName] = useState('');
  const [creating, setCreating] = useState(false);
  const [joining, setJoining] = useState<string | null>(null);

  useEffect(() => {
    if (showBrowse) loadAllTeams();
  }, [showBrowse, loadAllTeams]);

  const handleCreate = async () => {
    if (!name.trim() || !displayName.trim()) return;
    setCreating(true);
    try {
      const team = await createTeam(name.trim(), displayName.trim());
      await selectTeam(team);
      setShowCreate(false);
      setName('');
      setDisplayName('');
    } finally {
      setCreating(false);
    }
  };

  const handleJoin = async (teamId: string) => {
    setJoining(teamId);
    try {
      await joinTeam(teamId);
      const joined = allTeams.find((t) => t.id === teamId);
      if (joined) await selectTeam(joined);
      setShowBrowse(false);
    } finally {
      setJoining(null);
    }
  };

  const myTeamIds = new Set(teams.map((t) => t.id));
  const joinableTeams = allTeams.filter((t) => !myTeamIds.has(t.id));

  return (
    <>
      <div className="flex flex-col items-center py-3 gap-2 bg-slate-900 w-16 flex-shrink-0">
        {teams.map((team) => (
          <button
            key={team.id}
            onClick={() => selectTeam(team)}
            title={team.display_name}
            className={`w-10 h-10 rounded-lg flex items-center justify-center text-sm font-bold transition-all ${
              currentTeam?.id === team.id
                ? 'bg-blue-600 text-white'
                : 'bg-slate-700 text-slate-300 hover:bg-slate-600'
            }`}
          >
            {team.display_name.charAt(0).toUpperCase()}
          </button>
        ))}
        <button
          onClick={() => setShowBrowse(true)}
          className="w-10 h-10 rounded-lg flex items-center justify-center text-slate-400 hover:text-white hover:bg-slate-700 transition-colors"
          title="Browse teams"
        >
          <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
        </button>
        <button
          onClick={() => setShowCreate(true)}
          className="w-10 h-10 rounded-lg flex items-center justify-center text-slate-400 hover:text-white hover:bg-slate-700 transition-colors border-2 border-dashed border-slate-600"
          title="Create team"
        >
          <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
          </svg>
        </button>
      </div>

      <Modal open={showBrowse} onClose={() => setShowBrowse(false)} title="Browse Teams">
        <div className="space-y-2 max-h-80 overflow-y-auto">
          {joinableTeams.length === 0 && (
            <p className="text-sm text-gray-500 text-center py-4">No teams available to join</p>
          )}
          {joinableTeams.map((team) => (
            <div key={team.id} className="flex items-center justify-between p-3 rounded-lg bg-gray-50 hover:bg-gray-100">
              <div>
                <p className="font-medium text-sm">{team.display_name}</p>
                <p className="text-xs text-gray-500">{team.name}</p>
              </div>
              <Button
                size="sm"
                onClick={() => handleJoin(team.id)}
                loading={joining === team.id}
              >
                Join
              </Button>
            </div>
          ))}
        </div>
      </Modal>

      <Modal open={showCreate} onClose={() => setShowCreate(false)} title="Create a Team">
        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Team Name</label>
            <input
              type="text"
              value={name}
              onChange={(e) => setName(e.target.value.toLowerCase().replace(/[^a-z0-9-]/g, ''))}
              className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="my-team"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Display Name</label>
            <input
              type="text"
              value={displayName}
              onChange={(e) => setDisplayName(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="My Team"
            />
          </div>
          <div className="flex justify-end gap-2 pt-2">
            <Button variant="secondary" onClick={() => setShowCreate(false)}>Cancel</Button>
            <Button onClick={handleCreate} loading={creating} disabled={!name.trim() || !displayName.trim()}>
              Create
            </Button>
          </div>
        </div>
      </Modal>
    </>
  );
}
