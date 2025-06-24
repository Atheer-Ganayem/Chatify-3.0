"use client";

import { getCookie } from "cookies-next";
import { useSession } from "next-auth/react";
import { useSearchParams } from "next/navigation";
import React, { createContext, useContext, useEffect, useState } from "react";
import { toast } from "sonner";

type ConversationContextType = {
  conversations: Conversation[];
  loading: boolean;
  appendConversation: (cnv: Conversation) => void;
  currentConversation: Conversation | null;
  updateLastMessage: (message: Message) => void;
  getSenderName: (receivedMessageConversationID: string) => string;
};

const ConversationsContext = createContext<ConversationContextType | undefined>(
  undefined
);

const ConversationsProvider = ({ children }: { children: React.ReactNode }) => {
  const session = useSession();
  const [conversations, setConversations] = useState<Conversation[]>([]);
  const [loading, setLoading] = useState(true);
  const params = useSearchParams();
  const conversationID = params.get("conversationID");
  const currentConversation =
    conversations.find((cnv) => cnv._id === conversationID) || null;

  const fetchData = async () => {
    try {
      setLoading(true);
      const response = await fetch(`${process.env.BACKEND_URL}/conversations`, {
        credentials: "include",
        headers: {
          Authorization: `Bearer ${getCookie("next-auth.session-token")}`,
        },
      });
      const responseData = await response.json();
      if (response.ok) {
        setConversations(responseData.conversations || []);
      } else {
        toast.error(responseData.message);
      }
    } catch (err) {
      console.log(err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (session.data && loading) {
      fetchData();
    }
  }, [session.data]);

  function appendConversation(newCnv: Conversation) {
    setConversations((prev) => [newCnv, ...prev]);
  }

  function updateLastMessage(message: Message) {
    setConversations((prev) => {
      const index = prev.findIndex((cnv) => cnv._id === message.conversationId);
      const copy = [...prev];
      if (index >= 0) {
        copy[index].lastMessage = message;
      }
      return copy;
    });
  }

  function getSenderName(receivedMsgCnvID: string) {
    const cnv = conversations.find((cnv) => cnv._id === receivedMsgCnvID);
    if (cnv) {
      return cnv.participant.name;
    }
    return "";
  }

  return (
    <ConversationsContext.Provider
      value={{
        conversations,
        loading,
        appendConversation,
        currentConversation: currentConversation,
        updateLastMessage,
        getSenderName,
      }}
    >
      {children}
    </ConversationsContext.Provider>
  );
};

export const useConversations = () => {
  const context = useContext(ConversationsContext);
  if (!context) throw new Error("useData must be used within a DataProvider");
  return context;
};

export default ConversationsProvider;
