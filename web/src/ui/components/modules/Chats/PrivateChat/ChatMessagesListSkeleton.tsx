import ChatMessageSkeleton from './ChatMessageSkeleton'

export default function ChatMessagesListSkeleton() {
  return (
    <>
      <ChatMessageSkeleton isMe className="h-12 w-[60%] max-w-[260px]" />
      <ChatMessageSkeleton isMe className="h-10 w-[45%] max-w-[180px]" />

      <ChatMessageSkeleton className="h-16 w-[75%] max-w-[320px] mt-2" />
      <ChatMessageSkeleton className="h-10 w-[50%] max-w-[200px]" />

      <ChatMessageSkeleton isMe className="h-14 w-[65%] max-w-[280px] mt-2" />
      <ChatMessageSkeleton className="h-12 w-[55%] max-w-[220px] mt-2" />
    </>
  )
}
