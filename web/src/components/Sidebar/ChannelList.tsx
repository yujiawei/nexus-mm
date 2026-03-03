import { useState } from 'react';
import { useTeamStore } from '../../store/team';
import type { Channel } from '../../api/types';
import Modal from '../common/Modal';
import Button from '../common/Button';

function ChannelIcon({ type }: { type: string }) {
  if (type === 'private') {
    return (
      <svg className="w-4 h-4 flex-shrink-0 opacity-70" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
      </svg>
    );
  }
  return (
    <svg className="w-4 h-4 flex-shrink-0 opacity-70" fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M7 20l4-16m2 16l4-16M6 9h14M4 15h14" />
    </svg>
  );
}

function CategorySection({
  label,
  channels,
  currentChannel,
  onSelect,
}: {
  label: string;
  channels: Channel[];
  currentChannel: Channel | null;
  onSelect: (ch: Channel) => void;
}) {
  const [collapsed, setCollapsed] = useState(false);

  return (
    <div className="mb-1">
      <button
        onClick={() => setCollapsed(!collapsed)}
        className="flex items-center gap-1 px-3 py-1.5 w-full text-xs font-semibold uppercase tracking-wider text-slate-400 hover:text-slate-200 transition-colors"
      >
        <svg
          className={`w-3 h-3 transition-transform ${collapsed ? '-rotate-90' : ''}`}
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
        >
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
        </svg>
        {label}
      </button>
      {!collapsed &&
        channels.map((ch) => (
          <button
            key={ch.id}
            onClick={() => onSelect(ch)}
            className={`flex items-center gap-2 w-full px-4 py-1.5 text-sm transition-colors ${
              currentChannel?.id === ch.id
                ? 'bg-slate-700/80 text-white'
                : 'text-slate-300 hover:bg-slate-700/40 hover:text-white'
            }`}
          >
            <ChannelIcon type={ch.type} />
            <span className="truncate">{ch.display_name}</span>
          </button>
        ))}
    </div>
  );
}

export default function ChannelList() {
  const { channels, currentChannel, currentTeam, selectChannel, createChannel } = useTeamStore();
  const [showModal, setShowModal] = useState(false);
  const [name, setName] = useState('');
  const [displayName, setDisplayName] = useState('');
  const [chType, setChType] = useState<'open' | 'private'>('open');
  const [creating, setCreating] = useState(false);

  const openChannels = channels.filter((ch) => ch.type === 'open');
  const privateChannels = channels.filter((ch) => ch.type === 'private');
  const directChannels = channels.filter((ch) => ch.type === 'direct');

  const handleCreate = async () => {
    if (!name.trim() || !displayName.trim() || !currentTeam) return;
    setCreating(true);
    try {
      const ch = await createChannel(currentTeam.id, name.trim(), displayName.trim(), chType);
      selectChannel(ch);
      setShowModal(false);
      setName('');
      setDisplayName('');
    } finally {
      setCreating(false);
    }
  };

  return (
    <>
      <div className="flex-1 overflow-y-auto sidebar-scrollbar py-2">
        {openChannels.length > 0 && (
          <CategorySection
            label="Channels"
            channels={openChannels}
            currentChannel={currentChannel}
            onSelect={selectChannel}
          />
        )}
        {privateChannels.length > 0 && (
          <CategorySection
            label="Private"
            channels={privateChannels}
            currentChannel={currentChannel}
            onSelect={selectChannel}
          />
        )}
        {directChannels.length > 0 && (
          <CategorySection
            label="Direct Messages"
            channels={directChannels}
            currentChannel={currentChannel}
            onSelect={selectChannel}
          />
        )}
        {channels.length === 0 && (
          <p className="px-4 py-8 text-sm text-slate-500 text-center">No channels yet</p>
        )}
      </div>

      <div className="px-3 py-2 border-t border-slate-700">
        <button
          onClick={() => setShowModal(true)}
          className="flex items-center gap-2 w-full px-2 py-1.5 text-sm text-slate-400 hover:text-white hover:bg-slate-700/40 rounded-md transition-colors"
        >
          <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
          </svg>
          Add Channel
        </button>
      </div>

      <Modal open={showModal} onClose={() => setShowModal(false)} title="Create a Channel">
        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Channel Name</label>
            <input
              type="text"
              value={name}
              onChange={(e) => setName(e.target.value.toLowerCase().replace(/[^a-z0-9-]/g, ''))}
              className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="new-channel"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Display Name</label>
            <input
              type="text"
              value={displayName}
              onChange={(e) => setDisplayName(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="New Channel"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Type</label>
            <div className="flex gap-4">
              <label className="flex items-center gap-2 text-sm text-gray-600 cursor-pointer">
                <input
                  type="radio"
                  checked={chType === 'open'}
                  onChange={() => setChType('open')}
                  className="text-blue-600"
                />
                Public
              </label>
              <label className="flex items-center gap-2 text-sm text-gray-600 cursor-pointer">
                <input
                  type="radio"
                  checked={chType === 'private'}
                  onChange={() => setChType('private')}
                  className="text-blue-600"
                />
                Private
              </label>
            </div>
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
