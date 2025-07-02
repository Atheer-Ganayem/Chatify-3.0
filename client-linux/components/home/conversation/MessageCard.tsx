import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { useConversations } from "@/context/ConversationsContext";
import { getDate } from "@/utils/date-formatting";
import { useSession } from "next-auth/react";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { ClipboardCopy, Loader2, Trash2 } from "lucide-react";
import { toast } from "sonner";
import { useState } from "react";
import { useMessages } from "@/context/MessagesContext";
import { getCookie } from "cookies-next";

interface Props {
  message: Message;
}

const MessageCard: React.FC<Props> = ({ message }) => {
  const session = useSession();
  const conversationCtx = useConversations();
  const msgCtx = useMessages();
  const [isDeleting, setIsDeleting] = useState<boolean>(false);

  if (!session || !session.data) {
    return;
  }
  const sender =
    message.sender === session.data.user.id
      ? {
          name: session.data.user.name,
          avatar: session.data.user.avatar,
        }
      : {
          name: conversationCtx.currentConversation?.participant.name,
          avatar: conversationCtx.currentConversation?.participant.avatar,
        };

  async function copyHandler() {
    await navigator.clipboard.writeText(message.text);
    toast.success("Message copied to clipboard.");
  }

  async function deleteHandler() {
    try {
      setIsDeleting(true);
      const response = await fetch(
        `${process.env.BACKEND_URL}/message/${message._id}`,
        {
          method: "DELETE",
          headers: {
            Authorization: `Bearer ${getCookie("next-auth.session-token")}`,
          },
          credentials: "include",
        }
      );
      const responseData = await response.json();
      if (!response.ok) {
        toast.error(responseData.message);
      }
      toast.success("Message deleted successfully.");
      msgCtx.onDelete(message._id);
    } catch (error) {
      console.log(error);
      toast.error("Coudln't delete message, please try again later.");
    } finally {
      setIsDeleting(false);
    }
  }

  return (
    <div
      className={`flex items-start max-w-2xl ${
        message.sender === session.data.user.id && "ms-auto flex-row-reverse"
      }`}
    >
      <Avatar className="h-12 w-12 rounded-full mx-2">
        <AvatarImage
          src={`${process.env.AWS}${sender.avatar}`}
          alt={sender.name}
        />
        <AvatarFallback className="rounded-lg">
          {sender
            .name!.split(" ")
            .map((n) => n[0])
            .join("")}
        </AvatarFallback>
      </Avatar>

      <Card
        className={`flex-1 shadow-sm my-2 ${
          message.sender === session.data.user.id &&
          (isDeleting ? "bg-red-800 opacity-60" : "bg-primary")
        }`}
      >
        <CardHeader>
          <div className="flex items-center justify-between">
            <h4 className="text-lg font-bold">{sender.name}</h4>
            <span className="text-sm flex gap-3 items-center">
              {getDate(new Date(message.createdAt))}
              <DropdownMenu>
                <DropdownMenuTrigger className="flex flex-col items-center justify-center w-6 h-6 rounded">
                  <span className="w-1 h-1 rounded-full mb-0.5 bg-foreground"></span>
                  <span className="w-1 h-1 rounded-full mb-0.5 bg-foreground"></span>
                  <span className="w-1 h-1 rounded-full bg-foreground"></span>
                </DropdownMenuTrigger>
                <DropdownMenuContent>
                  <DropdownMenuItem onClick={async () => await copyHandler()}>
                    <ClipboardCopy /> Copy
                  </DropdownMenuItem>
                  {message.sender === session.data.user.id && (
                    <DropdownMenuItem
                      disabled={isDeleting}
                      variant="destructive"
                      onClick={deleteHandler}
                    >
                      {isDeleting ? (
                        <>
                          <Loader2 className="animate-spin" /> Deleting...
                        </>
                      ) : (
                        <>
                          <Trash2 /> Delete
                        </>
                      )}
                    </DropdownMenuItem>
                  )}
                </DropdownMenuContent>
              </DropdownMenu>
            </span>
          </div>
        </CardHeader>
        <CardContent className="pt-0">
          <p className="text-md leading-relaxed">{message.text}</p>
        </CardContent>
      </Card>
    </div>
  );
};

export default MessageCard;
