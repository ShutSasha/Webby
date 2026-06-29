export default function ModerationVideosTableSkeleton() {
  const skeletonItems = Array.from({ length: 10 })

  return (
    <div className="bg-surface border border-border rounded-2xl overflow-hidden flex flex-col">
      <div className="overflow-x-auto custom-scrollbar">
        <table className="w-full text-left border-collapse min-w-[800px]">
          <thead>
            <tr className="border-b border-border bg-surface/50 text-foreground-faint text-sm">
              <th className="py-4 px-6 font-medium">Video</th>
              <th className="py-4 px-6 font-medium">Status</th>
              <th className="py-4 px-6 font-medium">Created At</th>
              <th className="py-4 px-6 font-medium text-right">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-background/60">
            {skeletonItems.map((_, index) => (
              <tr key={index} className="animate-pulse">
                <td className="py-4 px-6">
                  <div className="flex items-center gap-4">
                    <div className="w-28 h-16 rounded-md bg-border/50 shrink-0" />

                    <div className="flex flex-col gap-2">
                      <div className="h-4 w-40 bg-border/50 rounded" />
                      <div className="h-3 w-24 bg-border/30 rounded" />
                    </div>
                  </div>
                </td>

                <td className="py-4 px-6">
                  <div className="flex items-center gap-2">
                    <div className="size-2 rounded-full bg-border/50" />
                    <div className="h-4 w-16 bg-border/30 rounded" />
                  </div>
                </td>

                <td className="py-4 px-6">
                  <div className="h-4 w-24 bg-border/30 rounded" />
                </td>

                <td className="py-4 px-6 flex justify-end items-center gap-2 h-[97px]">
                  <div className="w-[68px] h-8 bg-border/50 rounded-lg" />
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}
