"use client";

import { getCookie } from "cookies-next";
import { useState } from "react";

type Args = {
  method?: "GET" | "POST" | "DELETE" | "PUT";
  path: string;
  auth: boolean;
  defaultLoadong?: boolean;
};

export default function useFetch(options: Args) {
  const [isLoading, setIsLoading] = useState<boolean>(!!options.defaultLoadong);

  const headers: Record<string, string> = {};
  if (options.auth) {
    headers.Authorization = `Bearer ${getCookie("next-auth.session-token")}`;
  }

  async function exec(body?: FormData | string) {
    console.log("exec...");

    if (typeof body === "string") {
      headers["Content-Type"] = "application/json";
    }

    try {
      setIsLoading(true);
      const response = await fetch(
        `${process.env.BACKEND_URL}${options.path}`,
        {
          method: options.method || "GET",
          headers: headers,
          body: body || null,
        }
      );
      const responseData = await response.json();

      return { ok: response.ok, responseData };
    } catch (error) {
      return { ok: false, responseData: null, error };
    } finally {
      setIsLoading(false);
    }
  }

  return { isLoading, exec };
}
