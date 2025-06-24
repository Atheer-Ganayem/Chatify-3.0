import { Avatar } from "@/components/ui/avatar";
import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { AvatarFallback, AvatarImage } from "@radix-ui/react-avatar";
import { useSession } from "next-auth/react";
import React from "react";

interface Props {
  text: string;
}

const WaitingMessageCard: React.FC<Props> = ({ text }) => {
  const { data } = useSession();
  if (!data) {
    return;
  }

  return (
    <div className="flex items-start max-w-2xl ms-auto flex-row-reverse">
      <Avatar className="h-12 w-12 rounded-full mx-2">
        <AvatarImage
          src={`${process.env.AWS}${data.user.avatar}`}
          alt={data.user.name}
        />
        <AvatarFallback className="rounded-lg">
          {data.user.name
            .split(" ")
            .map((n) => n[0])
            .join("")}
        </AvatarFallback>
      </Avatar>

      <Card className="flex-1 shadow-sm my-2 bg-primary opacity-50">
        <CardHeader className="pb-2">
          <div className="flex items-center justify-between">
            <h4 className="text-lg font-bold">{data.user.name}</h4>
            <span className="text-sm">Sending...</span>
          </div>
        </CardHeader>
        <CardContent className="pt-0">
          <p className="text-md leading-relaxed">{text}</p>
        </CardContent>
      </Card>
    </div>
  );
};

export default WaitingMessageCard;
