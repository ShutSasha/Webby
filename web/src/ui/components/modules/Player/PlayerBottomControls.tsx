import React from 'react'

import FullscreenIcon from '@/assets/icons/Player/fullscreen.svg'
import PauseIcon from '@/assets/icons/Player/pause.svg'
import PlayIcon from '@/assets/icons/Player/play.svg'
import SettingsIcon from '@/assets/icons/Player/settings.svg'
import MaxVolume from '@/assets/icons/Player/volume-max.svg'
import MinVolume from '@/assets/icons/Player/volume-min.svg'
import MutedVolume from '@/assets/icons/Player/volume-muted.svg'

import Duration from './Duration'
import PlayerSyncButton from './PlayerSyncButton'

type Props = {
  playing: boolean
  duration: number
  playedFraction: number
  baseUserVolume: number
  showSettings: boolean
  isRoom?: boolean
  onPlayPause: (e: React.MouseEvent) => void
  onVolumeChange: (e: React.SyntheticEvent<HTMLInputElement>) => void
  onToggleMute: () => void
  onToggleSettings: (e: React.MouseEvent) => void
  onToggleFullScreen: () => void
}

export default function PlayerBottomControls({
  playing,
  duration,
  playedFraction,
  baseUserVolume,
  showSettings,
  isRoom,
  onPlayPause,
  onVolumeChange,
  onToggleMute,
  onToggleSettings,
  onToggleFullScreen,
}: Props) {
  return (
    <div className="flex items-center justify-between">
      <div className="flex items-center gap-3">
        <button onClick={onPlayPause} className="cursor-pointer">
          {playing ? (
            <PauseIcon
              className="w-5 h-5 md:w-6 md:h-6 text-foreground-subtle hover:text-emerald-500 duration-300 ease-out
                transition-colors stroke-[1.5px]"
            />
          ) : (
            <PlayIcon
              className="w-5 h-5 md:w-6 md:h-6 text-foreground-subtle hover:text-emerald-500 duration-300 ease-out
                transition-colors stroke-[1.5px]"
            />
          )}
        </button>

        {/* Duration */}
        <div className="flex items-center gap-0.5">
          <Duration
            seconds={duration * playedFraction}
            className="text-foreground-subtle text-sm font-medium leading-5"
          />
          <span className="text-foreground-subtle text-sm font-medium leading-5">/</span>
          <Duration seconds={duration} className="text-foreground-subtle text-sm font-medium leading-5" />
        </div>

        {/* Volume */}
        <div className="flex items-center gap-2 group/volume">
          <button onClick={onToggleMute} className="cursor-pointer">
            {baseUserVolume >= 0.5 && (
              <MaxVolume
                className="text-foreground-subtle w-5 h-5 md:w-6 md:h-6 group-hover/volume:text-emerald-500
                  transition-colors stroke-[1.5px]"
              />
            )}
            {baseUserVolume < 0.5 && baseUserVolume > 0 && (
              <MinVolume
                className="text-foreground-subtle w-5 h-5 md:w-6 md:h-6 group-hover/volume:text-emerald-500
                  transition-colors stroke-[1.5px]"
              />
            )}
            {baseUserVolume === 0 && (
              <MutedVolume
                className="text-foreground-subtle w-5 h-5 md:w-6 md:h-6 group-hover/volume:text-emerald-500
                  transition-colors stroke-[1.5px]"
              />
            )}
          </button>

          {/* Volume Slider */}
          <div className="relative w-22 h-1 md:h-1.5 bg-white/20 rounded-full group/slider">
            <div
              className="absolute h-full bg-emerald-500 rounded-full"
              style={{ width: `${baseUserVolume * 100}%` }}
            />
            <div
              className="absolute w-2.5 h-2.5 md:h-3 md:w-3 bg-emerald-500 rounded-full top-1/2 -translate-x-1/2
                -translate-y-1/2 pointer-events-none transition-transform"
              style={{ left: `calc(${baseUserVolume * 100}% + (${(0.5 - baseUserVolume) * 10}px))` }}
            />
            <input
              type="range"
              min={0}
              max={1}
              step="any"
              value={baseUserVolume}
              onChange={onVolumeChange}
              className="absolute inset-0 w-full h-full opacity-0 cursor-pointer z-10"
            />
          </div>
        </div>
      </div>

      <div className="flex items-center gap-2 md:gap-3">
        <button type="button" onClick={onToggleSettings} className="cursor-pointer group/settings">
          <SettingsIcon
            className={`w-5 h-5 md:w-6 md:h-6 transition-colors duration-300 stroke-[1.5px] ${
              showSettings
                ? 'text-emerald-500 rotate-45'
                : 'text-foreground-subtle group-hover/settings:text-emerald-500'
              }`}
          />
        </button>
        <button type="button" onClick={onToggleFullScreen} className="cursor-pointer group/fullscreen">
          <FullscreenIcon
            className="text-foreground-subtle w-5 h-5 md:w-6 md:h-6 group-hover/fullscreen:text-emerald-500
              transition-colors stroke-[1.5px]"
          />
        </button>
      </div>
    </div>
  )
}
