"use client";

import FilterConversatonBox from "./FilterConversatonBox";
import {
  Sidebar,
  SidebarContent,
  SidebarHeader,
  SidebarRail,
} from "@/components/ui/sidebar";
import Image from "next/image";
import SidebarFooter from "./SidebarFooter";
import ConversationsList from "./ConversationsList";
import { useState } from "react";

export function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {
  const [filter, setFilter] = useState<string>("");

  function onFilter(str: string) {
    setFilter(str);
  }

  return (
    <Sidebar {...props}>
      <SidebarHeader>
        <div className="font-bold text-3xl flex justify-center items-center gap-2 py-2">
          <Image src={"/logo.png"} width={40} height={40} alt="logo" />
          Chatify
        </div>
        <FilterConversatonBox currentFilter={filter} onFilter={onFilter} />
      </SidebarHeader>
      <SidebarContent>
        <ConversationsList filter={filter.toLowerCase()} />
      </SidebarContent>
      <SidebarFooter />
      <SidebarRail />
    </Sidebar>
  );
}
