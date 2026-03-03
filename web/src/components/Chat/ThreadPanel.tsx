import type { Message } from '../../api/types';
import MessageItem from './MessageItem';
import MessageInput from './MessageInput';

interface ThreadPanelProps {
  rootMessage: Message;
  replies: Message[];
  currentUserId: string;
  onClose: () => void;
  onSendReply: (content: string) => void;
  onToggleReaction: (messageId: string, emoji: string) => void;
}

export default function ThreadPanel({
  rootMessage,
  replies,
  currentUserId,
  onClose,
  onSendReply,
  onToggleReaction,
}: ThreadPanelProps) {
  return (
    <div className="w-96 border-l border-gray-200 bg-white flex flex-col flex-shrink-0">
      <div className="flex items-center justify-between px-4 h-14 border-b border-gray-200 flex-shrink-0">
        <div>
          <h3 className="text-sm font-semibold text-gray-900">Thread</h3>
          <p className="text-xs text-gray-500">
            {replies.length} {replies.length === 1 ? 'reply' : 'replies'}
          </p>
        </div>
        <button
          onClick={onClose}
          className="p-1 text-gray-400 hover:text-gray-600 transition-colors"
        >
          <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <div className="flex-1 overflow-y-auto">
        <div className="border-b border-gray-100 pb-2 mb-2">
          <MessageItem
            message={rootMessage}
            currentUserId={currentUserId}
            onOpenThread={() => {}}
            onToggleReaction={onToggleReaction}
          />
        </div>

        {replies.map((msg) => (
          <MessageItem
            key={msg.id}
            message={msg}
            currentUserId={currentUserId}
            onOpenThread={() => {}}
            onToggleReaction={onToggleReaction}
          />
        ))}
      </div>

      <MessageInput onSend={onSendReply} placeholder="Reply..." />
    </div>
  );
}
