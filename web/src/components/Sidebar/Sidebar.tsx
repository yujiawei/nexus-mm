import { useTeamStore } from '../../store/team';
import { useAuthStore } from '../../store/auth';
import TeamSwitcher from './TeamSwitcher';
import ChannelList from './ChannelList';
import Avatar from '../common/Avatar';

export default function Sidebar() {
  const { currentTeam } = useTeamStore();
  const { user, logout } = useAuthStore();

  return (
    <div className="flex h-full">
      <TeamSwitcher />
      <div className="flex flex-col w-60 bg-slate-800 text-white">
        <div className="flex items-center px-4 h-14 border-b border-slate-700 flex-shrink-0">
          <h2 className="text-base font-bold truncate">
            {currentTeam?.display_name || 'Nexus-MM'}
          </h2>
        </div>

        <ChannelList />

        <div className="flex items-center gap-2 px-3 py-3 border-t border-slate-700 flex-shrink-0">
          <Avatar name={user?.nickname || user?.username || '?'} url={user?.avatar_url} size="sm" />
          <span className="text-sm text-slate-300 truncate flex-1">
            {user?.nickname || user?.username}
          </span>
          <button
            onClick={logout}
            title="Sign out"
            className="text-slate-400 hover:text-white transition-colors p-1"
          >
            <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
            </svg>
          </button>
        </div>
      </div>
    </div>
  );
}
