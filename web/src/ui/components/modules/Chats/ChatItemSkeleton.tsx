export default function ChatItemSkeleton() {
  return (
    <div className="flex items-center gap-3 p-3 rounded-xl">
      {/* Avatar skeleton */}
      <div className="size-12 rounded-full bg-neutral-800 animate-pulse shrink-0" />

      {/* Text container */}
      <div className="flex flex-col overflow-hidden w-full">
        <div className="h-4 w-24 bg-neutral-800 rounded-md animate-pulse" />
        <div className="h-3 w-full max-w-[180px] bg-neutral-800 rounded-md animate-pulse mt-2" />
      </div>
    </div>
  )
}
