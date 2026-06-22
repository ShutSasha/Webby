import Link from 'next/link'

import { getUserColor } from '@/lib/utils/chat.utils'

type Props = {
  senderId: string
  username: string
  message: string
}

export default function RoomChatMessage({ senderId, username, message }: Props) {
  const userColor = getUserColor(senderId)

  return (
    <div>
      <Link href={`/profile/${senderId}`} className={`hover:underline font-medium ${userColor}`}>
        {username}
      </Link>
      <span className="text-neutral-300">:&nbsp;{message}</span>
    </div>
  )
}

export function RoomChatMessageSkeleton({ index = 0 }: { index?: number }) {
  const messageWidths = ['w-[60%]', 'w-[80%]', 'w-[40%]', 'w-[90%]', 'w-[50%]']
  const nameWidths = ['w-16', 'w-24', 'w-20', 'w-14', 'w-28']

  const messageWidth = messageWidths[index % messageWidths.length]
  const nameWidth = nameWidths[index % nameWidths.length]

  return (
    <div className="flex items-center gap-2 animate-pulse py-0.5 w-full">
      <div className={`h-3.5 ${nameWidth} bg-neutral-800/80 rounded shrink-0`} />
      <div className={`h-3.5 ${messageWidth} bg-neutral-800/50 rounded`} />
    </div>
  )
}
