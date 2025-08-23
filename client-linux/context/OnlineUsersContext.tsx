"use client";

import React, { createContext, useContext, useState } from "react";

type OnlineUsersContextType = {
  online: string[];
  addOnline: (...ids: string[]) => void;
  removeOnline: (id: string) => void;
  isOnline: (id: string) => boolean;
};

const OnlineUsersContext = createContext<OnlineUsersContextType | undefined>(
  undefined
);

const OnlineUsersProvider = ({ children }: { children: React.ReactNode }) => {
  const [online, setOnline] = useState<string[]>([]);

  function addOnline(...ids: string[]) {
    setOnline((prev) => {
      const newIds = ids.filter((id) => !prev.includes(id));
      return [...prev, ...newIds];
    });
  }

  function removeOnline(id: string) {
    setOnline((prev) => prev.filter((currentId) => currentId !== id));
  }

  function isOnline(id: string): boolean {
    return online.includes(id);
  }

  return (
    <OnlineUsersContext.Provider
      value={{
        online,
        addOnline,
        removeOnline,
        isOnline,
      }}
    >
      {children}
    </OnlineUsersContext.Provider>
  );
};

export const useOnlineUsers = () => {
  const context = useContext(OnlineUsersContext);
  if (!context) throw new Error("useData must be used within a DataProvider");
  return context;
};

export default OnlineUsersProvider;
