export default function RoomSettingsSkeleton() {
  return (
    <div className="flex flex-col h-full px-1 animate-pulse">
      <div className="flex flex-col gap-5">
        {/* Title */}
        <div className="h-3 w-24 bg-background/80 rounded ml-1" />

        {/* Room Name Input Skeleton */}
        <div className="flex flex-col gap-1.5">
          <div className="h-3.5 w-20 bg-background/80 rounded ml-1" />
          <div className="h-[42px] w-full bg-background/50 rounded-xl" />
        </div>

        {/* Room Privacy Dropdown Skeleton */}
        <div className="flex flex-col gap-1.5">
          <div className="h-3.5 w-24 bg-background/80 rounded ml-1" />
          <div className="h-[42px] w-full bg-background/50 rounded-xl" />
        </div>
      </div>

      {/* Action Footer Skeleton */}
      <div className="mt-8 pt-5 border-t border-border/50">
        <div className="h-[42px] w-full bg-background/50 rounded-xl" />
      </div>
    </div>
  )
}
