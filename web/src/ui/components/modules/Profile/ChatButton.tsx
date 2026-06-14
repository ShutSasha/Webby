'use client'

import MailIcon from '@/assets/icons/ic_mail.svg'
import { useCreateChatMutation } from '@/lib/hooks/api/chat/useCreateChat'

import ActionButton from '../../shared/ActionButton'

type Props = {
  targetUserId: string
}

export default function ChatButton({ targetUserId }: Props) {
  const { mutate: createChat, isPending } = useCreateChatMutation()

  const handleCreateChat = () => {
    if (isPending) return
    createChat(targetUserId)
  }

  return (
    <ActionButton
      label={isPending ? 'Opening...' : 'Chat'}
      onClick={handleCreateChat}
      disabled={isPending}
      btnClassName={isPending ? 'opacity-70 cursor-wait' : ''}
    >
      <MailIcon className="w-4 h-4" />
    </ActionButton>
  )
}
