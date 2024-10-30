'use client';

import { CiSettings } from 'react-icons/ci';
import { Chat } from '@trio/types';
import { v4 } from 'uuid';
import { Message } from './Message';
import { UpdateReflectionModal } from './UpdateReflectionModal';
import { useReflectionChat } from '@trio/hooks';
import { PlaceholdersAndVanishInput } from '@/shadcn/ui/placeholders-and-vanish-input';
import { useState } from 'react';

const placeholders = [
  'How many Rs are in the word strawberry',
  'What should I do if my best friend is ignoring me?',
  'How do I get over a crush?',
  "What's the best way to handle rumors about me?",
  'How do I tell my parents I got a bad grade?',
];

export function ReflectionChatDisplay({ chat }: { chat: Chat }) {
  const { sendMessage } = useReflectionChat(chat.id);

  const [message, setMessage] = useState('');

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setMessage(e.target.value);
  };
  const onSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    sendMessage({
      chatId: chat.id,
      content: message,
      onData: (data) => console.log(JSON.parse(data)),
      onError: (error) => console.error(error),
      onComplete: () => console.log('complete'),
    });
  };

  return (
    <>
      <div className="overflow-y-scroll w-full h-full flex flex-col gap-5 pb-32">
        <div className="text-white fixed font-semibold w-full text-sm sm:text-lg">
          <UpdateReflectionModal
            chat={chat}
            trigger={
              <button
                type="button"
                className="flex items-center gap-1 bg-transparent"
              >
                {chat.chatName}
                <CiSettings />
              </button>
            }
          />
        </div>
        {chat.messages?.length > 0 &&
          chat.messages.map((message, index) => {
            const isSameSender =
              index !== 0 &&
              chat.messages[index - 1].sender.username ===
                message.sender.username &&
              chat.messages[index - 1].sender.name === message.sender.name;

            return (
              <Message
                key={v4()}
                senderType={message.senderType}
                messageContent={message.content}
                senderName={
                  isSameSender
                    ? undefined
                    : message.sender.username ?? message.sender.name
                }
              />
            );
          })}
      </div>
      <PlaceholdersAndVanishInput
        placeholders={placeholders}
        onChange={handleChange}
        onSubmit={onSubmit}
      />
    </>
  );
}
