"use client";

import { useConversations } from "@/context/ConversationsContext";
import React from "react";
import { SidebarMenuButton } from "../../ui/sidebar";
import Link from "next/link";
import { Avatar, AvatarFallback, AvatarImage } from "../../ui/avatar";
import { getDate } from "@/utils/date-formatting";
import { useSession } from "next-auth/react";
import ConversationListLoader from "@/components/loaders/ConversationListLoader";

interface Props {
  filter: string;
}

const ConversationsList: React.FC<Props> = ({ filter }) => {
  const ctx = useConversations();
  const session = useSession();

  if (!ctx || ctx.loading) {
    return <ConversationListLoader />;
  }

  let conversations = filter
    ? ctx!.conversations!.filter((cnv) =>
        cnv.participant.name.toLowerCase().includes(filter)
      )
    : ctx!.conversations!;

  conversations = conversations.sort((a, b) => {
    if (!a.lastMessage && b.lastMessage) return 1;
    else if (a.lastMessage && !b.lastMessage) return -1;
    const dateA = new Date(a.lastMessage?.createdAt || Date.now()).getTime();
    const dateB = new Date(b.lastMessage?.createdAt || Date.now()).getTime();
    return dateB - dateA;
  });

  return conversations.map((cnv) => (
    <SidebarMenuButton
      key={cnv._id}
      asChild
      isActive={ctx.currentConversation?._id === cnv._id}
      className="h-auto"
    >
      <Link
        href={{ pathname: "/", query: { conversationID: cnv._id } }}
        className="flex gap-4"
      >
        <Avatar className="h-14 w-14 rounded-full">
          <AvatarImage
            src={`${process.env.AWS}${cnv.participant.avatar}`}
            alt={cnv.participant.name}
          />
          <AvatarFallback className="rounded-lg">
            {cnv.participant.name
              .split(" ")
              .map((n) => n[0])
              .join("")}
          </AvatarFallback>
        </Avatar>
        <div className="w-full">
          <h4 className="font-bold text-xl">{cnv.participant.name}</h4>
          {cnv.lastMessage && (
            <>
              {cnv.lastMessage.sender === cnv.participant._id
                ? cnv.participant.name
                : session?.data?.user.name}
              :{" "}
              <span className="text-primary-foreground opacity-75">
                {cnv.lastMessage.text.length > 30
                  ? cnv.lastMessage.text.slice(0, 30) + "..."
                  : cnv.lastMessage.text}
              </span>
              <div className="flex justify-end">
                <span>{getDate(new Date(cnv.lastMessage.createdAt))}</span>
              </div>
            </>
          )}
        </div>
      </Link>
    </SidebarMenuButton>
  ));
};

export default ConversationsList;
