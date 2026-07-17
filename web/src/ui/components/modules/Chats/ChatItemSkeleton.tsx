export default function ChatItemSkeleton() {
  return (
    <div className="flex items-center gap-3 p-3 rounded-xl">
      {/* Avatar skeleton */}
      <div className="size-12 rounded-full bg-background animate-pulse shrink-0" />

      {/* Text container */}
      <div className="flex flex-col overflow-hidden w-full">
        <div className="h-4 w-24 bg-background rounded-md animate-pulse" />
        <div className="h-3 w-full max-w-[180px] bg-background rounded-md animate-pulse mt-2" />
      </div>
    </div>
  )
}
