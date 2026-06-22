import QueueItemSkeleton from './Queue/QueueItemSkeleton'

export default function RoomInteractionSkeleton() {
  return (
    <div className="flex flex-col w-full h-full animate-pulse">
      <div className="border-b border-neutral-700/80 flex flex-row items-center justify-between py-2 px-3 mb-2 shrink-0">
        <div className="size-6 bg-neutral-800/80 rounded" />
        <div className="flex flex-row items-center gap-4">
          <div className="size-6 bg-neutral-800/80 rounded" />
          <div className="size-6 bg-neutral-800/80 rounded" />
          <div className="size-6 bg-neutral-800/80 rounded" />
        </div>
      </div>

      <div className="mb-2 shrink-0 w-full">
        <div className="h-[38px] w-full bg-neutral-900/80 rounded-md" />
      </div>

      <div className="flex flex-col gap-2 flex-1 overflow-hidden">
        {Array.from({ length: 6 }).map((_, i) => (
          <QueueItemSkeleton key={i} />
        ))}
      </div>
    </div>
  )
}
