'use client'

import { useSession } from 'next-auth/react'

import { useRoomWebSocket } from '@/lib/hooks/useRoomWebSocket'

type Props = {
  roomId: string | undefined
  chatId: string
}

export default function RoomWebSocketManager({ roomId, chatId }: Props) {
  const { data: session } = useSession()
  const userId = session?.user?.id
  useRoomWebSocket(roomId, chatId, userId)

  return null
}
