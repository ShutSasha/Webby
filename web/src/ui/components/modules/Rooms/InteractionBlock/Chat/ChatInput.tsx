'use client'

import { useEffect, useRef, useState } from 'react'

import { useSession } from 'next-auth/react'

import SendIcon from '@/assets/icons/shared/send_message.svg'
import { useSendMessageMutation } from '@/lib/hooks/api/chat/useSendMessage'
import { useToastStore } from '@/stores/toast-store'

type Props = {
  chatId: string
}

export default function ChatInput({ chatId }: Props) {
  const { data: session } = useSession()
  const [content, setContent] = useState('')
  const inputRef = useRef<HTMLInputElement>(null)
  const addToast = useToastStore(state => state.addToast)

  const { mutate: sendMessage } = useSendMessageMutation()

  useEffect(() => {
    inputRef.current?.focus()
  }, [])

  const handleSend = (e?: React.FormEvent) => {
    e?.preventDefault()

    const trimmedContent = content.trim()
    if (!trimmedContent || !chatId || !session?.user) return

    setContent('')
    inputRef.current?.focus()

    sendMessage(
      {
        chatId,
        content: trimmedContent,
        sender: {
          id: session.user.id,
          username: session.user.username,
          avatarUrl: session.user.image || '',
        },
      },
      {
        onError: error => {
          setContent(trimmedContent)
          addToast(error.message || 'Failed to send message', 'error')
        },
      },
    )
  }

  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      handleSend()
    }
  }

  return (
    <form onSubmit={handleSend} className="relative">
      <input
        ref={inputRef}
        value={content}
        onChange={e => setContent(e.target.value)}
        onKeyDown={handleKeyDown}
        placeholder="Send a message"
        autoComplete="off"
        className="py-2 pl-4 pr-12 dark:bg-surface bg-background placeholder:text-foreground-ghost w-full ring-0
          outline-0 rounded-lg text-foreground-subtle transition-opacity"
      />
      <button
        type="submit"
        disabled={!content.trim()}
        className="absolute -translate-y-1/2 top-1/2 right-3 disabled:opacity-50"
      >
        <SendIcon
          className={`size-6 transition-colors duration-300 ease-in-out ${
            content.trim() ? 'text-emerald-500 cursor-pointer' : 'text-foreground-ghost'
          }`}
        />
      </button>
    </form>
  )
}
