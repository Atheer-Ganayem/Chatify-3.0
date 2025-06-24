import AuthCard from "@/components/auth/AuthCard";
import { getServerSession } from "next-auth";
import React from "react";
import { redirect } from "next/navigation";
import Header from "@/components/auth/Header";
import authOptions from "@/utils/authOptions";

export const metadata = {
  title: "Chatify | Auth",
};

const page = async () => {
  const session = await getServerSession(authOptions);
  if (session) {
    redirect("/");
  }

  return (
    <>
      <Header />
      <div className="p-10 flex justify-center items-center">
        <AuthCard />
      </div>
    </>
  );
};

export default page;
