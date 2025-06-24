"use client";

import React from "react";
import { ThemeProvider } from "next-themes";
import { SessionProvider } from "next-auth/react";
import ConversationsProvider from "@/context/ConversationsContext";
import MessagesProvider from "@/context/MessagesContext";
import NotificationProvider from "@/context/NotificationContext";

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
          <ConversationsProvider>
            <MessagesProvider>{children}</MessagesProvider>
          </ConversationsProvider>
        </NotificationProvider>
      </SessionProvider>
    </ThemeProvider>
  );
};

export default Providers;
