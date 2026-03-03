import { useNavigate } from 'react-router-dom';
import { useTeamStore } from '../../store/team';
import { useAuthStore } from '../../store/auth';
import TeamSwitcher from './TeamSwitcher';
import ChannelList from './ChannelList';
import Avatar from '../common/Avatar';

export default function Sidebar() {
  const navigate = useNavigate();
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
            onClick={() => navigate('/agents')}
            title="Agents"
            className="text-slate-400 hover:text-white transition-colors p-1"
          >
            <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9.75 3.104v5.714a2.25 2.25 0 01-.659 1.591L5 14.5M9.75 3.104c-.251.023-.501.05-.75.082m.75-.082a24.301 24.301 0 014.5 0m0 0v5.714a2.25 2.25 0 00.659 1.591L19 14.5M14.25 3.104c.251.023.501.05.75.082M19 14.5l-1.5 4.5H6.5L5 14.5m14 0H5" />
            </svg>
          </button>
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
