export default function TableSkeleton() {
  return (
    <div className="bg-surface border border-border rounded-2xl overflow-hidden flex flex-col">
      <div className="overflow-x-auto custom-scrollbar">
        <table className="w-full text-left border-collapse min-w-[800px]">
          <thead>
            <tr className="border-b border-border bg-surface/50 text-foreground-faint text-sm">
              <th className="py-4 px-6 font-medium">User</th>
              <th className="py-4 px-6 font-medium">Role</th>
              <th className="py-4 px-6 font-medium">Status</th>
              <th className="py-4 px-6 font-medium">Registered</th>
              <th className="py-4 px-6 font-medium text-right">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-background/60">
            {Array.from({ length: 5 }).map((_, index) => (
              <tr key={index} className="animate-pulse">
                <td className="py-4 px-6">
                  <div className="flex items-center gap-3">
                    <div className="size-10 rounded-full bg-border/50 shrink-0" />
                    <div className="flex flex-col gap-2">
                      <div className="h-4 w-32 bg-border/50 rounded-md" />
                      <div className="h-3 w-40 bg-border/40 rounded-md" />
                    </div>
                  </div>
                </td>
                <td className="py-4 px-6">
                  <div className="h-6 w-20 bg-border/50 rounded-md" />
                </td>
                <td className="py-4 px-6">
                  <div className="flex items-center gap-2">
                    <div className="size-2 rounded-full bg-border/50" />
                    <div className="h-4 w-16 bg-border/50 rounded-md" />
                  </div>
                </td>
                <td className="py-4 px-6">
                  <div className="h-4 w-24 bg-border/50 rounded-md" />
                </td>
                <td className="py-4 px-6 flex justify-end gap-2">
                  <div className="h-8 w-20 bg-border/50 rounded-lg" />
                  <div className="h-8 w-16 bg-border/50 rounded-lg" />
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}
