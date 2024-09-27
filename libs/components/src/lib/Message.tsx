import { SenderTypeEnum } from '@trio/types';

interface Props {
  senderType: SenderTypeEnum;
  messageContent: string;
  senderName: string;
}

export function Message({ senderType, messageContent, senderName }: Props) {
  return (
    <div
      className={`${
        senderType === SenderTypeEnum.USER
          ? 'bg-green-400 self-end' // Gradient for User
          : 'self-start bg-gray-100'
      } text-black rounded-md px-4 py-2 max-w-[80%] shadow-md`}
    >
      <p className="text-xs font-semibold uppercase text-gray-600">
        {senderType === SenderTypeEnum.USER ? 'Me' : senderName}
      </p>
      <p className="text-sm">{messageContent}</p>
    </div>
  );
}
