export default function AdminOverviewLoading() {
  return (
    <div className="flex flex-col gap-6 pb-10">
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        {Array.from({ length: 4 }).map((_, i) => (
          <div key={i} className="bg-surface border border-border rounded-2xl p-5 flex flex-col gap-4 h-[108px]">
            <div className="h-4 w-24 bg-border/50 rounded-md animate-pulse" />
            <div className="h-8 w-32 bg-border/50 rounded-md animate-pulse mt-auto" />
          </div>
        ))}
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mt-2">
        {Array.from({ length: 2 }).map((_, i) => (
          <div
            key={`chart-${i}`}
            className="bg-surface border border-border rounded-2xl p-6 flex flex-col w-full min-h-[380px]"
          >
            <div className="h-6 w-48 bg-border/50 rounded-md animate-pulse mb-6" />

            <div
              className="w-full flex-1 bg-border/20 rounded-xl animate-pulse flex items-end justify-between pb-4 px-2"
            >
              <div className="w-full border-b border-border/40 border-dashed" />
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}
