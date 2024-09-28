import { SenderTypeEnum } from '@trio/types';
import { GiArtificialHive } from 'react-icons/gi';
interface Props {
  senderType: SenderTypeEnum;
  messageContent: string;
  senderName?: string;
}

export function Message({ senderType, messageContent, senderName }: Props) {
  return (
    <div
      className={`${
        senderType === SenderTypeEnum.USER
          ? 'text-white text-xs text-left bg-neutral-800 self-end' // Gradient for User
          : 'self-start text-white'
      } text-black rounded-xl p-2 max-w-full sm:max-w-[80%] shadow-md`}
    >
      <p
        className={`${
          senderType === SenderTypeEnum.USER
            ? '' // Gradient for User
            : 'text-xs flex gap-1 items-center'
        } text-xs font-semibold mb-1 text-gray-200`}
      >
        {senderType === SenderTypeEnum.AGENT && (
          <GiArtificialHive className="text-base" />
        )}
        {senderName}
      </p>
      <p className="text-sm text-gray-100">{messageContent}</p>
    </div>
  );
}
