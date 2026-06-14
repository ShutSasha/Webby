'use client'

import { useEffect, useRef, useState } from 'react'

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

    if (!trimmedMessage || isPending) return

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
        { chatId, content: trimmedMessage },
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
    <div
      className="p-4 bg-neutral-900/20 border-t border-neutral-800/60 shrink-0 flex flex-col transition-all
        duration-300"
    >
      {editingMessage && (
        <div
          className="flex items-center justify-between mb-3 px-2 animate-in fade-in slide-in-from-bottom-2 duration-200"
        >
          <div className="flex items-center gap-3 overflow-hidden">
            <EditPenIcon className="size-5 text-emerald-500 shrink-0" />
            <div className="flex flex-col overflow-hidden border-l-2 border-emerald-500/50 pl-2">
              <span className="text-xs font-semibold text-emerald-500 leading-tight">Edit message</span>
              <span className="text-sm text-neutral-400 truncate max-w-[200px] sm:max-w-[400px]">
                {editingMessage.content}
              </span>
            </div>
          </div>
          <button
            onClick={onCancelEdit}
            className="p-1.5 hover:bg-neutral-800 rounded-full text-neutral-500 hover:text-neutral-300
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
          className="w-full bg-[#141414] border border-neutral-800 text-neutral-200 placeholder:text-neutral-600
            rounded-xl py-3.5 pl-5 pr-12 outline-none focus:border-neutral-600 transition-colors disabled:opacity-50"
        />
        <button
          onClick={handleSendMessage}
          disabled={!messageText.trim() || isPending}
          className="absolute right-3 p-1.5 text-neutral-500 hover:text-emerald-500 transition-colors
            disabled:hover:text-neutral-500 disabled:opacity-50"
        >
          <SendIcon className="size-5" />
        </button>
      </div>
    </div>
  )
}
