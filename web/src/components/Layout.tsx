import { useEffect } from 'react';
import { useAuthStore } from '../store/auth';
import { useTeamStore } from '../store/team';
import { useWebSocket } from '../hooks/useWebSocket';
import NavBar from './NavBar';
import Sidebar from './Sidebar/Sidebar';
import ChatView from './Chat/ChatView';
import Spinner from './common/Spinner';

export default function Layout() {
  const { user, loadUser, loading: authLoading } = useAuthStore();
  const { currentChannel, loadTeams, loading: teamLoading } = useTeamStore();

  useEffect(() => {
    loadUser();
  }, [loadUser]);

  useEffect(() => {
    if (user) loadTeams();
  }, [user, loadTeams]);

  useWebSocket(user?.id, currentChannel?.id);

  if (authLoading || teamLoading) {
    return (
      <div className="h-full flex items-center justify-center bg-gray-50">
        <div className="text-center">
          <Spinner size="lg" />
          <p className="mt-4 text-sm text-gray-500">Loading...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="h-full flex flex-col">
      <NavBar />
      <div className="flex-1 flex min-h-0">
        <Sidebar />
        <ChatView />
      </div>
    </div>
  );
}
