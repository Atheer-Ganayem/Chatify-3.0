"use client";

import { useSearchParams } from "next/navigation";
import { cn } from "@/lib/utils";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import SignupForm from "./SignupForm";
import LoginForm from "./LoginForm";

const AuthCard = () => {
  const params = useSearchParams();
  const mode = params.get("mode") === "register" ? "Register" : "Login";

  return (
    <Card className={cn("w-[380px]")}>
      <CardHeader>
        <CardTitle className="text-3xl text-center">{mode}</CardTitle>
      </CardHeader>
      <CardContent className="grid gap-4">
        {mode == "Login" ? <LoginForm /> : <SignupForm />}
      </CardContent>
    </Card>
  );
};

export default AuthCard;
