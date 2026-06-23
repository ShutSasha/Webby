import Link from 'next/link'

import MessagesSquareIcon from '@/assets/icons/shared/message-circle-more.svg'

export function ChatNotFound() {
  return (
    <div className="flex flex-col h-full w-full items-center justify-center text-foreground-faint bg-neutral-900/5">
      <div className="size-16 rounded-2xl bg-neutral-800/50 flex items-center justify-center mb-4">
        <MessagesSquareIcon className="size-8 text-foreground-disabled stroke-2" />
      </div>
      <h3 className="text-xl font-semibold text-foreground-subtle">Chat Not Found</h3>
      <p className="text-sm mt-2 text-center max-w-[300px]">This chat may have been deleted or does not exist.</p>
      <Link
        href="/chats"
        className="mt-6 px-4 py-2 bg-neutral-800 hover:bg-neutral-700 text-foreground-tertiary rounded-xl
          transition-colors text-sm font-medium"
      >
        Back to Chats
      </Link>
    </div>
  )
}
