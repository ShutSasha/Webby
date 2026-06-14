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
