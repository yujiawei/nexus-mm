import { useEffect, useState, useCallback } from 'react';
import { useTeamStore } from '../../store/team';
import { useAuthStore } from '../../store/auth';
import { listTeamMembers, removeTeamMember, createInviteLink } from '../../api/teams';
import type { TeamMember } from '../../api/types';
import Modal from '../common/Modal';
import Button from '../common/Button';
import Avatar from '../common/Avatar';
import Spinner from '../common/Spinner';

interface Props {
  open: boolean;
  onClose: () => void;
}

export default function TeamSettingsModal({ open, onClose }: Props) {
  const { currentTeam } = useTeamStore();
  const { user } = useAuthStore();
  const [members, setMembers] = useState<TeamMember[]>([]);
  const [loading, setLoading] = useState(false);
  const [inviteCode, setInviteCode] = useState<string | null>(null);
  const [generatingLink, setGeneratingLink] = useState(false);
  const [copied, setCopied] = useState(false);
  const [myRole, setMyRole] = useState<string>('member');
  const [removingUser, setRemovingUser] = useState<string | null>(null);

  const loadMembers = useCallback(async () => {
    if (!currentTeam) return;
    setLoading(true);
    try {
      const data = await listTeamMembers(currentTeam.id);
      setMembers(data || []);
      const me = data?.find((m) => m.user_id === user?.id);
      if (me) setMyRole(me.role);
    } finally {
      setLoading(false);
    }
  }, [currentTeam, user]);

  useEffect(() => {
    if (open) {
      loadMembers();
      setInviteCode(null);
      setCopied(false);
    }
  }, [open, loadMembers]);

  const handleGenerateLink = async () => {
    if (!currentTeam) return;
    setGeneratingLink(true);
    try {
      const result = await createInviteLink(currentTeam.id);
      setInviteCode(result.code);
      setCopied(false);
    } finally {
      setGeneratingLink(false);
    }
  };

  const handleCopy = () => {
    if (!inviteCode) return;
    const link = `${window.location.origin}/invite/${inviteCode}`;
    navigator.clipboard.writeText(link);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  const handleRemoveMember = async (userId: string) => {
    if (!currentTeam) return;
    setRemovingUser(userId);
    try {
      await removeTeamMember(currentTeam.id, userId);
      await loadMembers();
    } finally {
      setRemovingUser(null);
    }
  };

  const isAdmin = myRole === 'owner' || myRole === 'admin';
  const inviteLink = inviteCode ? `${window.location.origin}/invite/${inviteCode}` : '';

  return (
    <Modal open={open} onClose={onClose} title={currentTeam?.display_name || 'Team Settings'}>
      <div className="space-y-5">
        {/* Invite Link Section */}
        <div>
          <h4 className="text-sm font-semibold text-gray-700 mb-2">Invite Link</h4>
          {inviteCode ? (
            <div className="flex items-center gap-2">
              <input
                type="text"
                readOnly
                value={inviteLink}
                className="flex-1 px-3 py-2 border border-gray-300 rounded-md text-sm bg-gray-50 text-gray-700"
              />
              <Button size="sm" onClick={handleCopy} variant={copied ? 'secondary' : 'primary'}>
                {copied ? 'Copied!' : 'Copy'}
              </Button>
            </div>
          ) : (
            <Button size="sm" onClick={handleGenerateLink} loading={generatingLink}>
              Generate Invite Link
            </Button>
          )}
        </div>

        {/* Members Section */}
        <div>
          <h4 className="text-sm font-semibold text-gray-700 mb-2">
            Members ({members.length})
          </h4>
          {loading ? (
            <div className="flex justify-center py-4">
              <Spinner size="sm" />
            </div>
          ) : (
            <div className="space-y-1 max-h-64 overflow-y-auto">
              {members.map((m) => (
                <div key={m.user_id} className="flex items-center justify-between py-2 px-2 rounded hover:bg-gray-50">
                  <div className="flex items-center gap-2 min-w-0">
                    <Avatar name={m.user_id} size="sm" />
                    <div className="min-w-0">
                      <span className="text-sm text-gray-900 truncate block">
                        {m.user_id === user?.id ? (user.nickname || user.username) : m.user_id.slice(0, 12)}
                      </span>
                    </div>
                    <span className={`text-xs px-1.5 py-0.5 rounded-full flex-shrink-0 ${
                      m.role === 'owner' ? 'bg-yellow-100 text-yellow-700' :
                      m.role === 'admin' ? 'bg-blue-100 text-blue-700' :
                      'bg-gray-100 text-gray-600'
                    }`}>
                      {m.role}
                    </span>
                  </div>
                  {isAdmin && m.user_id !== user?.id && m.role !== 'owner' && (
                    <Button
                      size="sm"
                      variant="danger"
                      onClick={() => handleRemoveMember(m.user_id)}
                      loading={removingUser === m.user_id}
                    >
                      Remove
                    </Button>
                  )}
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </Modal>
  );
}
