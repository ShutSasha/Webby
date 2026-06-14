'use client'

import { useEffect, useRef, useState } from 'react'

import SendIcon from '@/assets/icons/shared/send_message.svg'
import { useSendMessageMutation } from '@/lib/hooks/api/chat/useSendMessage'
import { useToastStore } from '@/stores/toast-store'

type Props = {
  chatId: string
}

export default function ChatInput({ chatId }: Props) {
  const [content, setContent] = useState('')
  const inputRef = useRef<HTMLInputElement>(null)
  const addToast = useToastStore(state => state.addToast)

  const { mutate: sendMessage, isPending } = useSendMessageMutation()

  const prevPendingRef = useRef(isPending)

  useEffect(() => {
    if (prevPendingRef.current && !isPending) {
      inputRef.current?.focus()
    }
    prevPendingRef.current = isPending
  }, [isPending])

  const handleSend = (e?: React.FormEvent) => {
    e?.preventDefault()

    if (!content.trim() || isPending || !chatId) return

    sendMessage(
      { chatId, content: content.trim() },
      {
        onSuccess: () => {
          setContent('')
        },
        onError: error => {
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
        disabled={isPending}
        placeholder="Send a message"
        className="py-2 pl-4 pr-12 bg-neutral-900 placeholder:text-neutral-700 w-full ring-0 outline-0 rounded-lg
          text-neutral-300 disabled:opacity-50"
      />
      <button
        type="submit"
        disabled={!content.trim() || isPending}
        className="absolute -translate-y-1/2 top-1/2 right-3 disabled:opacity-50"
      >
        <SendIcon
          className={`size-6 transition-colors duration-300 ease-in-out ${
            content.trim() ? 'text-emerald-500 cursor-pointer' : 'text-neutral-700'
          }`}
        />
      </button>
    </form>
  )
}
