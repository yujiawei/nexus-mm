import { useEffect, useState, useCallback } from 'react';
import { useAuthStore } from '../../store/auth';
import { listChannelMembers, removeChannelMember } from '../../api/channels';
import type { ChannelMember, Channel } from '../../api/types';
import Avatar from '../common/Avatar';
import Button from '../common/Button';
import Spinner from '../common/Spinner';

interface Props {
  channel: Channel;
  onClose: () => void;
}

export default function ChannelMembersPanel({ channel, onClose }: Props) {
  const { user } = useAuthStore();
  const [members, setMembers] = useState<ChannelMember[]>([]);
  const [loading, setLoading] = useState(true);
  const [removingUser, setRemovingUser] = useState<string | null>(null);

  const loadMembers = useCallback(async () => {
    setLoading(true);
    try {
      const data = await listChannelMembers(channel.id);
      setMembers(data || []);
    } finally {
      setLoading(false);
    }
  }, [channel.id]);

  useEffect(() => {
    loadMembers();
  }, [loadMembers]);

  const isCreator = channel.creator_id === user?.id;

  const handleRemove = async (userId: string) => {
    setRemovingUser(userId);
    try {
      await removeChannelMember(channel.id, userId);
      await loadMembers();
    } finally {
      setRemovingUser(null);
    }
  };

  return (
    <div className="w-64 border-l border-gray-200 bg-white flex flex-col flex-shrink-0">
      <div className="flex items-center justify-between px-4 h-14 border-b border-gray-200 flex-shrink-0">
        <h3 className="text-sm font-semibold text-gray-900">Members ({members.length})</h3>
        <button onClick={onClose} className="text-gray-400 hover:text-gray-600">
          <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>
      <div className="flex-1 overflow-y-auto p-3">
        {loading ? (
          <div className="flex justify-center py-4">
            <Spinner size="sm" />
          </div>
        ) : (
          <div className="space-y-1">
            {members.map((m) => (
              <div key={m.user_id} className="flex items-center justify-between py-1.5 px-2 rounded hover:bg-gray-50 group">
                <div className="flex items-center gap-2 min-w-0">
                  <Avatar name={m.user_id} size="sm" />
                  <span className="text-sm text-gray-700 truncate">
                    {m.user_id === user?.id ? (user.nickname || user.username) : m.user_id.slice(0, 10)}
                  </span>
                </div>
                {isCreator && m.user_id !== user?.id && (
                  <button
                    onClick={() => handleRemove(m.user_id)}
                    disabled={removingUser === m.user_id}
                    className="text-gray-300 hover:text-red-500 opacity-0 group-hover:opacity-100 transition-opacity"
                    title="Remove member"
                  >
                    <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                    </svg>
                  </button>
                )}
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
