export default async function PageLoading() {
  return (
    <div className="flex flex-1 flex-col justify-center items-center gap-4 min-h-[40vh]">
      <div className="flex items-center gap-1.5 h-10">
        <div
          className="w-1.5 h-full bg-emerald-500 rounded-full animate-pulse shadow-[0_0_8px_rgba(16,185,129,0.5)]"
          style={{ animationDelay: '0ms' }}
        ></div>
        <div
          className="w-1.5 h-3/4 bg-emerald-400 rounded-full animate-pulse shadow-[0_0_8px_rgba(16,185,129,0.5)]"
          style={{ animationDelay: '200ms' }}
        ></div>
        <div
          className="w-1.5 h-1/2 bg-emerald-300 rounded-full animate-pulse shadow-[0_0_8px_rgba(16,185,129,0.5)]"
          style={{ animationDelay: '400ms' }}
        ></div>
        <div
          className="w-1.5 h-3/4 bg-emerald-400 rounded-full animate-pulse shadow-[0_0_8px_rgba(16,185,129,0.5)]"
          style={{ animationDelay: '600ms' }}
        ></div>
        <div
          className="w-1.5 h-full bg-emerald-500 rounded-full animate-pulse shadow-[0_0_8px_rgba(16,185,129,0.5)]"
          style={{ animationDelay: '800ms' }}
        ></div>
      </div>

      <p className="text-neutral-500 text-[11px] font-bold tracking-[0.3em] uppercase">Loading...</p>
    </div>
  )
}
