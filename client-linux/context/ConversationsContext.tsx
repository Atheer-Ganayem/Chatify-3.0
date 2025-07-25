"use client";

import { useSession } from "next-auth/react";
import { useSearchParams } from "next/navigation";
import React, { createContext, useContext, useEffect, useState } from "react";
import { toast } from "sonner";
import { useOnlineUsers } from "./OnlineUsersContext";
import useFetch from "@/hooks/useFetch";

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
  const params = useSearchParams();
  const onlineCtx = useOnlineUsers();
  const conversationID = params.get("conversationID");
  const currentConversation =
    conversations.find((cnv) => cnv._id === conversationID) || null;
  const { isLoading: loading, exec } = useFetch({
    path: "/conversations",
    auth: true,
    defaultLoading: true,
  });

  const fetchData = async () => {
    console.log("fetch data");

    try {
      const { ok, responseData, error } = await exec();

      if (error) throw error;
      if (ok) {
        setConversations(responseData.conversations || []);
        onlineCtx.addOnline(...(responseData.online || []));
      } else {
        toast.error(responseData.message);
      }
    } catch (err) {
      console.log(err);
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
