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
