export function PlaylistItemSkeleton() {
  return (
    <div className="flex items-center justify-between py-2 px-2.5 mr-1 rounded-xl">
      <div className="flex gap-4 w-full">
        <div className="size-10 rounded-lg bg-neutral-800/80 animate-pulse shrink-0" />
        <div className="flex flex-col justify-center gap-2 w-full">
          <div className="h-3.5 bg-neutral-800/80 rounded w-[60%] animate-pulse" />
          <div className="h-2.5 bg-neutral-800/60 rounded w-[30%] animate-pulse" />
        </div>
      </div>
      <div className="flex items-center justify-center size-6 shrink-0 ml-3">
        <div className="size-5 border-2 border-neutral-700/80 rounded-full animate-pulse" />
      </div>
    </div>
  )
}
