"use client";

import { getCookie } from "cookies-next";

interface createConversationReturnValue {
  ok: boolean;
  message: string;
  conversationID: string;
}

export async function createConversation(
  targetUserID: string
): Promise<createConversationReturnValue> {
  const response = await fetch(`${process.env.BACKEND_URL}/conversation`, {
    method: "POST",
    headers: {
      "Content-Type": "aplication/json",
      Authorization: `Bearer ${getCookie("next-auth.session-token")}`,
    },
    credentials: "include",
    body: JSON.stringify({
      targetUserID,
    }),
  });
  const responseData = await response.json();

  return {
    ok: response.ok,
    message: responseData.message,
    conversationID: responseData.conversationID,
  };
}
