type Conversation = {
  _id: string;
  participant: Participant;
  lastMessage?: Message;
};

type Participant = {
  _id: string;
  name: string;
  email: string;
  avatar: string;
};

type Message = {
  _id: string;
  sender: string;
  conversationId: string;
  text: string;
  createdAt: Date;
};

type WaitingMessage = {
  requestId: string;
  text: string;
  conversationID: string;
  createdAt: Date;
};

type WSRequest = {
  type: "msg";
  id: string;
  message: string;
  conversationId: string;
};

type WSResponse = {
  id: string;
  message: Message | string;
  type: "err" | "msg" | "acknowledged" | "delete" | "status" | "cnv";
  messageId?: string;
  userId?: string;
  online?: boolean;
  cnvId?: string;
  user?: Participant;
};
