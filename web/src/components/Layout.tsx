import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuthStore } from '../store/auth';
import { useTeamStore } from '../store/team';
import { useWebSocket } from '../hooks/useWebSocket';
import NavBar from './NavBar';
import Sidebar from './Sidebar/Sidebar';
import ChatView from './Chat/ChatView';
import BrowseTeamsPage from './Team/BrowseTeamsPage';
import Spinner from './common/Spinner';

export default function Layout() {
  const { user, loadUser, loading: authLoading } = useAuthStore();
  const { teams, loadTeams, loading: teamLoading } = useTeamStore();
  const navigate = useNavigate();

  useEffect(() => {
    loadUser();
  }, [loadUser]);

  useEffect(() => {
    if (user) {
      loadTeams();
      // Handle pending invite after login
      const pendingInvite = sessionStorage.getItem('pendingInvite');
      if (pendingInvite) {
        sessionStorage.removeItem('pendingInvite');
        navigate(`/invite/${pendingInvite}`, { replace: true });
      }
    }
  }, [user, loadTeams, navigate]);

  useWebSocket();

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

  // Show browse teams page if user has no teams
  if (teams.length === 0) {
    return (
      <div className="h-full flex flex-col">
        <NavBar />
        <BrowseTeamsPage />
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
