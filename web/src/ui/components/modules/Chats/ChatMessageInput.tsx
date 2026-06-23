'use client'

import { useEffect, useRef, useState } from 'react'

import { useSession } from 'next-auth/react'

import CloseIcon from '@/assets/icons/shared/circle-xmark.svg'
import EditPenIcon from '@/assets/icons/shared/edit-pen.svg'
import SendIcon from '@/assets/icons/shared/send_message.svg'
import { useEditMessageMutation } from '@/lib/hooks/api/chat/useEditMessage'
import { useSendMessageMutation } from '@/lib/hooks/api/chat/useSendMessage'

type Props = {
  chatId: string
  editingMessage: { id: string; content: string } | null
  onCancelEdit: () => void
}

export default function ChatMessageInput({ chatId, editingMessage, onCancelEdit }: Props) {
  const { data: session } = useSession()
  const [messageText, setMessageText] = useState('')
  const inputRef = useRef<HTMLInputElement>(null)

  const { mutate: sendMessage, isPending: isSending } = useSendMessageMutation()
  const { mutate: editMessage, isPending: isEditing } = useEditMessageMutation()

  const isPending = isSending || isEditing

  const prevPendingRef = useRef(isPending)

  useEffect(() => {
    if (prevPendingRef.current && !isPending) {
      inputRef.current?.focus()
    }
    prevPendingRef.current = isPending
  }, [isPending])

  useEffect(() => {
    if (editingMessage) {
      setMessageText(editingMessage.content)
      inputRef.current?.focus()
    } else {
      setMessageText('')
    }
  }, [editingMessage])

  const handleSendMessage = () => {
    const trimmedMessage = messageText.trim()

    if (!trimmedMessage || isPending || !session?.user) return

    if (editingMessage) {
      editMessage(
        { chatId, messageId: editingMessage.id, content: trimmedMessage },
        {
          onSuccess: () => {
            onCancelEdit()
            setMessageText('')
          },
        },
      )
    } else {
      sendMessage(
        {
          chatId,
          content: trimmedMessage,
          sender: {
            id: session.user.id,
            username: session.user.username,
            avatarUrl: session.user.image || '',
          },
        },
        {
          onSuccess: () => {
            setMessageText('')
          },
        },
      )
    }
  }

  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter') {
      e.preventDefault()
      handleSendMessage()
    } else if (e.key === 'Escape' && editingMessage) {
      e.preventDefault()
      onCancelEdit()
    }
  }

  return (
    <div className="p-4 bg-surface/20 border-t border-border/60 shrink-0 flex flex-col transition-all duration-300">
      {editingMessage && (
        <div
          className="flex items-center justify-between mb-3 px-2 animate-in fade-in slide-in-from-bottom-2 duration-200"
        >
          <div className="flex items-center gap-3 overflow-hidden">
            <EditPenIcon className="size-5 text-emerald-500 shrink-0" />
            <div className="flex flex-col overflow-hidden border-l-2 border-emerald-500/50 pl-2">
              <span className="text-xs font-semibold text-emerald-500 leading-tight">Edit message</span>
              <span className="text-sm text-foreground-muted truncate max-w-[200px] sm:max-w-[400px]">
                {editingMessage.content}
              </span>
            </div>
          </div>
          <button
            onClick={onCancelEdit}
            className="p-1.5 hover:bg-background rounded-full text-foreground-faint hover:text-foreground-subtle
              transition-colors"
          >
            <CloseIcon className="size-5" />
          </button>
        </div>
      )}

      <div className="relative flex items-center">
        <input
          ref={inputRef}
          type="text"
          value={messageText}
          onChange={e => setMessageText(e.target.value)}
          onKeyDown={handleKeyDown}
          placeholder={editingMessage ? 'Edit your message...' : 'Write a message...'}
          disabled={isPending}
          className="w-full bg-background border border-border text-foreground-tertiary
            placeholder:text-foreground-disabled rounded-xl py-3.5 pl-5 pr-12 outline-none focus:border-neutral-600
            transition-colors disabled:opacity-50"
        />
        <button
          onClick={handleSendMessage}
          disabled={!messageText.trim() || isPending}
          className="absolute right-3 p-1.5 text-foreground-faint hover:text-emerald-500 transition-colors
            disabled:hover:text-foreground-faint disabled:opacity-50"
        >
          <SendIcon className="size-5" />
        </button>
      </div>
    </div>
  )
}
