'use client'

import { useEffect, useRef, useState } from 'react'

import Link from 'next/link'
import { useSession } from 'next-auth/react'

import MoreVerticalIcon from '@/assets/icons/shared/more-vertical.svg'
import { useDeleteChatMutation } from '@/lib/hooks/api/chat/useDeleteChat'
import { useGetChatDetailsQuery } from '@/lib/hooks/api/chat/useGetChatDetails'
import Modal from '@/ui/components/shared/Modal'
import SafeImage from '@/ui/components/shared/SafeImage'

import { ChatAreaSkeleton } from './ChatAreaSkeleton'
import ChatMessageInput from './ChatMessageInput'
import { ChatNotFound } from './ChatNotFound'
import ChatMessagesList from './PrivateChat/ChatMessagesList'

type Props = {
  chatId: string
}

export default function ChatArea({ chatId }: Props) {
  const [editingMessage, setEditingMessage] = useState<{ id: string; content: string } | null>(null)
  const { data: session } = useSession()
  const currentUserId = session?.user?.id

  const [isMenuOpen, setIsMenuOpen] = useState(false)
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false)
  const menuRef = useRef<HTMLDivElement>(null)

  const { mutate: deleteChat, isPending: isDeleting } = useDeleteChatMutation()
  const { data: chatDetailsResponse, isLoading: isChatDetailsLoading } = useGetChatDetailsQuery(chatId)

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (menuRef.current && !menuRef.current.contains(event.target as Node)) {
        setIsMenuOpen(false)
      }
    }
    if (isMenuOpen) document.addEventListener('mousedown', handleClickOutside)
    return () => document.removeEventListener('mousedown', handleClickOutside)
  }, [isMenuOpen])

  if (isChatDetailsLoading) {
    return <ChatAreaSkeleton />
  }

  const chatDetails = chatDetailsResponse?.data

  if (!chatDetails) {
    return <ChatNotFound />
  }

  const otherUser = chatDetails.user

  const handleDeleteChat = () => {
    deleteChat(chatId, {
      onSuccess: () => {
        setIsDeleteModalOpen(false)
      },
    })
  }

  return (
    <div className="flex flex-col h-full w-full relative">
      {/* Header */}
      <div
        className="h-[72px] shrink-0 border-b border-neutral-800/60 flex items-center justify-between px-6
          bg-neutral-900/20"
      >
        <div className="flex items-center gap-3">
          {isChatDetailsLoading ? (
            <>
              <div className="size-10 rounded-full bg-neutral-800 animate-pulse shrink-0" />
              <div className="h-5 w-32 bg-neutral-800 rounded-md animate-pulse" />
            </>
          ) : (
            <>
              <SafeImage
                src={otherUser?.avatarUrl || ''}
                fallbackType="user"
                alt="Avatar"
                width={40}
                height={40}
                className="rounded-full object-cover size-10"
              />
              <Link
                href={`/profile/${otherUser?.id}`}
                className="font-medium text-neutral-200 hover:underline hover:cursor-pointer"
              >
                {otherUser?.username?.startsWith('@') ? otherUser.username : `@${otherUser?.username}`}
              </Link>
            </>
          )}
        </div>

        {/* Dropdown Menu Area */}
        <div className="relative" ref={menuRef}>
          <button
            onClick={() => setIsMenuOpen(prev => !prev)}
            className="p-2 rounded-full text-neutral-500 hover:text-neutral-200 hover:bg-neutral-800 transition-colors"
          >
            <MoreVerticalIcon className="size-5" />
          </button>

          {isMenuOpen && (
            <div
              className="absolute right-0 top-full mt-2 w-48 bg-neutral-800 border border-neutral-700/60 shadow-xl
                shadow-black/50 z-50 py-1.5 rounded-xl animate-in fade-in zoom-in-95 duration-200"
            >
              <button
                onClick={() => {
                  setIsMenuOpen(false)
                  setIsDeleteModalOpen(true)
                }}
                className="w-full text-left px-4 py-2 text-sm text-red-400 hover:bg-neutral-700/50 transition-colors"
              >
                Delete chat
              </button>
            </div>
          )}
        </div>
      </div>

      <ChatMessagesList chatId={chatId} currentUserId={currentUserId} onEditMessage={setEditingMessage} />
      <ChatMessageInput chatId={chatId} editingMessage={editingMessage} onCancelEdit={() => setEditingMessage(null)} />

      {/* Delete Confirmation Modal */}
      <Modal isOpen={isDeleteModalOpen} onClose={() => !isDeleting && setIsDeleteModalOpen(false)}>
        <div className="flex flex-col">
          <h3 className="text-xl font-semibold text-neutral-200 mb-2">Delete Chat</h3>
          <p className="text-sm text-neutral-400 mb-6">
            Are you sure you want to delete this chat? All messages will be permanently removed. This action cannot be
            undone.
          </p>
          <div className="flex items-center justify-end gap-3">
            <button
              onClick={() => setIsDeleteModalOpen(false)}
              disabled={isDeleting}
              className="px-4 py-2 text-sm font-medium text-neutral-300 hover:bg-neutral-800 rounded-xl
                transition-colors disabled:opacity-50"
            >
              Cancel
            </button>
            <button
              onClick={handleDeleteChat}
              disabled={isDeleting}
              className="px-4 py-2 text-sm font-medium bg-red-500/10 text-red-500 hover:bg-red-500/20 border
                border-red-500/20 rounded-xl transition-colors disabled:opacity-50"
            >
              {isDeleting ? 'Deleting...' : 'Delete'}
            </button>
          </div>
        </div>
      </Modal>
    </div>
  )
}
