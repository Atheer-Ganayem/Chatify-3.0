import { useConversations } from "@/context/ConversationsContext";
import React, { useState } from "react";
import { Avatar, AvatarFallback, AvatarImage } from "../../ui/avatar";
import { Button } from "../../ui/button";
import { Loader2Icon, Mail } from "lucide-react";
import { redirect, useRouter } from "next/navigation";
import { useSession } from "next-auth/react";
import { createConversation } from "@/utils/requests";
import { toast } from "sonner";

interface Props {
  user: Participant;
  onClose: () => void;
}

const AddConversationResults: React.FC<Props> = ({ user, onClose }) => {
  const ctx = useConversations();
  const { data } = useSession();
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const router = useRouter();

  async function onClickHandler() {
    ctx?.conversations?.forEach((cnv) => {
      if (cnv.participant._id === user._id) {
        onClose();
        redirect(`/?conversationID=${cnv._id}`);
      }
    });
    try {
      setIsLoading(true);

      const response = await createConversation(user._id);
      if (!response.ok) {
        toast.error(response.message);
      }

      ctx.appendConversation({
        _id: response.conversationID,
        participant: user,
      });
      onClose();
      router.push(`/?conversationID=${response.conversationID}`);
    } catch (error) {
      console.log(error);
      toast.error("Something went wrong, please try again later.");
    } finally {
      setIsLoading(false);
    }
  }

  return (
    <div className="flex justify-between items-center">
      <Avatar className="h-14 w-14 rounded-full">
        <AvatarImage src={`${process.env.AWS}${user.avatar}`} alt={user.name} />
        <AvatarFallback className="rounded-lg">
          {user.name
            .split(" ")
            .map((n) => n[0])
            .join("")}
        </AvatarFallback>
      </Avatar>
      <p>{user.name}</p>
      {!isLoading ? (
        <Button
          onClick={onClickHandler}
          type="button"
          disabled={user._id == data?.user.id}
        >
          <Mail />
        </Button>
      ) : (
        <Button size="sm" disabled>
          <Loader2Icon className="animate-spin" />
        </Button>
      )}
    </div>
  );
};

export default AddConversationResults;
