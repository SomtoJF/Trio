import { Message } from './message.interface';
import { Agent } from './agent.interface';

export interface Chat {
  id: string; // UUID
  chatName: string;
  messages: Message[];
  agents: Agent[];
}
