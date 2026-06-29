export function FollowsContainerSkeleton() {
  return (
    <div className="grid grid-cols-1 md:grid-cols-2 2xl:grid-cols-3 gap-4">
      {Array.from({ length: 24 }).map((_, index) => (
        <FollowItemSkeleton key={index} />
      ))}
    </div>
  )
}
function FollowItemSkeleton() {
  return (
    <div
      className="bg-background/50 animate-pulse rounded-xl p-2 flex items-center justify-between gap-2 border
        border-transparent"
    >
      <div className="flex items-center gap-2">
        <div className="w-13 h-13 rounded-full bg-surface-tertiary/50" />

        <div className="flex flex-col gap-2">
          <div className="h-4 w-24 bg-surface-tertiary/60 rounded-md" />

          <div className="h-3 w-16 bg-surface-tertiary/60 rounded-md" />
        </div>
      </div>

      <div className="flex items-center gap-5 pr-2">
        <div className="w-5 h-5 bg-surface-tertiary/60 rounded-full" />
        <div className="w-5 h-5 bg-surface-tertiary/60 rounded-full" />
      </div>
    </div>
  )
}
