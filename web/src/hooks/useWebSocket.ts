import { useEffect, useRef, useCallback } from 'react';
import { useMessagesStore } from '../store/messages';

const RECONNECT_DELAY = 3000;

export function useWebSocket(userId: string | undefined, channelId: string | undefined) {
  const wsRef = useRef<WebSocket | null>(null);
  const reconnectTimer = useRef<ReturnType<typeof setTimeout>>();
  const addIncomingMessage = useMessagesStore((s) => s.addIncomingMessage);

  const connect = useCallback(() => {
    if (!userId) return;

    const wsUrl = import.meta.env.VITE_WS_URL || `ws://${window.location.hostname}:5200`;
    const url = `${wsUrl}/ws?uid=${userId}`;

    try {
      const ws = new WebSocket(url);
      wsRef.current = ws;

      ws.onopen = () => {
        if (reconnectTimer.current) {
          clearTimeout(reconnectTimer.current);
        }
      };

      ws.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data);
          if (data.content) {
            const payload = typeof data.content === 'string' ? JSON.parse(data.content) : data.content;
            if (payload.channel_id === channelId || !channelId) {
              addIncomingMessage(payload);
            }
          }
        } catch {
          /* ignore parse errors */
        }
      };

      ws.onclose = () => {
        reconnectTimer.current = setTimeout(connect, RECONNECT_DELAY);
      };

      ws.onerror = () => {
        ws.close();
      };
    } catch {
      reconnectTimer.current = setTimeout(connect, RECONNECT_DELAY);
    }
  }, [userId, channelId, addIncomingMessage]);

  useEffect(() => {
    connect();
    return () => {
      if (reconnectTimer.current) clearTimeout(reconnectTimer.current);
      if (wsRef.current) {
        wsRef.current.onclose = null;
        wsRef.current.close();
      }
    };
  }, [connect]);
}
