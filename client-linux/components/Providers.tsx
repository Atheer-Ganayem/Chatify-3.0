"use client";

import React from "react";
import { ThemeProvider } from "next-themes";
import { SessionProvider } from "next-auth/react";
import ConversationsProvider from "@/context/ConversationsContext";
import MessagesProvider from "@/context/MessagesContext";
import NotificationProvider from "@/context/NotificationContext";
import OnlineUsersProvider from "@/context/OnlineUsersContext";
import WebsocketProvider from "@/context/WebsocketContext";

const Providers = ({ children }: { children: React.ReactNode }) => {
  return (
    <ThemeProvider
      attribute="class"
      defaultTheme="system"
      enableSystem
      disableTransitionOnChange
    >
      <SessionProvider>
        <NotificationProvider>
          <OnlineUsersProvider>
            <ConversationsProvider>
              <MessagesProvider>
                <WebsocketProvider>{children}</WebsocketProvider>
              </MessagesProvider>
            </ConversationsProvider>
          </OnlineUsersProvider>
        </NotificationProvider>
      </SessionProvider>
    </ThemeProvider>
  );
};

export default Providers;
