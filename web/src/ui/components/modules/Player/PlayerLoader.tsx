type Props = {
  isReady: boolean
  buffering: boolean
}

export default function PlayerLoader({ isReady, buffering }: Props) {
  if (isReady && !buffering) return null

  return (
    <div className="absolute inset-0 z-1 flex items-center justify-center pointer-events-none bg-black/40">
      <div className="w-16 h-16 border-6 border-emerald-500/20 border-t-emerald-500 rounded-full animate-spin" />
    </div>
  )
}
