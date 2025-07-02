"use client";

import { Button } from "@/components/ui/button";
import { getCookie } from "cookies-next";

const page = () => {
  async function ping() {
    try {
      const response = await fetch(`${process.env.BACKEND_URL}/ping`, {
        headers: {
          Authorization: `Bearer ${getCookie("next-auth.session-token")}`,
        },
      });
      const responseData = await response.json();
      console.log(response, response);

      console.log(responseData);
    } catch (error) {
      console.log(error);
    }
  }

  return (
    <div>
      <Button onClick={async () => await ping()}>Ping</Button>
    </div>
  );
};

export default page;
