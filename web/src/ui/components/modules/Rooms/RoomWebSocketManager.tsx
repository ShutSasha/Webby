'use client'

import { useRoomWebSocket } from '@/lib/hooks/useRoomWebSocket'

type Props = {
  chatId: string
}

export default function RoomWebSocketManager({ chatId }: Props) {
  useRoomWebSocket(chatId)

  return null
}
