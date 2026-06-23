export default function QueueItemSkeleton() {
  return (
    <div className="flex flex-col gap-1">
      <div className="flex items-center justify-between p-2 rounded-xl bg-surface/50 border-b-2 border-transparent">
        <div className="flex items-center gap-3 overflow-hidden w-full">
          <div className="size-10 rounded-lg bg-background/80 animate-pulse shrink-0" />

          <div className="h-3.5 w-3/5 bg-background/80 rounded animate-pulse" />
        </div>

        <div className="flex items-center gap-1 shrink-0 ml-4">
          <div className="size-7 rounded-full bg-background/80 animate-pulse" />
        </div>
      </div>
    </div>
  )
}
