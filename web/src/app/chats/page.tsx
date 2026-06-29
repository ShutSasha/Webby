import LockIcon from '@/assets/icons/shared/lock.svg'
import MessagesSquareIcon from '@/assets/icons/shared/message-circle-more.svg'
import AuthPlaceholder from '@/ui/components/shared/AuthPlaceholder'
import { auth } from '@/workspace/auth'

export default async function ChatsEmptyPage() {
  const session = await auth()

  if (!session) {
    return (
      <AuthPlaceholder
        title="Sign in to view your chats"
        description="Please log in to access your secure inbox, read and reply to your direct messages, review your full chat history, 
        and stay connected with your friends in real time."
        icon={<LockIcon className="size-10 text-foreground-faint stroke-1" />}
      />
    )
  }

  return (
    <div className="flex-1 flex flex-col items-center justify-center text-foreground-faint bg-surface/20">
      <div className="size-16 rounded-2xl bg-background/50 flex items-center justify-center mb-4">
        <MessagesSquareIcon className="size-8 text-foreground-disabled stroke-2" />
      </div>
      <h3 className="text-xl font-semibold text-foreground-subtle">Your Messages</h3>
      <p className="text-sm mt-2">Select a chat from the sidebar to start messaging.</p>
    </div>
  )
}
