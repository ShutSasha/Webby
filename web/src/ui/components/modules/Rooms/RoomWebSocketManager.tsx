'use client'

import { useSession } from 'next-auth/react'

import { useChatSocketListeners } from '@/lib/hooks/api/room/ws/useChatSocketListeners'
import { useEntertainmentRoomListeners } from '@/lib/hooks/api/room/ws/useEntertainmentRoomListeners'
import { usePlayerSocketListeners } from '@/lib/hooks/api/room/ws/usePlayerSocketListeners'
import { useRoomMemberSocketListeners } from '@/lib/hooks/api/room/ws/useRoomMemberSocketListeners'
import { useRoomSocketConnection } from '@/lib/hooks/api/room/ws/useRoomSocketConnection'
import { useVoteSocketListeners } from '@/lib/hooks/api/room/ws/useVoteSocketListeners'

type Props = {
  roomId: string | undefined
  chatId: string
}

export default function RoomWebSocketManager({ roomId, chatId }: Props) {
  const { data: session } = useSession()
  const userId = session?.user?.id
  useRoomSocketConnection(roomId, chatId)
  useChatSocketListeners(chatId)
  useEntertainmentRoomListeners(userId, roomId)
  usePlayerSocketListeners(roomId)
  useVoteSocketListeners(roomId)
  useRoomMemberSocketListeners(roomId, userId)

  return null
}
