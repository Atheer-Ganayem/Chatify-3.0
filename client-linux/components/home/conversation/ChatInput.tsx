"use client";

import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import { useConversations } from "@/context/ConversationsContext";
import { useWebsocket } from "@/context/WebsocketContext";
import { Send } from "lucide-react";
import React, { useRef, useState } from "react";
import ImageUploader from "./ImageUploader";
import useFetch from "@/hooks/useFetch";
import { toast } from "sonner";
import ImagePreviewer from "./ImagePreviewer";

const ChatInput = () => {
  const [input, setInput] = useState<string>("");
  const [image, setImage] = useState<string>(""); // this is the image path we receive when uploading it.
  const [previewUrl, setPreviewUrl] = useState<string | undefined>();
  const { currentConversation } = useConversations();
  const websocketCtx = useWebsocket();
  const btnRef = useRef<HTMLButtonElement>(null);
  const { isLoading: isUploading, exec } = useFetch({
    path: "/upload",
    method: "POST",
    auth: true,
  });

  function onSubmitHandler(e: React.FormEvent) {
    e.preventDefault();
    if (input.trim() === "") {
      return;
    }
    websocketCtx.sendMessage(input.trim(), image);
    setInput("");
    reset();
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

  async function onSelectImage(file: File | undefined) {
    if (isUploading) return;
    if (file) {
      setPreviewUrl(URL.createObjectURL(file));
    } else {
      setPreviewUrl(undefined);
      setImage("");
      return;
    }
    try {
      const formData = new FormData();
      formData.append("image", file);
      const { ok, responseData, error } = await exec(formData);
      if (error) throw error;

      if (ok) {
        setImage(responseData.path);
        toast.success(responseData.message);
      } else {
        toast.error(responseData.message);
        reset();
      }
    } catch (error) {
      console.log(error);
      toast.error("Couldn't upload image");
      reset();
    }
  }

  function reset() {
    setPreviewUrl(undefined);
    setImage("");
  }

  return (
    <div>
      {previewUrl && (
        <ImagePreviewer url={previewUrl} isLoading={isUploading} />
      )}
      <form
        className="flex justify-center w-full max-w-full"
        onSubmit={onSubmitHandler}
      >
        <div className="relative w-full">
          <Textarea
            onKeyDown={onKeyHandler}
            placeholder={`Type something to send to ${currentConversation?.participant.name}...`}
            maxLength={1000}
            className="text-xl w-full pr-20 resize-none overflow-y-auto overflow-x-hidden max-h-40 min-h-[56px] break-words break-all"
            value={input}
            onChange={(e) => setInput(e.target.value)}
          />
          <div className="absolute bottom-2 right-2 flex gap-2 flex-row-reverse">
            <Button
              size="icon"
              type="submit"
              ref={btnRef}
              disabled={!input.trim()}
            >
              <Send className="h-4 w-4" />
            </Button>
            <ImageUploader onSelectImage={onSelectImage} />
          </div>
        </div>
      </form>
    </div>
  );
};

export default ChatInput;
