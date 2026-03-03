import { useEffect, useState } from 'react';
import { useParams, useNavigate, Navigate } from 'react-router-dom';
import { useAuthStore } from '../../store/auth';
import { useTeamStore } from '../../store/team';
import Spinner from '../common/Spinner';

export default function InvitePage() {
  const { code } = useParams<{ code: string }>();
  const navigate = useNavigate();
  const token = useAuthStore((s) => s.token);
  const { joinByCode, selectTeam } = useTeamStore();
  const [error, setError] = useState<string | null>(null);
  const [joining, setJoining] = useState(true);

  useEffect(() => {
    if (!token || !code) return;

    const doJoin = async () => {
      try {
        const team = await joinByCode(code);
        await selectTeam(team);
        navigate('/', { replace: true });
      } catch (err: any) {
        setError(err?.response?.data?.error || 'Invalid or expired invite link');
        setJoining(false);
      }
    };
    doJoin();
  }, [token, code, joinByCode, selectTeam, navigate]);

  if (!token) {
    // Save invite code and redirect to login
    sessionStorage.setItem('pendingInvite', code || '');
    return <Navigate to="/login" replace />;
  }

  if (joining) {
    return (
      <div className="h-full flex items-center justify-center bg-gray-50">
        <div className="text-center">
          <Spinner size="lg" />
          <p className="mt-4 text-gray-500">Joining team...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="h-full flex items-center justify-center bg-gray-50">
      <div className="text-center max-w-sm mx-4">
        <div className="w-16 h-16 bg-red-100 rounded-full flex items-center justify-center mx-auto mb-4">
          <svg className="w-8 h-8 text-red-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 2.5 1.732 2.5z" />
          </svg>
        </div>
        <h2 className="text-xl font-bold text-gray-900 mb-2">Unable to Join</h2>
        <p className="text-gray-500 mb-4">{error}</p>
        <button
          onClick={() => navigate('/', { replace: true })}
          className="text-blue-600 hover:text-blue-700 text-sm font-medium"
        >
          Go to Home
        </button>
      </div>
    </div>
  );
}
