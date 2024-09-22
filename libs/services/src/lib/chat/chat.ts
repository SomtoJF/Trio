import { BaseRoute, Route } from '../routes';

interface Data {
  chatName: string;
  agents: Agent[];
}
interface Agent {
  name: string;
  lingo: string;
  traits: string[];
}

interface Chat {
  id: string;
  chatName: string;
  messages: any[];
  agents: any[];
}

export async function createChatWithAgents(data: Data) {
  const res = await fetch(`${BaseRoute}/${Route.Chats.CreateWithAgents}`, {
    method: 'POST',
    credentials: 'include',
    body: JSON.stringify(data),
  });
  if (res.status > 299) throw new Error(res.statusText);
  const result = await res.json();
  return result;
}

export async function currentUserChats(): Promise<Chat[]> {
  const res = await fetch(`${BaseRoute}/${Route.CurrentUser.Chats}`, {
    method: 'GET',
    credentials: 'include',
  });
  if (res.status > 299) throw new Error(res.statusText);
  const result = await res.json();
  return result.data as Chat[];
}
