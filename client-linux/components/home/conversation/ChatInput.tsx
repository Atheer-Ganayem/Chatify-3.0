"use client";

import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import { useConversations } from "@/context/ConversationsContext";
import { useWebsocket } from "@/context/WebsocketContext";
import { Send } from "lucide-react";
import React, { useRef, useState } from "react";

const ChatInput = () => {
  const [input, setInput] = useState<string>("");
  const { currentConversation } = useConversations();
  const websocketCtx = useWebsocket();
  const btnRef = useRef<HTMLButtonElement>(null);

  function onSubmitHandler(e: React.FormEvent) {
    e.preventDefault();
    if (input.trim() === "") {
      return;
    }
    websocketCtx.sendMessage(input.trim());
    setInput("");
  }

  function onKeyHandler(e: React.KeyboardEvent<HTMLTextAreaElement>) {
    if (e.key === "Enter") {
      if (e.shiftKey) {
        return;
      }
      e.preventDefault(); // prevent newline
      btnRef.current?.click();
    }
  }

  return (
    <form
      className="flex justify-center w-full max-w-full"
      onSubmit={onSubmitHandler}
    >
      <div className="relative w-full">
        <Textarea
          onKeyDown={onKeyHandler}
          placeholder={`Type something to send to ${currentConversation?.participant.name}...`}
          maxLength={1000}
          className="text-xl w-full pr-12 resize-none overflow-y-auto overflow-x-hidden max-h-40 min-h-[56px] break-words break-all"
          value={input}
          onChange={(e) => setInput(e.target.value)}
        />
        <Button
          size="icon"
          className="absolute bottom-2 right-2"
          type="submit"
          ref={btnRef}
        >
          <Send className="h-4 w-4" />
        </Button>
      </div>
    </form>
  );
};

export default ChatInput;
