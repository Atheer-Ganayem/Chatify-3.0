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
import { useNotification } from "./NotificationContext";
import { useRouter } from "next/navigation";
import useFetch from "@/hooks/useFetch";

type MessagesContextType = {
  messages: Message[];
  waitingMessages: WaitingMessage[];
  loading: boolean;
  loadMore: (page: number) => Promise<void>;
  onDelete: (deletedMessageId: string) => void;
  onReceiveMessage: (payload: WSResponse) => void;
  appendWaitingMessage: (payload: WSRequest) => void;
};

const MessagesContext = createContext<MessagesContextType | undefined>(
  undefined
);

const MessagesProvider = ({ children }: { children: React.ReactNode }) => {
  const session = useSession();
  const router = useRouter();
  const [messages, setMessages] = useState<Message[]>([]);
  const [waitingMessages, setWaitingMessages] = useState<WaitingMessage[]>([]);
  const ConversationsCtx = useConversations();
  const notificationsCtx = useNotification();
  const conversationId = ConversationsCtx.currentConversation?._id;
  const currentConversationIdRef = useRef<string>(
    ConversationsCtx.currentConversation?._id
  );
  const { isLoading, exec } = useFetch({
    path: `/messages/${conversationId}`,
    auth: true,
  });

  useEffect(() => {
    currentConversationIdRef.current = conversationId;
  }, [conversationId]);

  const fetchData = async () => {
    if (!conversationId) {
      setMessages([]);
      return;
    }

    try {
      setMessages([]);
      const { ok, responseData, error } = await exec();

      if (error) throw error;
      if (ok) {
        setMessages((responseData.messages as []).reverse() || []);
      } else {
        toast.error(responseData.message);
        setMessages([]);
      }
    } catch (err) {
      console.log(err);
      setMessages([]);
    }
  };

  const onReceiveMessage = (payload: WSResponse) => {
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

  useEffect(() => {
    if (!session.data || ConversationsCtx.loading) return;
    fetchData();
  }, [conversationId]);

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
        loading: isLoading,
        loadMore,
        onDelete,
        onReceiveMessage,
        appendWaitingMessage,
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
