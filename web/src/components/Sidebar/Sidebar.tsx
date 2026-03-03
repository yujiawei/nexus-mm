import { useState } from 'react';
import { useTeamStore } from '../../store/team';
import { useAuthStore } from '../../store/auth';
import TeamSwitcher from './TeamSwitcher';
import ChannelList from './ChannelList';
import TeamSettingsModal from '../Team/TeamSettingsModal';
import Avatar from '../common/Avatar';

export default function Sidebar() {
  const { currentTeam } = useTeamStore();
  const { user } = useAuthStore();
  const [showTeamSettings, setShowTeamSettings] = useState(false);

  return (
    <div className="flex h-full">
      <TeamSwitcher />
      <div className="flex flex-col w-60 bg-slate-800 text-white">
        <button
          onClick={() => setShowTeamSettings(true)}
          className="flex items-center justify-between px-4 h-14 border-b border-slate-700 flex-shrink-0 hover:bg-slate-700 transition-colors w-full text-left"
        >
          <h2 className="text-base font-bold truncate">
            {currentTeam?.display_name || 'Nexus-MM'}
          </h2>
          <svg className="w-4 h-4 text-slate-400 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
          </svg>
        </button>

        <ChannelList />

        <div className="flex items-center gap-2 px-3 py-3 border-t border-slate-700 flex-shrink-0">
          <Avatar name={user?.nickname || user?.username || '?'} url={user?.avatar_url} size="sm" />
          <div className="flex flex-col min-w-0 flex-1">
            <span className="text-sm text-slate-200 truncate leading-tight">
              {user?.nickname || user?.username}
            </span>
            <span className="text-[10px] text-slate-500 truncate leading-tight">
              @{user?.username}
            </span>
          </div>
        </div>
      </div>

      <TeamSettingsModal open={showTeamSettings} onClose={() => setShowTeamSettings(false)} />
    </div>
  );
}
