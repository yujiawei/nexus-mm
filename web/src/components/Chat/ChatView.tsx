import { useEffect, useCallback, useState } from 'react';
import { useTeamStore } from '../../store/team';
import { useMessagesStore } from '../../store/messages';
import { useAuthStore } from '../../store/auth';
import { pinMessage, unpinMessage } from '../../api/channels';
import MessageList from './MessageList';
import MessageInput from './MessageInput';
import ThreadPanel from './ThreadPanel';
import SearchBar from '../Search/SearchBar';
import SearchResults from '../Search/SearchResults';
import Button from '../common/Button';
import type { Message } from '../../api/types';

export default function ChatView() {
  const { currentChannel, joinChannel } = useTeamStore();
  const { user } = useAuthStore();
  const {
    messages,
    threadRoot,
    threadMessages,
    searchQuery,
    loading,
    hasMore,
    loadMessages,
    loadMore,
    sendMessage,
    openThread,
    closeThread,
    toggleReaction,
  } = useMessagesStore();
  const [notMember, setNotMember] = useState(false);
  const [joiningChannel, setJoiningChannel] = useState(false);

  useEffect(() => {
    if (currentChannel) {
      setNotMember(false);
      loadMessages(currentChannel.id).catch((err) => {
        if (err?.response?.status === 403) {
          setNotMember(true);
        }
      });
      closeThread();
    }
  }, [currentChannel, loadMessages, closeThread]);

  const handleJoinChannel = useCallback(async () => {
    if (!currentChannel) return;
    setJoiningChannel(true);
    try {
      await joinChannel(currentChannel.id);
      setNotMember(false);
      await loadMessages(currentChannel.id);
    } catch {
      /* ignore */
    } finally {
      setJoiningChannel(false);
    }
  }, [currentChannel, joinChannel, loadMessages]);

  const handleSend = useCallback(
    (content: string) => {
      if (currentChannel) sendMessage(currentChannel.id, content);
    },
    [currentChannel, sendMessage]
  );

  const handleSendReply = useCallback(
    (content: string) => {
      if (currentChannel && threadRoot) sendMessage(currentChannel.id, content, threadRoot.id);
    },
    [currentChannel, threadRoot, sendMessage]
  );

  const handleOpenThread = useCallback(
    (msg: Message) => {
      if (currentChannel) openThread(currentChannel.id, msg);
    },
    [currentChannel, openThread]
  );

  const handleToggleReaction = useCallback(
    (messageId: string, emoji: string) => {
      if (currentChannel && user) toggleReaction(currentChannel.id, messageId, emoji, user.id);
    },
    [currentChannel, user, toggleReaction]
  );

  const handlePinToggle = useCallback(
    async (messageId: string, isPinned: boolean) => {
      if (!currentChannel) return;
      try {
        if (isPinned) {
          await unpinMessage(currentChannel.id, messageId);
        } else {
          await pinMessage(currentChannel.id, messageId);
        }
        loadMessages(currentChannel.id);
      } catch {
        /* ignore */
      }
    },
    [currentChannel, loadMessages]
  );

  const handleLoadMore = useCallback(() => {
    if (currentChannel) loadMore(currentChannel.id);
  }, [currentChannel, loadMore]);

  if (!currentChannel) {
    return (
      <div className="flex-1 flex items-center justify-center bg-gray-50">
        <div className="text-center text-gray-400">
          <svg className="w-20 h-20 mx-auto mb-4 opacity-20" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1} d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
          </svg>
          <p className="text-lg font-medium">Select a channel to start chatting</p>
        </div>
      </div>
    );
  }

  return (
    <div className="flex-1 flex min-w-0">
      <div className="flex-1 flex flex-col min-w-0">
        <div className="flex items-center justify-between px-5 h-14 border-b border-gray-200 flex-shrink-0 bg-white">
          <div className="flex items-center gap-2 min-w-0">
            <span className="text-gray-400">
              {currentChannel.type === 'private' ? (
                <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                </svg>
              ) : (
                <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M7 20l4-16m2 16l4-16M6 9h14M4 15h14" />
                </svg>
              )}
            </span>
            <h2 className="font-semibold text-gray-900 truncate">{currentChannel.display_name}</h2>
            {currentChannel.purpose && (
              <span className="text-sm text-gray-400 truncate hidden sm:inline">
                | {currentChannel.purpose}
              </span>
            )}
          </div>
          <SearchBar />
        </div>

        {notMember ? (
          <div className="flex-1 flex items-center justify-center bg-gray-50">
            <div className="text-center">
              <p className="text-gray-500 mb-4">You are not a member of this channel.</p>
              <Button onClick={handleJoinChannel} loading={joiningChannel}>
                Join #{currentChannel.display_name}
              </Button>
            </div>
          </div>
        ) : searchQuery ? (
          <SearchResults />
        ) : (
          <>
            <MessageList
              messages={messages}
              loading={loading}
              hasMore={hasMore}
              currentUserId={user?.id || ''}
              onOpenThread={handleOpenThread}
              onToggleReaction={handleToggleReaction}
              onPinToggle={handlePinToggle}
              onLoadMore={handleLoadMore}
            />
            <MessageInput onSend={handleSend} placeholder={`Message #${currentChannel.display_name}`} />
          </>
        )}
      </div>

      {threadRoot && (
        <ThreadPanel
          rootMessage={threadRoot}
          replies={threadMessages}
          currentUserId={user?.id || ''}
          onClose={closeThread}
          onSendReply={handleSendReply}
          onToggleReaction={handleToggleReaction}
        />
      )}
    </div>
  );
}
