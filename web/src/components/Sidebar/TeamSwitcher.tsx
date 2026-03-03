import { useState } from 'react';
import { useTeamStore } from '../../store/team';
import Modal from '../common/Modal';
import Button from '../common/Button';

export default function TeamSwitcher() {
  const { teams, currentTeam, selectTeam, createTeam } = useTeamStore();
  const [showModal, setShowModal] = useState(false);
  const [name, setName] = useState('');
  const [displayName, setDisplayName] = useState('');
  const [creating, setCreating] = useState(false);

  const handleCreate = async () => {
    if (!name.trim() || !displayName.trim()) return;
    setCreating(true);
    try {
      const team = await createTeam(name.trim(), displayName.trim());
      await selectTeam(team);
      setShowModal(false);
      setName('');
      setDisplayName('');
    } finally {
      setCreating(false);
    }
  };

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
          onClick={() => setShowModal(true)}
          className="w-10 h-10 rounded-lg flex items-center justify-center text-slate-400 hover:text-white hover:bg-slate-700 transition-colors border-2 border-dashed border-slate-600"
          title="Create team"
        >
          <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
          </svg>
        </button>
      </div>

      <Modal open={showModal} onClose={() => setShowModal(false)} title="Create a Team">
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
            <Button variant="secondary" onClick={() => setShowModal(false)}>Cancel</Button>
            <Button onClick={handleCreate} loading={creating} disabled={!name.trim() || !displayName.trim()}>
              Create
            </Button>
          </div>
        </div>
      </Modal>
    </>
  );
}
