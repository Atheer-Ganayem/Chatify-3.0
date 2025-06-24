import React from "react";
import { SidebarMenuButton } from "../ui/sidebar";
import { Avatar } from "../ui/avatar";

const ConversationListLoader = () => {
  return Array.from({ length: 5 }).map((_, index) => (
    <SidebarMenuButton
      key={index}
      asChild
      className="h-auto hover:bg-transparent animate-pulse"
    >
      <span className="flex gap-4">
        <Avatar className="h-14 w-14 rounded-full bg-muted" />
        <div className="w-full flex flex-col gap-5">
          <h4 className="w-[12rem] bg-muted h-4" />
          <p className="w-[18rem] bg-muted h-3" />
        </div>
      </span>
    </SidebarMenuButton>
  ));
};

export default ConversationListLoader;
