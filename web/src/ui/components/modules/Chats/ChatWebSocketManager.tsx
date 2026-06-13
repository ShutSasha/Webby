'use client'

import { useChatWebSocket } from '@/lib/hooks/api/chat/useChatWebSocket'

type Props = {
  chatId: string
}

export default function ChatWebSocketManager({ chatId }: Props) {
  useChatWebSocket(chatId)
  return null
}
