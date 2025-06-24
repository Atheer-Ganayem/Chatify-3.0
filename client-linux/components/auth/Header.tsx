"use client";

import Image from "next/image";
import { Button } from "../ui/button";
import Link from "next/link";
import ThemeToggle from "../ui/ThemeToggle";

const Header = () => {
  return (
    <header className="border-b-2">
      <div className="max-w-7xl mx-auto py-3 px-5 flex justify-between items-center">
        <Link href="/" className="font-bold text-xl flex items-center gap-2">
          <Image src={"/logo.png"} width={40} height={40} alt="logo" />
          Chatify
        </Link>
        <div className="flex gap-5 items-center">
          <ThemeToggle />
          <Link href="/auth?mode=login">
            <Button>Login</Button>
          </Link>
        </div>
      </div>
    </header>
  );
};

export default Header;
