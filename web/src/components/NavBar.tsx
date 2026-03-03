import { useLocation, useNavigate } from 'react-router-dom';
import { useAuthStore } from '../store/auth';
import Avatar from './common/Avatar';

const navItems = [
  { path: '/', label: 'Chat' },
  { path: '/agents', label: 'Agents' },
];

export default function NavBar() {
  const location = useLocation();
  const navigate = useNavigate();
  const { user, logout } = useAuthStore();

  return (
    <div className="flex items-center justify-between h-10 px-4 bg-slate-900 text-white flex-shrink-0">
      <div className="flex items-center gap-1">
        <span className="font-bold text-sm mr-4">Nexus-MM</span>
        {navItems.map((item) => {
          const active = item.path === '/'
            ? location.pathname === '/' || (!navItems.some((n) => n.path !== '/' && location.pathname.startsWith(n.path)))
            : location.pathname.startsWith(item.path);
          return (
            <button
              key={item.path}
              onClick={() => navigate(item.path)}
              className={`px-3 py-1 text-xs font-medium rounded transition-colors ${
                active
                  ? 'bg-slate-700 text-white'
                  : 'text-slate-400 hover:text-white hover:bg-slate-800'
              }`}
            >
              {item.label}
            </button>
          );
        })}
      </div>
      <div className="flex items-center gap-3">
        {user && (
          <div className="flex items-center gap-2">
            <Avatar name={user.nickname || user.username} url={user.avatar_url} size="sm" />
            <span className="text-xs text-slate-300">{user.nickname || user.username}</span>
          </div>
        )}
        <button
          onClick={logout}
          className="text-xs text-slate-400 hover:text-white transition-colors"
        >
          Sign out
        </button>
      </div>
    </div>
  );
}
