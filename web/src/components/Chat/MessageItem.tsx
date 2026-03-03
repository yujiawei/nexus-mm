import { useState } from 'react';
import type { Message, ReactionGroup } from '../../api/types';
import { groupReactions } from '../../store/messages';
import Avatar from '../common/Avatar';

const QUICK_EMOJIS = ['+1', 'heart', 'smile', 'eyes', 'rocket', 'tada'];

const EMOJI_MAP: Record<string, string> = {
  '+1': '\uD83D\uDC4D',
  heart: '\u2764\uFE0F',
  smile: '\uD83D\uDE04',
  eyes: '\uD83D\uDC40',
  rocket: '\uD83D\uDE80',
  tada: '\uD83C\uDF89',
  fire: '\uD83D\uDD25',
  thinking: '\uD83E\uDD14',
  clap: '\uD83D\uDC4F',
  check: '\u2705',
};

function emojiToUnicode(name: string): string {
  return EMOJI_MAP[name] || `:${name}:`;
}

function formatTime(dateStr: string): string {
  const d = new Date(dateStr);
  return d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
}

interface MessageItemProps {
  message: Message;
  currentUserId: string;
  onOpenThread: (msg: Message) => void;
  onToggleReaction: (messageId: string, emoji: string) => void;
  onPinToggle?: (messageId: string, isPinned: boolean) => void;
}

export default function MessageItem({
  message,
  currentUserId,
  onOpenThread,
  onToggleReaction,
  onPinToggle,
}: MessageItemProps) {
  const [showActions, setShowActions] = useState(false);
  const [showEmojiPicker, setShowEmojiPicker] = useState(false);
  const reactions: ReactionGroup[] = message.reactions
    ? groupReactions(message.reactions, currentUserId)
    : [];

  const displayName = message.user?.nickname || message.user?.username || message.user_id.slice(0, 8);

  return (
    <div
      className="group relative flex gap-3 px-5 py-1.5 hover:bg-gray-50 transition-colors"
      onMouseEnter={() => setShowActions(true)}
      onMouseLeave={() => { setShowActions(false); setShowEmojiPicker(false); }}
    >
      <div className="pt-0.5 flex-shrink-0">
        <Avatar name={displayName} url={message.user?.avatar_url} size="md" />
      </div>

      <div className="flex-1 min-w-0">
        <div className="flex items-baseline gap-2">
          <span className="font-semibold text-sm text-gray-900">{displayName}</span>
          <span className="text-xs text-gray-400">{formatTime(message.created_at)}</span>
          {message.is_pinned && (
            <span className="text-xs text-amber-600 flex items-center gap-0.5">
              <svg className="w-3 h-3" fill="currentColor" viewBox="0 0 20 20">
                <path d="M9.828 3.414a2 2 0 012.829 0l1.414 1.414a2 2 0 010 2.829L12.07 9.657l3.535 3.536-1.414 1.414-3.536-3.535-2 2a2 2 0 01-2.828 0L4.414 11.657a2 2 0 010-2.829L9.828 3.414z" />
              </svg>
              Pinned
            </span>
          )}
        </div>

        <p className="text-sm text-gray-800 whitespace-pre-wrap break-words mt-0.5">{message.content}</p>

        {reactions.length > 0 && (
          <div className="flex flex-wrap gap-1 mt-1.5">
            {reactions.map((r) => (
              <button
                key={r.emoji_name}
                onClick={() => onToggleReaction(message.id, r.emoji_name)}
                className={`inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-xs border transition-colors ${
                  r.includes_me
                    ? 'bg-blue-50 border-blue-200 text-blue-700'
                    : 'bg-gray-50 border-gray-200 text-gray-600 hover:bg-gray-100'
                }`}
              >
                <span>{emojiToUnicode(r.emoji_name)}</span>
                <span>{r.count}</span>
              </button>
            ))}
          </div>
        )}

        {message.reply_count > 0 && (
          <button
            onClick={() => onOpenThread(message)}
            className="inline-flex items-center gap-1 text-xs text-blue-600 hover:text-blue-700 mt-1.5 font-medium"
          >
            <svg className="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 10h10a8 8 0 018 8v2M3 10l6 6m-6-6l6-6" />
            </svg>
            {message.reply_count} {message.reply_count === 1 ? 'reply' : 'replies'}
          </button>
        )}
      </div>

      {showActions && (
        <div className="absolute -top-3 right-4 flex items-center bg-white border border-gray-200 rounded-md shadow-sm">
          <button
            onClick={() => setShowEmojiPicker(!showEmojiPicker)}
            className="p-1.5 hover:bg-gray-100 text-gray-400 hover:text-gray-600 transition-colors"
            title="Add reaction"
          >
            <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M14.828 14.828a4 4 0 01-5.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </button>
          <button
            onClick={() => onOpenThread(message)}
            className="p-1.5 hover:bg-gray-100 text-gray-400 hover:text-gray-600 transition-colors"
            title="Reply in thread"
          >
            <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 10h10a8 8 0 018 8v2M3 10l6 6m-6-6l6-6" />
            </svg>
          </button>
          {onPinToggle && (
            <button
              onClick={() => onPinToggle(message.id, !!message.is_pinned)}
              className="p-1.5 hover:bg-gray-100 text-gray-400 hover:text-gray-600 transition-colors"
              title={message.is_pinned ? 'Unpin message' : 'Pin message'}
            >
              <svg className="w-4 h-4" fill={message.is_pinned ? 'currentColor' : 'none'} viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 5a2 2 0 012-2h10a2 2 0 012 2v16l-7-3.5L5 21V5z" />
              </svg>
            </button>
          )}
        </div>
      )}

      {showEmojiPicker && (
        <div className="absolute -top-3 right-28 bg-white border border-gray-200 rounded-md shadow-lg p-2 flex gap-1 z-10">
          {QUICK_EMOJIS.map((emoji) => (
            <button
              key={emoji}
              onClick={() => {
                onToggleReaction(message.id, emoji);
                setShowEmojiPicker(false);
              }}
              className="w-8 h-8 flex items-center justify-center hover:bg-gray-100 rounded transition-colors text-lg"
              title={`:${emoji}:`}
            >
              {emojiToUnicode(emoji)}
            </button>
          ))}
        </div>
      )}
    </div>
  );
}
