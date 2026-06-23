import ChatMessagesListSkeleton from './PrivateChat/ChatMessagesListSkeleton'

export function ChatAreaSkeleton() {
  return (
    <div className="flex flex-col h-full w-full relative">
      {/* Header Skeleton */}
      <div className="h-[72px] shrink-0 border-b border-border/60 flex items-center justify-between px-6 bg-surface/20">
        <div className="flex items-center gap-3">
          <div className="size-10 rounded-full bg-background animate-pulse shrink-0" />
          <div className="h-5 w-32 bg-background rounded-md animate-pulse" />
        </div>
        <div className="size-9 rounded-full bg-background animate-pulse shrink-0" />
      </div>

      {/* Messages Skeleton */}
      <div className="flex-1 overflow-y-hidden p-6 flex flex-col-reverse gap-2">
        <ChatMessagesListSkeleton />
      </div>

      {/* Input Skeleton */}
      <div className="p-4 bg-surface/20 border-t border-border/60 shrink-0">
        <div className="h-[52px] w-full bg-background border border-border rounded-xl animate-pulse" />
      </div>
    </div>
  )
}
