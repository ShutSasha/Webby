export default function VideoDetailsSkeleton() {
  return (
    <div className="h-fit min-w-0 w-full animate-pulse">
      {/* Player */}
      <div className="w-full aspect-video bg-background rounded-2xl" />

      {/* Title Area */}
      <div className="flex items-start justify-between mt-3 mb-2">
        <div className="h-7.5 bg-background rounded-lg w-3/4 md:w-1/2" />
      </div>

      {/* User & Actions Area */}
      <div className="flex items-center justify-between flex-wrap gap-2">
        <div className="flex items-center gap-3">
          <div className="flex items-center gap-3">
            <div className="size-9 bg-background rounded-full shrink-0" />
            <div className="h-5 w-24 sm:w-32 bg-background rounded-md" />
          </div>
          {/* Follow Button */}
          <div className="h-8 w-24 bg-background rounded-full" />
        </div>

        <div className="flex gap-3 items-center">
          <div className="h-9 w-36 bg-background rounded-full" />
          <div className="h-9 w-[150px] bg-background rounded-full" />
        </div>
      </div>

      {/* Description Area */}
      <div className="flex flex-col gap-2 rounded-xl bg-background/40 p-3 mt-4 overflow-hidden">
        {/* Views & Date */}
        <div className="h-5 w-40 bg-surface-tertiary/50 rounded-md" />

        {/* Description Text */}
        <div className="h-[69px] w-full bg-surface-tertiary/50 rounded-md" />
      </div>
    </div>
  )
}
