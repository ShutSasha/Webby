import React, { useState } from 'react'

import Duration from './Duration'

type Props = {
  played: number
  loaded: number
  duration: number
  onSeekMouseDown: () => void
  onSeekChange: (e: React.SyntheticEvent<HTMLInputElement>) => void
  onSeekMouseUp: (newTimeFraction: number) => void
}

export default function PlayerProgressBar({
  played,
  loaded,
  duration,
  onSeekMouseDown,
  onSeekChange,
  onSeekMouseUp,
}: Props) {
  const [hoverTime, setHoverTime] = useState<number | null>(null)
  const [hoverX, setHoverX] = useState<number>(0)

  const handleMouseMove = (e: React.MouseEvent<HTMLDivElement>) => {
    const rect = e.currentTarget.getBoundingClientRect()
    const x = e.clientX - rect.left
    const percentage = Math.max(0, Math.min(1, x / rect.width))

    setHoverTime(percentage * duration)
    setHoverX(x)
  }

  const handleMouseLeave = () => {
    setHoverTime(null)
  }

  const handleMouseUp = (e: React.SyntheticEvent<HTMLInputElement>) => {
    const inputTarget = e.target as HTMLInputElement
    const fraction = hoverTime !== null ? hoverTime / duration : Number.parseFloat(inputTarget.value)
    onSeekMouseUp(fraction)
  }

  return (
    <div
      className="relative h-1.5 w-full bg-white/20 rounded-full group/bar cursor-pointer"
      onMouseMove={handleMouseMove}
      onMouseLeave={handleMouseLeave}
    >
      {/* Tooltip */}
      {hoverTime !== null && (
        <div
          className="absolute bottom-4 -translate-x-1/2 bg-white text-black px-1.5 py-0.5 rounded-md text-[12px]
            font-bold shadow-lg pointer-events-none transition-opacity"
          style={{ left: `${hoverX}px` }}
        >
          <Duration seconds={hoverTime} />
          <div className="absolute -bottom-1 left-1/2 -translate-x-1/2 w-2 h-2 bg-white rotate-45" />
        </div>
      )}

      {/* Pre-Loaded line */}
      <div className="absolute h-full bg-white/30 rounded-full transition-all" style={{ width: `${loaded * 100}%` }} />

      {/* Played line */}
      <div className="absolute h-full bg-emerald-500 rounded-full" style={{ width: `${played * 100}%` }} />

      {/* Played Circle */}
      <div
        className="absolute h-3 w-3 bg-emerald-500 rounded-full top-1/2 -translate-x-1/2 -translate-y-1/2
          pointer-events-none"
        style={{ left: `${played * 100}%` }}
      />

      {/* Invisible Input for control */}
      <input
        type="range"
        min={0}
        max={0.999999}
        step="any"
        value={played}
        onMouseDown={onSeekMouseDown}
        onChange={onSeekChange}
        onMouseUp={handleMouseUp}
        className="absolute inset-0 w-full h-full opacity-0 cursor-pointer z-10"
      />
    </div>
  )
}
