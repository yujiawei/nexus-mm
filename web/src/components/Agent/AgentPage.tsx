import { useState, useEffect, useCallback } from 'react';
import { useAuthStore } from '../../store/auth';
import * as botsApi from '../../api/bots';
import type { BotInfo } from '../../api/types';
import Button from '../common/Button';

export default function AgentPage() {
  const { user } = useAuthStore();
  const [bots, setBots] = useState<BotInfo[]>([]);
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [creating, setCreating] = useState(false);
  const [error, setError] = useState('');
  const [copiedId, setCopiedId] = useState<string | null>(null);
  const [editingWebhook, setEditingWebhook] = useState<string | null>(null);
  const [webhookUrl, setWebhookUrl] = useState('');

  const loadBots = useCallback(async () => {
    try {
      const data = await botsApi.listBots();
      setBots(data || []);
    } catch {
      // ignore
    }
  }, []);

  useEffect(() => {
    loadBots();
  }, [loadBots]);

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim()) return;
    setCreating(true);
    setError('');
    try {
      await botsApi.createBot({ name: name.trim(), description: description.trim() });
      setName('');
      setDescription('');
      await loadBots();
    } catch (err: unknown) {
      const message = (err as { response?: { data?: { error?: string } } })?.response?.data?.error || 'Failed to create bot';
      setError(message);
    } finally {
      setCreating(false);
    }
  };

  const handleCopy = (text: string, id: string) => {
    navigator.clipboard.writeText(text);
    setCopiedId(id);
    setTimeout(() => setCopiedId(null), 2000);
  };

  const handleRegenerate = async (botId: string) => {
    try {
      await botsApi.regenerateToken(botId);
      await loadBots();
    } catch {
      // ignore
    }
  };

  const handleSaveWebhook = async (botId: string) => {
    try {
      await botsApi.updateWebhook(botId, webhookUrl);
      setEditingWebhook(null);
      setWebhookUrl('');
      await loadBots();
    } catch {
      // ignore
    }
  };

  const baseUrl = window.location.origin;
  const quickPrompt = `Read ${baseUrl}/skill.md to install Nexus-MM channel messaging. My email is ${user?.email || 'your@email.com'}`;

  return (
    <div className="flex-1 overflow-y-auto bg-gray-50 p-6">
      <div className="max-w-3xl mx-auto space-y-6">
        <h1 className="text-2xl font-bold text-gray-900">My Agents</h1>

        {/* Create New Bot */}
        <div className="bg-white rounded-lg border border-gray-200 p-6">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">Create New Bot</h2>
          <form onSubmit={handleCreate} className="space-y-3">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Name</label>
              <input
                type="text"
                value={name}
                onChange={(e) => setName(e.target.value)}
                placeholder="MyAgent"
                className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm text-gray-900 focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Description</label>
              <input
                type="text"
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                placeholder="What does this bot do?"
                className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm text-gray-900 focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>
            {error && <p className="text-sm text-red-600">{error}</p>}
            <Button type="submit" loading={creating} disabled={!name.trim()}>
              Create Bot
            </Button>
          </form>
        </div>

        {/* My Bots */}
        <div className="bg-white rounded-lg border border-gray-200 p-6">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">My Bots</h2>
          {bots.length === 0 ? (
            <p className="text-sm text-gray-500">No bots created yet.</p>
          ) : (
            <div className="space-y-4">
              {bots.map((bot) => (
                <div key={bot.id} className="border border-gray-100 rounded-md p-4 bg-gray-50">
                  <div className="flex items-center gap-2 mb-2">
                    <span className="text-lg">🤖</span>
                    <span className="font-medium text-gray-900">{bot.nickname}</span>
                    <span className="text-xs text-gray-500">@{bot.username}</span>
                  </div>
                  {bot.description && (
                    <p className="text-sm text-gray-600 mb-2">{bot.description}</p>
                  )}
                  <div className="space-y-2">
                    <div className="flex items-center gap-2">
                      <span className="text-xs font-medium text-gray-500">Token:</span>
                      <code className="text-xs bg-gray-200 px-2 py-0.5 rounded text-gray-800 font-mono">
                        {bot.token.slice(0, 12)}...
                      </code>
                      <button
                        onClick={() => handleCopy(bot.token, `token-${bot.id}`)}
                        className="text-xs text-blue-600 hover:text-blue-800"
                      >
                        {copiedId === `token-${bot.id}` ? 'Copied!' : 'Copy'}
                      </button>
                      <button
                        onClick={() => handleRegenerate(bot.id)}
                        className="text-xs text-orange-600 hover:text-orange-800"
                      >
                        Regenerate
                      </button>
                    </div>
                    <div className="flex items-center gap-2">
                      <span className="text-xs font-medium text-gray-500">Webhook:</span>
                      {editingWebhook === bot.id ? (
                        <div className="flex items-center gap-1">
                          <input
                            type="text"
                            value={webhookUrl}
                            onChange={(e) => setWebhookUrl(e.target.value)}
                            placeholder="https://..."
                            className="text-xs px-2 py-1 border border-gray-300 rounded text-gray-900"
                          />
                          <button
                            onClick={() => handleSaveWebhook(bot.id)}
                            className="text-xs text-green-600 hover:text-green-800"
                          >
                            Save
                          </button>
                          <button
                            onClick={() => setEditingWebhook(null)}
                            className="text-xs text-gray-500 hover:text-gray-700"
                          >
                            Cancel
                          </button>
                        </div>
                      ) : (
                        <>
                          <span className="text-xs text-gray-600">
                            {bot.webhook_url || '(not set)'}
                          </span>
                          <button
                            onClick={() => {
                              setEditingWebhook(bot.id);
                              setWebhookUrl(bot.webhook_url || '');
                            }}
                            className="text-xs text-blue-600 hover:text-blue-800"
                          >
                            Edit
                          </button>
                        </>
                      )}
                    </div>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>

        {/* Quick Connect */}
        <div className="bg-white rounded-lg border border-gray-200 p-6">
          <h2 className="text-lg font-semibold text-gray-900 mb-2">Quick Connect</h2>
          <p className="text-sm text-gray-600 mb-3">
            Copy this prompt to your AI agent to connect it to Nexus-MM:
          </p>
          <div className="bg-gray-100 rounded-md p-3 font-mono text-sm text-gray-800 whitespace-pre-wrap">
            {quickPrompt}
          </div>
          <div className="mt-3">
            <Button
              variant="secondary"
              size="sm"
              onClick={() => handleCopy(quickPrompt, 'prompt')}
            >
              {copiedId === 'prompt' ? 'Copied!' : 'Copy Prompt'}
            </Button>
          </div>
        </div>

        {/* Bot API Reference */}
        <div className="bg-white rounded-lg border border-gray-200 p-6">
          <h2 className="text-lg font-semibold text-gray-900 mb-3">Bot API Reference</h2>
          <div className="space-y-3 text-sm text-gray-700">
            <div>
              <code className="text-xs bg-blue-50 text-blue-800 px-2 py-0.5 rounded font-mono">GET /bot/&#123;token&#125;/getMe</code>
              <p className="text-gray-500 mt-1">Returns bot info.</p>
            </div>
            <div>
              <code className="text-xs bg-blue-50 text-blue-800 px-2 py-0.5 rounded font-mono">POST /bot/&#123;token&#125;/sendMessage</code>
              <p className="text-gray-500 mt-1">Send a message. Body: &#123;"channel_id", "content", "root_id"&#125;</p>
            </div>
            <div>
              <code className="text-xs bg-blue-50 text-blue-800 px-2 py-0.5 rounded font-mono">GET /bot/&#123;token&#125;/getUpdates</code>
              <p className="text-gray-500 mt-1">Poll for new messages. Query: ?offset=0&limit=100</p>
            </div>
            <div>
              <code className="text-xs bg-blue-50 text-blue-800 px-2 py-0.5 rounded font-mono">POST /bot/&#123;token&#125;/setWebhook</code>
              <p className="text-gray-500 mt-1">Set webhook URL. Body: &#123;"url"&#125;</p>
            </div>
            <div>
              <code className="text-xs bg-blue-50 text-blue-800 px-2 py-0.5 rounded font-mono">POST /bot/&#123;token&#125;/sendReaction</code>
              <p className="text-gray-500 mt-1">Add reaction. Body: &#123;"channel_id", "message_id", "emoji_name"&#125;</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
