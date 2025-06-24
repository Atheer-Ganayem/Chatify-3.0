import {Loader2Icon } from "lucide-react";
import React from "react";

const ChatLoader = () => {
  return (
    <div className="h-full w-full flex justify-center items-center flex-col gap-5">
      <Loader2Icon className="animate-spin" size={100} />
      <p className="text-3xl font-semibold text-center">
        Loading messages, <br />
        please wait...
      </p>
    </div>
  );
};

export default ChatLoader;
