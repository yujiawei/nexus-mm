import { useMessagesStore } from '../../store/messages';
import { useAuthStore } from '../../store/auth';
import MessageItem from '../Chat/MessageItem';
import Spinner from '../common/Spinner';

export default function SearchResults() {
  const { searchResults, searchQuery, searching, clearSearch } = useMessagesStore();
  const { user } = useAuthStore();

  return (
    <div className="flex-1 overflow-y-auto bg-gray-50">
      <div className="flex items-center justify-between px-5 py-3 bg-white border-b border-gray-200">
        <div className="flex items-center gap-2">
          <h3 className="text-sm font-semibold text-gray-900">Search Results</h3>
          {!searching && (
            <span className="text-xs text-gray-500">
              {searchResults.length} result{searchResults.length !== 1 ? 's' : ''} for &quot;{searchQuery}&quot;
            </span>
          )}
        </div>
        <button onClick={clearSearch} className="text-sm text-blue-600 hover:text-blue-700 font-medium">
          Clear
        </button>
      </div>

      {searching ? (
        <div className="flex items-center justify-center py-16">
          <Spinner />
        </div>
      ) : searchResults.length === 0 ? (
        <div className="flex flex-col items-center justify-center py-16 text-gray-400">
          <svg className="w-12 h-12 mb-3 opacity-30" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
          <p className="text-sm">No messages found</p>
        </div>
      ) : (
        <div className="py-2">
          {searchResults.map((msg) => (
            <div key={msg.id} className="bg-white mb-1">
              <MessageItem
                message={msg}
                currentUserId={user?.id || ''}
                onOpenThread={() => {}}
                onToggleReaction={() => {}}
              />
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
