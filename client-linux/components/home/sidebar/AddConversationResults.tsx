import { useConversations } from "@/context/ConversationsContext";
import { Avatar, AvatarFallback, AvatarImage } from "../../ui/avatar";
import { Button } from "../../ui/button";
import { Loader2Icon, Mail } from "lucide-react";
import { redirect, useRouter } from "next/navigation";
import { useSession } from "next-auth/react";
import { toast } from "sonner";
import { useOnlineUsers } from "@/context/OnlineUsersContext";
import useFetch from "@/hooks/useFetch";

interface Props {
  user: Participant;
  onClose: () => void;
}

const AddConversationResults: React.FC<Props> = ({ user, onClose }) => {
  const ctx = useConversations();
  const { data } = useSession();
  const router = useRouter();
  const onlineCtx = useOnlineUsers();
  const { exec, isLoading } = useFetch({
    path: "/conversation",
    method: "POST",
    auth: true,
  });

  async function onClickHandler() {
    ctx?.conversations?.forEach((cnv) => {
      if (cnv.participant._id === user._id) {
        onClose();
        redirect(`/?conversationID=${cnv._id}`);
      }
    });
    try {
      const { ok, responseData, error } = await exec(
        JSON.stringify({ targetUserID: user._id })
      );
      if (error) throw error;

      if (ok) {
        toast.error(responseData.message);
      }

      ctx.appendConversation({
        _id: responseData.conversationID,
        participant: user,
      });
      if (responseData.isOnline) {
        onlineCtx.addOnline(user._id);
      }
      onClose();
      router.push(`/?conversationID=${responseData.conversationID}`);
    } catch (error) {
      console.log(error);
      toast.error("Something went wrong, please try again later.");
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
          disabled={
            user._id == data?.user.id ||
            !!ctx.conversations.find((c) => c.participant._id === user._id)
          }
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
