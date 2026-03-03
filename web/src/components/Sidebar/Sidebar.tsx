import { useTeamStore } from '../../store/team';
import { useAuthStore } from '../../store/auth';
import TeamSwitcher from './TeamSwitcher';
import ChannelList from './ChannelList';
import Avatar from '../common/Avatar';

export default function Sidebar() {
  const { currentTeam } = useTeamStore();
  const { user } = useAuthStore();

  return (
    <div className="flex h-full">
      <TeamSwitcher />
      <div className="flex flex-col w-60 bg-slate-800 text-white">
        <div className="flex items-center justify-between px-4 h-14 border-b border-slate-700 flex-shrink-0">
          <h2 className="text-base font-bold truncate">
            {currentTeam?.display_name || 'Nexus-MM'}
          </h2>
        </div>

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
    </div>
  );
}
