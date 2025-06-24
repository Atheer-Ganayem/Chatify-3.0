"use client";

import React, {
  createContext,
  useContext,
  useEffect,
  useRef,
  useState,
} from "react";
import { toast } from "sonner";
import { useConversations } from "./ConversationsContext";
import { useSession } from "next-auth/react";
import { v4 as uuid } from "uuid";
import { useNotification } from "./NotificationContext";
import { useRouter } from "next/navigation";

type MessagesContextType = {
  messages: Message[];
  waitingMessages: WaitingMessage[];
  loading: boolean;
  send: (data: string) => void;
  loadMore: (page: number) => Promise<void>;
  onDelete: (deletedMessageId: string) => void;
};

const MAX_RECONNECTS = 5;
const RECONNECT_DELAY = 3000;

const MessagesContext = createContext<MessagesContextType | undefined>(
  undefined
);

const MessagesProvider = ({ children }: { children: React.ReactNode }) => {
  const session = useSession();
  const router = useRouter();
  const [messages, setMessages] = useState<Message[]>([]);
  const [waitingMessages, setWaitingMessages] = useState<WaitingMessage[]>([]);
  const [loading, setLoading] = useState(true);
  const ConversationsCtx = useConversations();
  const notificationsCtx = useNotification();
  const conversationId = ConversationsCtx.currentConversation?._id;
  const currentConversationIdRef = useRef<string>(
    ConversationsCtx.currentConversation?._id
  );
  const socketRef = useRef<WebSocket | null>(null);
  const reconnectRef = useRef<number>(0);

  useEffect(() => {
    currentConversationIdRef.current = conversationId;
  }, [conversationId]);

  const fetchData = async () => {
    if (!conversationId) {
      setMessages([]);
      setLoading(false);
      return;
    }

    try {
      setMessages([]);
      setLoading(true);
      const response = await fetch(
        `${process.env.BACKEND_URL}/messages/${conversationId}`,
        {
          credentials: "include",
        }
      );
      const responseData = await response.json();
      if (response.ok) {
        setMessages((responseData.messages as []).reverse() || []);
      } else {
        toast.error(responseData.message);
        setMessages([]);
      }
    } catch (err) {
      console.log(err);
      setMessages([]);
    } finally {
      setLoading(false);
    }
  };

  const onReceive = (payload: WSResponse) => {
    if (payload.type === "acknowledged") {
      setWaitingMessages((prev) =>
        prev.filter((msg) => msg.requestId !== payload.id)
      );
    }
    if (
      (payload.message as Message).conversationId ===
      currentConversationIdRef.current
    ) {
      setMessages((prev) => [...prev, payload.message as Message]);
    } else {
      const name = ConversationsCtx.getSenderName(
        (payload.message as Message).conversationId
      );
      if (name) {
        notificationsCtx.notify(`New message from: ${name}`);
      }
    }
    ConversationsCtx.updateLastMessage(payload.message as Message);
  };

  const connect = () => {
    if (!session.data) {
      return;
    }

    const ws = new WebSocket(`${process.env.WS_URL}/ws`);
    socketRef.current = ws;

    ws.onopen = () => {
      console.log("✅ WebSocket connected");
      reconnectRef.current = 0;
    };

    ws.onclose = () => {
      if (reconnectRef.current < MAX_RECONNECTS) {
        reconnectRef.current++;
        console.warn(
          `⚠️ Reconnect #${reconnectRef.current} in ${RECONNECT_DELAY / 1000}s`
        );
        setTimeout(connect, RECONNECT_DELAY);
      } else {
        console.error("❌ Max reconnect attempts reached.");
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
          onReceive(data);
          break;
        case "acknowledged":
          onReceive(data);
          break;
        case "delete":
          onDelete(data.messageId!);
          break;
      }
    };
  };

  useEffect(() => {
    if (!session.data || ConversationsCtx.loading) return;
    fetchData();
  }, [conversationId]);

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

  const send = (message: string) => {
    if (socketRef.current && socketRef.current.readyState === WebSocket.OPEN) {
      const payload: WSRequest = {
        id: uuid(),
        conversationId: conversationId as string,
        message,
      };
      socketRef.current.send(JSON.stringify(payload));
      appendWaitingMessage(payload);
    } else {
      console.warn("Tried to send on closed WebSocket");
    }
  };

  function appendWaitingMessage(payload: WSRequest) {
    setWaitingMessages((prev) => [
      ...prev,
      {
        requestId: payload.id,
        conversationID: payload.conversationId,
        text: payload.message as string,
        createdAt: new Date(Date.now()),
      },
    ]);
  }

  async function loadMore(page: number) {
    if (page <= 0) return;
    try {
      const response = await fetch(
        `${process.env.BACKEND_URL}/messages/${conversationId}?page=${page}`,
        {
          credentials: "include",
        }
      );
      const responseData = await response.json();

      if (
        response.ok &&
        responseData.messages &&
        responseData.messages.length > 0
      ) {
        setMessages((prev) => [...responseData.messages.reverse(), ...prev]);
        router.push(`/?conversationID=${conversationId}&page=${page}`);
      } else if (response.ok && !responseData.messages) {
        toast.info("No more messages to load.");
      } else {
        toast.error(responseData.message);
      }
    } catch (err) {
      console.log(err);
      toast.error("Coudn't load older messages.");
    }
  }

  function onDelete(deletedMessageId: string) {
    setMessages((prev) => prev.filter((msg) => msg._id !== deletedMessageId));
  }

  return (
    <MessagesContext.Provider
      value={{
        messages,
        waitingMessages,
        loading,
        send,
        loadMore,
        onDelete,
      }}
    >
      {children}
    </MessagesContext.Provider>
  );
};

export const useMessages = () => {
  const context = useContext(MessagesContext);
  if (!context) throw new Error("useData must be used within a DataProvider");
  return context;
};

export default MessagesProvider;
