import { useState, useRef, useEffect, type FormEvent, type KeyboardEvent } from 'react';

const SLASH_COMMANDS = [
  { trigger: '/join', description: 'Join a channel' },
  { trigger: '/leave', description: 'Leave current channel' },
  { trigger: '/invite', description: 'Invite a user' },
  { trigger: '/mute', description: 'Mute a channel' },
  { trigger: '/help', description: 'Show help' },
];

interface MessageInputProps {
  onSend: (content: string) => void;
  placeholder?: string;
  disabled?: boolean;
}

export default function MessageInput({ onSend, placeholder, disabled }: MessageInputProps) {
  const [value, setValue] = useState('');
  const [showSlash, setShowSlash] = useState(false);
  const [slashFilter, setSlashFilter] = useState('');
  const [selectedIdx, setSelectedIdx] = useState(0);
  const inputRef = useRef<HTMLTextAreaElement>(null);

  const filteredCommands = SLASH_COMMANDS.filter((cmd) =>
    cmd.trigger.startsWith(slashFilter || '/')
  );

  useEffect(() => {
    if (value.startsWith('/') && !value.includes(' ')) {
      setShowSlash(true);
      setSlashFilter(value);
      setSelectedIdx(0);
    } else {
      setShowSlash(false);
    }
  }, [value]);

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault();
    if (!value.trim() || disabled) return;
    onSend(value.trim());
    setValue('');
    setShowSlash(false);
  };

  const handleKeyDown = (e: KeyboardEvent) => {
    if (showSlash && filteredCommands.length > 0) {
      if (e.key === 'ArrowDown') {
        e.preventDefault();
        setSelectedIdx((i) => Math.min(i + 1, filteredCommands.length - 1));
        return;
      }
      if (e.key === 'ArrowUp') {
        e.preventDefault();
        setSelectedIdx((i) => Math.max(i - 1, 0));
        return;
      }
      if (e.key === 'Tab' || (e.key === 'Enter' && showSlash)) {
        e.preventDefault();
        setValue(filteredCommands[selectedIdx].trigger + ' ');
        setShowSlash(false);
        return;
      }
    }

    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSubmit(e);
    }
  };

  const selectCommand = (trigger: string) => {
    setValue(trigger + ' ');
    setShowSlash(false);
    inputRef.current?.focus();
  };

  return (
    <div className="relative px-5 pb-4">
      {showSlash && filteredCommands.length > 0 && (
        <div className="absolute bottom-full left-5 right-5 mb-1 bg-white border border-gray-200 rounded-lg shadow-lg overflow-hidden z-10">
          {filteredCommands.map((cmd, i) => (
            <button
              key={cmd.trigger}
              onClick={() => selectCommand(cmd.trigger)}
              className={`flex items-center gap-3 w-full px-4 py-2.5 text-left text-sm transition-colors ${
                i === selectedIdx ? 'bg-blue-50' : 'hover:bg-gray-50'
              }`}
            >
              <span className="font-mono font-medium text-blue-600">{cmd.trigger}</span>
              <span className="text-gray-500">{cmd.description}</span>
            </button>
          ))}
        </div>
      )}

      <form onSubmit={handleSubmit} className="flex items-end gap-2">
        <div className="flex-1 relative">
          <textarea
            ref={inputRef}
            value={value}
            onChange={(e) => setValue(e.target.value)}
            onKeyDown={handleKeyDown}
            placeholder={placeholder || 'Write a message...'}
            disabled={disabled}
            rows={1}
            className="w-full resize-none border border-gray-300 rounded-lg px-4 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent disabled:bg-gray-50 disabled:text-gray-400"
            style={{ minHeight: '42px', maxHeight: '120px' }}
          />
        </div>
        <button
          type="submit"
          disabled={!value.trim() || disabled}
          className="flex-shrink-0 p-2.5 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
        >
          <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8" />
          </svg>
        </button>
      </form>
    </div>
  );
}
