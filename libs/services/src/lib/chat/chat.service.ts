import { Agent, Chat } from '@trio/types';
import { BaseRoute, Route } from '../routes';

interface Data {
  chatName: string;
  agents: Partial<Agent>[];
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
  const res = await fetch(`${BaseRoute}/chats`, {
    method: 'GET',
    credentials: 'include',
  });
  if (res.status > 299) throw new Error(res.statusText);
  const result = await res.json();
  return result.data as Chat[];
}

export async function getOneChat(id: string): Promise<Chat> {
  const res = await fetch(`${BaseRoute}/${Route.Chats.Default}${id}`, {
    method: 'GET',
    credentials: 'include',
  });
  if (res.status > 299) throw new Error(res.statusText);
  const result = await res.json();
  return result.data as Chat;
}

export async function addMessageToChat(chatId: string, message: string) {
  const data = { content: message };
  const res = await fetch(
    `${BaseRoute}/${Route.Chats.Default}${chatId}/messages`,
    {
      method: 'POST',
      credentials: 'include',
      body: JSON.stringify(data),
    }
  );
  const result = await res.json();
  if (res.status > 299) throw new Error(result.error ?? res.statusText);
  return result;
}

export async function updateChat(
  chatId: string,
  data: {
    chatName: string;
    agents: { name: string; lingo: string; traits: string[] }[];
  }
) {
  const res = await fetch(`${BaseRoute}/${Route.Chats.Default}${chatId}`, {
    method: 'PUT',
    credentials: 'include',
    body: JSON.stringify(data),
  });

  const result = await res.json();
  if (res.status > 299) throw new Error(result.error ?? res.statusText);
  return result.data;
}
