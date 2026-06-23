export default function VideoNotFound() {
  return (
    <div className="h-fit min-w-0 w-full animate-in fade-in duration-500">
      <div
        className="w-full aspect-video bg-surface rounded-2xl mb-3 flex flex-col items-center justify-center border
          border-border/50"
      >
        <div className="size-16 rounded-full bg-background flex items-center justify-center mb-4">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            strokeWidth={1.5}
            stroke="currentColor"
            className="size-8 text-foreground-faint"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              d="M15.75 10.5l4.72-4.72a.75.75 0 011.28.53v11.38a.75.75 0 01-1.28.53l-4.72-4.72M12 18.75H4.5a2.25 2.25 0 01-2.25-2.25V7.5A2.25 2.25 0 014.5 5.25h7.5A2.25 2.25 0 0114.25 7.5v11.25a2.25 2.25 0 01-2.25 2.25z"
            />
            <path strokeLinecap="round" strokeLinejoin="round" d="M3 3l18 18" />
          </svg>
        </div>

        <p className="text-foreground-subtle text-[20px] font-bold">Video unavailable</p>
        <p className="text-foreground-faint text-sm mt-2 max-w-sm text-center">
          This video has been deleted, hidden, or the link you followed is invalid.
        </p>
      </div>

      <div className="h-8 bg-background/50 rounded-lg w-1/3 mt-6 mb-4" />
      <div className="flex gap-3 items-center mb-6">
        <div className="size-9 bg-background/50 rounded-full" />
        <div className="h-5 bg-background/50 rounded-md w-32" />
      </div>
      <div className="w-full h-24 bg-background/50 rounded-xl" />
    </div>
  )
}
