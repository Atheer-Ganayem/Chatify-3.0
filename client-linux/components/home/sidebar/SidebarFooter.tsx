"use client";

import React from "react";
import {
  SidebarFooter as SidebarFooterCN,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "../../ui/dropdown-menu";
import { Avatar, AvatarFallback, AvatarImage } from "../../ui/avatar";
import {
  BellOff,
  BellRing,
  ChevronUp,
  LogOut,
  Moon,
  Sun,
  Volume2,
  VolumeOff,
} from "lucide-react";
import { signOut, useSession } from "next-auth/react";
import { useTheme } from "next-themes";
import { useNotification } from "@/context/NotificationContext";
import AccountDialog from "./account-dialog/AccountDialog";

const SidebarFooter = () => {
  const session = useSession();

  return (
    <SidebarFooterCN>
      <SidebarMenu>
        <SidebarMenuItem>
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <SidebarMenuButton className="h-auto data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground">
                <Avatar className="h-16 w-16 rounded-full">
                  <AvatarImage
                    src={`${process.env.AWS}${session.data?.user.avatar}`}
                    alt={session.data?.user.name}
                  />
                  <AvatarFallback className="rounded-lg">
                    {session.data?.user.name
                      .split(" ")
                      .map((n) => n[0])
                      .join("")}
                  </AvatarFallback>
                </Avatar>
                <div className="grid flex-1 text-left text-sm leading-tight">
                  <span className="truncate font-semibold">
                    {session.data?.user.name}
                  </span>
                  <span className="truncate text-xs">
                    {session.data?.user.email}
                  </span>
                </div>
                <ChevronUp className="ml-auto size-4" />
              </SidebarMenuButton>
            </DropdownMenuTrigger>
            <DropdownMenuContent
              className="w-[--radix-dropdown-menu-trigger-width] min-w-56 rounded-lg"
              side="bottom"
              align="end"
              sideOffset={4}
            >
              <AccountDialog />
              <NotificationItems />
              <ThemeItem />
              <DropdownMenuSeparator />
              <LgoutItem />
            </DropdownMenuContent>
          </DropdownMenu>
        </SidebarMenuItem>
      </SidebarMenu>
    </SidebarFooterCN>
  );
};

function ThemeItem() {
  const { setTheme, theme } = useTheme();

  return (
    <DropdownMenuItem
      onClick={() => setTheme(theme === "light" ? "dark" : "light")}
    >
      {theme === "dark" ? (
        <>
          <Moon className="mr-2 h-4 w-4" />
          <span>Dark Mode</span>
        </>
      ) : (
        <>
          <Sun className="mr-2 h-4 w-4" />
          <span>Light mode</span>
        </>
      )}
    </DropdownMenuItem>
  );
}

function LgoutItem() {
  return (
    <DropdownMenuItem
      variant="destructive"
      onClick={() =>
        signOut({ redirect: true, callbackUrl: "/auth?mode=login" })
      }
    >
      <LogOut className="mr-2 h-4 w-4" />
      <span>Log out</span>
    </DropdownMenuItem>
  );
}

function NotificationItems() {
  const notification = useNotification();

  return (
    <>
      <DropdownMenuItem
        onClick={() => notification.setValues({ active: !notification.active })}
      >
        {notification.active ? <BellRing /> : <BellOff />}
        <span>Notifications: {notification.active ? "On" : "Off"}</span>
      </DropdownMenuItem>
      <DropdownMenuItem
        onClick={() =>
          notification.setValues({ sound: !notification.soundActive })
        }
      >
        {notification.soundActive ? <Volume2 /> : <VolumeOff />}
        <span>
          Notification sound: {notification.soundActive ? "On" : "Off"}
        </span>
      </DropdownMenuItem>
    </>
  );
}

export default SidebarFooter;
