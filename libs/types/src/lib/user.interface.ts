export interface User {
  id: string;
  fullName: string;
  userName: string;
  chats: Chat[];
}

// Additional interfaces to support the User interface
export interface Chat {
  id: string;
  chatName: string;
  messages: Message[];
  agents: Agent[];
}

export interface Message {
  id: string;
  content: string;
  senderType: string;
  senderId: number;
}

export interface Agent {
  id: string;
  name: string;
  lingo: string;
  traits: string[];
}
