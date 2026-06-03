'use client'

import { useRoomWebSocket } from '@/lib/hooks/useRoomWebSocket'

type Props = {
  roomId: string | undefined
  chatId: string
}

export default function RoomWebSocketManager({ roomId, chatId }: Props) {
  useRoomWebSocket(roomId, chatId)

  return null
}
