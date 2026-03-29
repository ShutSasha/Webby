export default function PlayerSkeleton() {
  return (
    <div className="h-fit min-w-0 w-full animate-pulse">
      {/* player */}
      <div className="w-full aspect-video bg-neutral-800 rounded-2xl mb-3" />

      {/* title and buttons */}
      <div className="flex items-center justify-between mt-3 mb-2 gap-5">
        <div className="h-[30px] bg-neutral-800 rounded-lg w-1/3" />
        <div className="flex gap-3">
          <div className="w-[152px] h-[34px] bg-neutral-800 rounded-full" />
          <div className="w-[166px] h-[34px] bg-neutral-800 rounded-full" />
        </div>
      </div>

      {/* user */}
      <div className="flex items-center gap-3 mt-4">
        <div className="size-9 bg-neutral-800 rounded-full shrink-0" />
        <div className="w-32 h-6 bg-neutral-800 rounded-md" />
        <div className="w-[100px] h-[34px] bg-neutral-800 rounded-full ml-2" />
      </div>

      {/* description */}
      <div className="w-full h-24 bg-neutral-800 rounded-xl mt-4" />
    </div>
  )
}
