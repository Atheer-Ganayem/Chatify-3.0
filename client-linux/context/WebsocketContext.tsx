"use client";

import { useSession } from "next-auth/react";
import { createContext, useContext, useEffect, useRef } from "react";
import { useMessages } from "./MessagesContext";
import { useOnlineUsers } from "./OnlineUsersContext";
import { getCookie } from "cookies-next";
import { toast } from "sonner";
import { useConversations } from "./ConversationsContext";
import { v4 as uuid } from "uuid";

type WebsocketContextType = {
  sendMessage: (message: string, image?: string) => void;
};

const MAX_RECONNECTS = 5;
const RECONNECT_DELAY = 3000;

const WebsocketContext = createContext<WebsocketContextType | undefined>(
  undefined
);

export default function WebsocketProvider({
  children,
}: {
  children: React.ReactNode;
}) {
  const reconnectRef = useRef<number>(0);
  const socketRef = useRef<WebSocket | null>(null);
  const session = useSession();
  const { onReceiveMessage, onDelete, appendWaitingMessage } = useMessages();
  const onlineCtx = useOnlineUsers();
  const ConversationsCtx = useConversations();
  const conversationId = ConversationsCtx.currentConversation?._id;

  const connect = () => {
    if (!session.data) {
      return;
    }

    const ws = new WebSocket(
      `${process.env.WS_URL}/ws?token=${getCookie("next-auth.session-token")}`
    );
    socketRef.current = ws;

    ws.onopen = () => {
      console.log("✅ WebSocket connected");
      reconnectRef.current = 0;
    };

    ws.onclose = () => {
      if (reconnectRef.current < MAX_RECONNECTS) {
        reconnectRef.current++;
        toast.warning(
          `Connection has been lost.\nReconnect #${reconnectRef.current} in ${
            RECONNECT_DELAY / 1000
          }s`,
          { duration: RECONNECT_DELAY }
        );
        setTimeout(connect, RECONNECT_DELAY);
      } else {
        toast.error(
          "Max reconnect attempts reached. Check your internet connection and reload the page.",
          { duration: 1000 * 60 }
        );
      }
    };

    ws.onerror = (err) => {
      console.error("❌ WebSocket error:", err);
      ws.close();
    };

    ws.onmessage = (e: MessageEvent) => {
      const data = JSON.parse(e.data) as WSResponse;
      switch (data.type) {
        case "err":
          toast.error(data.message as string);
          break;
        case "msg":
          onReceiveMessage(data);
          break;
        case "acknowledged":
          onReceiveMessage(data);
          break;
        case "delete":
          onDelete(data.messageId!);
          break;
        case "cnv":
          ConversationsCtx.appendConversation({
            _id: data.cnvId!,
            participant: data.user!,
          });
          console.log(data);
          if (data.isOnline) {
            onlineCtx.addOnline(data.user!._id);
          }
          break;
        case "status":
          if (data.online) {
            onlineCtx.addOnline(data.userId!);
          } else {
            onlineCtx.removeOnline(data.userId!);
          }
          break;
      }
    };
  };

  useEffect(() => {
    if (!session.data || ConversationsCtx.loading) return;
    if (
      socketRef.current &&
      (socketRef.current?.readyState === WebSocket.OPEN ||
        socketRef.current?.readyState === WebSocket.CONNECTING)
    )
      return;

    connect();
  }, [conversationId, session.data, ConversationsCtx.loading]);

  const sendMessage = (message: string, image?: string) => {
    if (socketRef.current && socketRef.current.readyState === WebSocket.OPEN) {
      const payload: WSRequest = {
        id: uuid(),
        conversationId: conversationId as string,
        message,
        type: "msg",
      };
      if (image) payload.image = image;
      socketRef.current.send(JSON.stringify(payload));
      appendWaitingMessage(payload);
    } else {
      console.warn("Tried to send on closed WebSocket");
    }
  };
  return (
    <WebsocketContext.Provider value={{ sendMessage }}>
      {children}
    </WebsocketContext.Provider>
  );
}

export function useWebsocket() {
  const context = useContext(WebsocketContext);
  if (!context) throw new Error("useData must be used within a DataProvider");
  return context;
}
