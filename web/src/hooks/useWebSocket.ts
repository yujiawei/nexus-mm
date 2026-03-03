import { useEffect, useRef } from 'react';
import WKSDK, { ConnectStatus } from 'wukongimjssdk';
import type { Message } from 'wukongimjssdk';
import { useMessagesStore } from '../store/messages';
import { useAuthStore } from '../store/auth';

export function useWebSocket() {
  const user = useAuthStore((s) => s.user);
  const wkToken = useAuthStore((s) => s.wkToken);
  const wsUrl = useAuthStore((s) => s.wsUrl);
  const addIncomingMessage = useMessagesStore((s) => s.addIncomingMessage);
  const connectedRef = useRef(false);

  useEffect(() => {
    if (!user?.id || !wkToken) return;

    // Derive WS address: prefer server-provided ws_url, fallback to current hostname
    let addr = wsUrl || `ws://${window.location.hostname}:15200`;
    // Replace localhost with actual hostname when accessed from browser
    if (addr.includes('localhost') || addr.includes('127.0.0.1')) {
      addr = `ws://${window.location.hostname}:15200`;
    }

    const im = WKSDK.shared();

    // Disconnect any prior session
    try { im.disconnect(); } catch { /* ignore */ }

    im.config.addr = addr;
    im.config.uid = user.id;
    im.config.token = wkToken;
    im.config.deviceFlag = 0;

    const statusListener = (status: ConnectStatus) => {
      switch (status) {
        case ConnectStatus.Connected:
          connectedRef.current = true;
          console.log('[WS] Connected to WuKongIM');
          break;
        case ConnectStatus.Disconnect:
          connectedRef.current = false;
          console.log('[WS] Disconnected from WuKongIM');
          break;
        case ConnectStatus.ConnectFail:
          connectedRef.current = false;
          console.warn('[WS] Connection failed');
          break;
        case ConnectStatus.ConnectKick:
          connectedRef.current = false;
          console.warn('[WS] Kicked by server');
          break;
      }
    };

    const messageListener = (message: Message) => {
      const content = message.content;
      const contentObj = content?.contentObj as Record<string, unknown> | undefined;
      if (!contentObj) return;

      // Build a Message object compatible with our store
      const msg = {
        id: contentObj.id as string || String(message.messageID),
        channel_id: contentObj.channel_id as string || message.channel?.channelID || '',
        user_id: contentObj.user_id as string || message.fromUID || '',
        content: contentObj.content as string || content?.conversationDigest || '',
        type: (contentObj.type as string) || 'text',
        root_id: (contentObj.root_id as string) || undefined,
        reply_count: (contentObj.reply_count as number) || 0,
        created_at: (contentObj.created_at as string) || new Date().toISOString(),
        updated_at: (contentObj.updated_at as string) || new Date().toISOString(),
        user: contentObj.user as Record<string, unknown> | undefined,
      };

      addIncomingMessage(msg as Parameters<typeof addIncomingMessage>[0]);
    };

    im.connectManager.addConnectStatusListener(statusListener);
    im.chatManager.addMessageListener(messageListener);
    im.connect();

    return () => {
      connectedRef.current = false;
      im.connectManager.removeConnectStatusListener(statusListener);
      im.chatManager.removeMessageListener(messageListener);
      try { im.disconnect(); } catch { /* ignore */ }
    };
  }, [user?.id, wkToken, wsUrl, addIncomingMessage]);
}
