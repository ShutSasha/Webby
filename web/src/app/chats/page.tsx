import MessagesSquareIcon from '@/assets/icons/shared/message-circle-more.svg'

export default function ChatsEmptyPage() {
  return (
    <div className="flex-1 flex flex-col items-center justify-center text-neutral-500 bg-neutral-900/20">
      <div className="size-16 rounded-2xl bg-neutral-800/50 flex items-center justify-center mb-4">
        <MessagesSquareIcon className="size-8 text-neutral-600 stroke-2" />
      </div>
      <h3 className="text-xl font-semibold text-neutral-300">Your Messages</h3>
      <p className="text-sm mt-2">Select a chat from the sidebar to start messaging.</p>
    </div>
  )
}
