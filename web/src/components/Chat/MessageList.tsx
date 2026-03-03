import { useEffect, useRef, useCallback } from 'react';
import type { Message } from '../../api/types';
import MessageItem from './MessageItem';
import Spinner from '../common/Spinner';

interface MessageListProps {
  messages: Message[];
  loading: boolean;
  hasMore: boolean;
  currentUserId: string;
  onOpenThread: (msg: Message) => void;
  onToggleReaction: (messageId: string, emoji: string) => void;
  onPinToggle: (messageId: string, isPinned: boolean) => void;
  onLoadMore: () => void;
}

export default function MessageList({
  messages,
  loading,
  hasMore,
  currentUserId,
  onOpenThread,
  onToggleReaction,
  onPinToggle,
  onLoadMore,
}: MessageListProps) {
  const bottomRef = useRef<HTMLDivElement>(null);
  const containerRef = useRef<HTMLDivElement>(null);
  const prevLenRef = useRef(0);

  useEffect(() => {
    if (messages.length > prevLenRef.current) {
      bottomRef.current?.scrollIntoView({ behavior: prevLenRef.current === 0 ? 'auto' : 'smooth' });
    }
    prevLenRef.current = messages.length;
  }, [messages.length]);

  const handleScroll = useCallback(() => {
    const el = containerRef.current;
    if (!el || !hasMore) return;
    if (el.scrollTop < 100) {
      onLoadMore();
    }
  }, [hasMore, onLoadMore]);

  useEffect(() => {
    const el = containerRef.current;
    if (!el) return;
    el.addEventListener('scroll', handleScroll, { passive: true });
    return () => el.removeEventListener('scroll', handleScroll);
  }, [handleScroll]);

  if (loading && messages.length === 0) {
    return (
      <div className="flex-1 flex items-center justify-center">
        <Spinner />
      </div>
    );
  }

  return (
    <div ref={containerRef} className="flex-1 overflow-y-auto">
      {hasMore && messages.length > 0 && (
        <div className="flex justify-center py-4">
          <button
            onClick={onLoadMore}
            className="text-xs text-blue-600 hover:text-blue-700 font-medium"
          >
            Load older messages
          </button>
        </div>
      )}

      {messages.length === 0 && !loading && (
        <div className="flex flex-col items-center justify-center h-full text-gray-400">
          <svg className="w-16 h-16 mb-4 opacity-30" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1} d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
          </svg>
          <p className="text-sm">No messages yet. Start the conversation!</p>
        </div>
      )}

      <div className="py-2">
        {messages.map((msg) => (
          <MessageItem
            key={msg.id}
            message={msg}
            currentUserId={currentUserId}
            onOpenThread={onOpenThread}
            onToggleReaction={onToggleReaction}
            onPinToggle={onPinToggle}
          />
        ))}
      </div>
      <div ref={bottomRef} />
    </div>
  );
}
