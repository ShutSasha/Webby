import LockIcon from '@/assets/icons/shared/lock.svg'
import ChatArea from '@/ui/components/modules/Chats/ChatArea'
import ChatWebSocketManager from '@/ui/components/modules/Chats/ChatWebSocketManager'
import AuthPlaceholder from '@/ui/components/shared/AuthPlaceholder'
import { auth } from '@/workspace/auth'

type Props = {
  params: Promise<{ chatId: string }>
}

export default async function ChatPage({ params }: Props) {
  const { chatId } = await params
  const session = await auth()

  if (!session) {
    return (
      <AuthPlaceholder
        title="Sign in to view your chats"
        description="Please log in to access your secure inbox, read and reply to your direct messages, review your full chat history, 
        and stay connected with your friends in real time."
        icon={<LockIcon className="size-10 text-neutral-500 stroke-1" />}
      />
    )
  }

  return (
    <>
      <ChatWebSocketManager chatId={chatId} />
      <ChatArea chatId={chatId} />
    </>
  )
}
