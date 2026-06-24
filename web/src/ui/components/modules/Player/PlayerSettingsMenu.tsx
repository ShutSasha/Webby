type Props = {
  showSettings: boolean
  playbackRate: number
  onSetPlaybackRate: (rate: number) => void
  menuRef: React.RefObject<HTMLDivElement | null>
}

const RATES = [0.25, 0.5, 0.75, 1, 1.25, 1.5, 1.75, 2]

export default function PlayerSettingsMenu({ showSettings, playbackRate, onSetPlaybackRate, menuRef }: Props) {
  if (!showSettings) return null

  return (
    <div
      ref={menuRef}
      className="absolute bottom-14 right-3 w-48 bg-surface-strong/30 backdrop-blur-md rounded-xl overflow-hidden z-20"
      onClick={e => e.stopPropagation()}
    >
      <div className="p-2 border-b border-white/5">
        <p className="text-foreground-muted text-xs font-bold px-3 py-1 uppercase tracking-wider">Speed</p>
      </div>
      <div className="py-1">
        {RATES.map(rate => (
          <button
            key={rate}
            onClick={() => onSetPlaybackRate(rate)}
            className={`w-full flex items-center justify-between px-4 py-2 text-sm transition-colors
            hover:bg-surface-inverse/10
            ${playbackRate === rate ? 'text-emerald-500 font-bold' : 'text-foreground-subtle'}`}
          >
            <span>{rate === 1 ? 'Normal' : `${rate}x`}</span>
            {playbackRate === rate && (
              <div className="w-1.5 h-1.5 rounded-full bg-emerald-500 shadow-[0_0_8px_#10b981]" />
            )}
          </button>
        ))}
      </div>
    </div>
  )
}
