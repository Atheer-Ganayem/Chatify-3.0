"use client";

import { Button } from "@/components/ui/button";
import { Paperclip } from "lucide-react";
import React, { useRef } from "react";

interface Props {
  onSelectImage: (file: File | undefined) => Promise<void>;
}

const ImageUploader: React.FC<Props> = ({ onSelectImage }) => {
  const inputRef = useRef<HTMLInputElement | null>(null);

  function onClickHandler() {
    inputRef.current?.click();
  }

  async function onChangeHandler(e: React.ChangeEvent<HTMLInputElement>) {
    const file = e.target.files?.[0];
    await onSelectImage(file);
  }

  return (
    <Button size="icon" type="button" variant="ghost" onClick={onClickHandler}>
      <Paperclip className="h-4 w-4" />
      <input type="file" hidden ref={inputRef} onChange={onChangeHandler} accept="image/*" />
    </Button>
  );
};

export default ImageUploader;
