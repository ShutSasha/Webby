'use client'

import { useState } from 'react'

import SendIcon from '@/assets/icons/shared/send_message.svg'
import { useSendMessageMutation } from '@/lib/hooks/api/chat/useSendMessage'

type Props = {
  chatId: string
}

export default function ChatMessageInput({ chatId }: Props) {
  const [messageText, setMessageText] = useState('')
  const { mutate: sendMessage, isPending: isSending } = useSendMessageMutation()

  const handleSendMessage = () => {
    const trimmedMessage = messageText.trim()

    if (!trimmedMessage || isSending) return

    sendMessage(
      { chatId, content: trimmedMessage },
      {
        onSuccess: () => {
          setMessageText('')
        },
      },
    )
  }

  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter') {
      e.preventDefault()
      handleSendMessage()
    }
  }

  return (
    <div className="p-4 bg-neutral-900/20 border-t border-neutral-800/60 shrink-0">
      <div className="relative flex items-center">
        <input
          type="text"
          value={messageText}
          onChange={e => setMessageText(e.target.value)}
          onKeyDown={handleKeyDown}
          placeholder="Write a message..."
          disabled={isSending}
          className="w-full bg-[#141414] border border-neutral-800 text-neutral-200 placeholder:text-neutral-600
            rounded-xl py-3.5 pl-5 pr-12 outline-none focus:border-neutral-600 transition-colors disabled:opacity-50"
        />
        <button
          onClick={handleSendMessage}
          disabled={!messageText.trim() || isSending}
          className="absolute right-3 p-1.5 text-neutral-500 hover:text-emerald-500 transition-colors
            disabled:hover:text-neutral-500 disabled:opacity-50"
        >
          <SendIcon className="size-5" />
        </button>
      </div>
    </div>
  )
}
