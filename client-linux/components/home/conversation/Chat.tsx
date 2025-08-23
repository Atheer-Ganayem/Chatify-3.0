"use client";

import React, { useEffect, useRef, useState } from "react";
import ChatInput from "./ChatInput";
import MessageCard from "./MessageCard";
import WaitingMessageCard from "./WaitingMessageCard";
import { useMessages } from "@/context/MessagesContext";
import ChatLoader from "@/components/loaders/ChatLoader";
import { useConversations } from "@/context/ConversationsContext";
import { MessageCircleIcon } from "lucide-react";
import { Button } from "@/components/ui/button";
import { useRouter, useSearchParams } from "next/navigation";
import { toast } from "sonner";

const Chat = () => {
  const router = useRouter();
  const [loadingMore, setLoadingMore] = useState<boolean>(false);
  const params = useSearchParams();
  const pageQuery = Number(params.get("page"));
  const [page, setPage] = useState<number>(1);
  const bottomRef = useRef<HTMLSpanElement | null>(null);
  const ctx = useMessages();
  const currentConversation = useConversations().currentConversation;

  useEffect(() => {
    router.push(`/?conversationID=${params.get("conversationID")}&page=${1}`);
  }, []);

  useEffect(() => {
    if (!isNaN(Number(pageQuery)) && pageQuery > 0) {
      setPage((prev) => {
        if (prev + 1 === pageQuery) {
          return pageQuery;
        }
        return prev;
      });
    }
  }, [pageQuery]);

  useEffect(() => {
    if (!loadingMore) {
      bottomRef.current?.scrollIntoView({ behavior: "smooth" });
    }
    setLoadingMore(false);
  }, [ctx.messages.length, ctx.waitingMessages.length]);

  async function onClickHandler() {
    try {
      setLoadingMore(true);
      await ctx.loadMore(page + 1);
    } catch (error) {
      console.log(error);
      toast.error("Coudn't load older messages.");
    }
  }
  return (
    <>
      <div className="flex-1 flex flex-col overflow-hidden px-4 py-2">
        <div className="flex-1 overflow-y-auto">
          {ctx.loading && currentConversation && <ChatLoader />}
          {!currentConversation && <PleaceHolder />}
          {ctx.messages.length >= 30 && (
            <div className="flex justify-center my-5">
              <Button variant="outline" onClick={onClickHandler}>
                Load More
              </Button>
            </div>
          )}
          {ctx.messages.map((msg) => (
            <MessageCard key={msg._id} message={msg} />
          ))}
          {ctx.waitingMessages.map((msg) => (
            <WaitingMessageCard key={msg.requestId} text={msg.text} />
          ))}
          <span ref={bottomRef} />
        </div>
        {currentConversation && !ctx.loading && <ChatInput />}
      </div>
    </>
  );
};

function PleaceHolder() {
  return (
    <div className="h-full w-full flex flex-col justify-center items-center text-center gap-4">
      <MessageCircleIcon size={70} />
      <p className="text-2xl font-semibold">
        Choose or create a conversation to start chatting.
      </p>
    </div>
  );
}

export default Chat;
