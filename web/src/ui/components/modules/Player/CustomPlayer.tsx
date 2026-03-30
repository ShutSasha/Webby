'use client'
import { useCallback, useEffect, useRef, useState } from 'react'

import ReactPlayer from 'react-player'

import FilledPlay from '@/assets/icons/Player/filled-play.svg'
import FullscreenIcon from '@/assets/icons/Player/fullscreen.svg'
import PauseIcon from '@/assets/icons/Player/pause.svg'
import PlayIcon from '@/assets/icons/Player/play.svg'
import SettingsIcon from '@/assets/icons/Player/settings.svg'
import MaxVolume from '@/assets/icons/Player/volume-max.svg'
import MinVolume from '@/assets/icons/Player/volume-min.svg'
import MutedVolume from '@/assets/icons/Player/volume-muted.svg'
import { useIsClient } from '@/lib/hooks/useIsClient'
import { usePlayerControls } from '@/lib/hooks/usePlayerControls'
import { usePlayerHotkeys } from '@/lib/hooks/usePlayerHotkeys'
import { usePlayerPlayStore, usePlayerStore } from '@/stores/player.store'

import Duration from './Duration'

type PlayerProps = {
  videoUrl: string
}

export default function CustomPlayer({ videoUrl }: PlayerProps) {
  const playerRef = useRef<HTMLVideoElement | null>(null)
  const playerContainerRef = useRef<HTMLDivElement>(null)
  const isMounted = useIsClient()

  const settingsContainerRef = useRef<HTMLDivElement>(null)
  const [isFullScreen, setIsFullScreen] = useState(false)
  const [prevVolume, setPrevVolume] = useState(1)
  const hideTimeoutRef = useRef<NodeJS.Timeout | null>(null)
  const [showCustomControls, setShowCustomControls] = useState(true)
  const isPlatformMode = usePlayerControls(videoUrl)

  const [hoverTime, setHoverTime] = useState<number | null>(null)
  const [hoverX, setHoverX] = useState<number>(0)

  const baseUserVolume = usePlayerStore(state => state.baseVolume)
  const setBaseUserVolume = usePlayerStore(state => state.setBaseVolume)

  const setPlayerRef = useCallback((player: HTMLVideoElement) => {
    if (!player) return
    playerRef.current = player
  }, [])

  const triggerEnded = usePlayerPlayStore(state => state.triggerEnded)
  const togglePlay = usePlayerPlayStore(state => state.togglePlay)
  const setPlaying = usePlayerPlayStore(state => state.setPlaying)
  const playing = usePlayerPlayStore(state => state.playing)

  const initialState = {
    pip: false,
    light: false,
    muted: false,
    played: 0,
    loaded: 0,
    duration: 0,
    playbackRate: 1.0,
    loop: false,
    seeking: false,
    loadedSeconds: 0,
    playedSeconds: 0,
    showSettings: false,
  }

  type PlayerState = typeof initialState

  const [state, setState] = useState<PlayerState>(initialState)

  const { light, muted, loop, played, loaded, duration, playbackRate, pip, showSettings } = state

  useEffect(() => {
    setState(prev => ({
      ...prev,
      played: 0,
      loaded: 0,
      duration: 0,
      loadedSeconds: 0,
      playedSeconds: 0,
    }))

    const hasInteracted = typeof navigator !== 'undefined' && (navigator as any).userActivation?.hasBeenActive

    if (hasInteracted === false) {
      setPlaying(false)
    } else {
      setPlaying(true)
    }
  }, [videoUrl, setPlaying])

  usePlayerHotkeys({
    playerRef,
    hideTimeoutRef,
    duration,
    togglePlay,
    setShowCustomControls,
    setState,
  })

  const handlePlayPause = (e: React.MouseEvent) => {
    if (showSettings) {
      if (settingsContainerRef.current && !settingsContainerRef.current.contains(e.target as Node)) {
        setState(prevState => ({ ...prevState, showSettings: false }))
        hideTimeoutRef.current = setTimeout(() => setShowCustomControls(false), 3000)
        return
      }
    }

    togglePlay()
  }

  const handleSetPlaybackRate = (event: React.SyntheticEvent<HTMLButtonElement>) => {
    const buttonTarget = event.currentTarget as HTMLButtonElement
    const btnData = Number.parseFloat(`${buttonTarget.dataset.value}`)

    if (isNaN(btnData)) return

    setState(prevState => ({
      ...prevState,
      playbackRate: btnData,
    }))
  }

  const handleRateChange = () => {
    const player = playerRef.current
    if (!player) return

    setState(prevState => ({ ...prevState, playbackRate: player.playbackRate }))
  }

  const toggleSettings = (e: React.MouseEvent) => {
    e.stopPropagation()
    setState(prev => ({ ...prev, showSettings: !prev.showSettings }))
  }

  const handleProgress = () => {
    const player = playerRef.current
    // We only want to update time slider if we are not currently seeking
    if (!player || state.seeking || !player.buffered?.length) return

    setState(prevState => ({
      ...prevState,
      loadedSeconds: player.buffered?.end(player.buffered?.length - 1),
      loaded: player.buffered?.end(player.buffered?.length - 1) / player.duration,
    }))
  }

  const handleSeekMouseDown = () => {
    setState({ ...state, seeking: true })
  }

  const handleSeekChange = (event: React.SyntheticEvent<HTMLInputElement>) => {
    const inputTarget = event.target as HTMLInputElement
    setState(prevState => ({ ...prevState, played: Number.parseFloat(inputTarget.value) }))
  }

  const handleSeekMouseUp = (event: React.SyntheticEvent<HTMLInputElement>) => {
    const inputTarget = event.target as HTMLInputElement
    const player = playerRef.current

    if (!player) return

    const newTime = hoverTime !== null ? hoverTime : Number.parseFloat(inputTarget.value) * duration

    setState(prevState => ({
      ...prevState,
      seeking: false,
      played: newTime / duration,
      playedSeconds: newTime,
    }))

    player.currentTime = newTime
  }

  const handleTimeUpdate = () => {
    const player = playerRef.current
    // We only want to update time slider if we are not currently seeking
    if (!player || state.seeking) return

    if (!player.duration) return

    setState(prevState => ({
      ...prevState,
      playedSeconds: player.currentTime,
      played: player.currentTime / player.duration,
    }))
  }

  const handleDurationChange = () => {
    const player = playerRef.current
    if (!player) return

    setState(prevState => ({ ...prevState, duration: player.duration }))
  }

  const handleVolumeChange = (event: React.SyntheticEvent<HTMLInputElement>) => {
    const inputTarget = event.target as HTMLInputElement
    const volume = Number.parseFloat(inputTarget.value)

    setBaseUserVolume(volume)
  }

  const toggleMute = () => {
    if (baseUserVolume > 0) {
      setPrevVolume(baseUserVolume)
      setBaseUserVolume(0)
    } else {
      setBaseUserVolume(prevVolume)
    }
  }

  const toggleFullScreen = () => {
    if (!isFullScreen) {
      playerContainerRef.current?.requestFullscreen()
    } else {
      document.exitFullscreen()
    }
  }

  const handleMouseMove = () => {
    if (isPlatformMode) return

    setShowCustomControls(true)

    if (hideTimeoutRef.current) {
      clearTimeout(hideTimeoutRef.current)
    }

    if (playing && !showSettings) {
      hideTimeoutRef.current = setTimeout(() => {
        setShowCustomControls(false)
      }, 3000)
    }
  }

  useEffect(() => {
    const handleFsChange = () => {
      setIsFullScreen(!!document.fullscreenElement)
    }
    document.addEventListener('fullscreenchange', handleFsChange)
    return () => document.removeEventListener('fullscreenchange', handleFsChange)
  }, [])

  useEffect(() => {
    if (!playing) {
      setShowCustomControls(true)
      if (hideTimeoutRef.current) clearTimeout(hideTimeoutRef.current)
    } else {
      hideTimeoutRef.current = setTimeout(() => setShowCustomControls(false), 3000)
    }
    return () => {
      if (hideTimeoutRef.current) clearTimeout(hideTimeoutRef.current)
    }
  }, [playing])

  const handleProgressMouseMove = (e: React.MouseEvent<HTMLDivElement>) => {
    const rect = e.currentTarget.getBoundingClientRect()
    const x = e.clientX - rect.left
    const percentage = Math.max(0, Math.min(1, x / rect.width))

    setHoverTime(percentage * duration)
    setHoverX(x)
  }

  const handleProgressMouseLeave = () => {
    setHoverTime(null)
  }

  if (!isMounted)
    return (
      <div className="aspect-video bg-black w-full rounded-2xl flex items-center justify-center">
        <div className="w-12 h-12 border-4 border-emerald-500/20 border-t-emerald-500 rounded-full animate-spin" />
      </div>
    )

  return (
    <div
      ref={playerContainerRef}
      onMouseMove={handleMouseMove}
      onMouseLeave={() => playing && setShowCustomControls(false)}
      className={`group relative w-full bg-transparent overflow-hidden transition-all
        ${isFullScreen ? 'w-screen h-screen rounded-0' : 'aspect-video rounded-2xl'}
        ${!isPlatformMode && !showCustomControls ? 'cursor-none' : 'cursor-default'}`}
    >
      <ReactPlayer
        className="react-player"
        ref={setPlayerRef}
        playing={playing}
        controls={isPlatformMode}
        width="100%"
        height="100%"
        light={light}
        muted={muted}
        loop={loop}
        volume={baseUserVolume}
        playbackRate={playbackRate}
        pip={pip}
        src={videoUrl}
        onRateChange={handleRateChange}
        onProgress={handleProgress}
        onTimeUpdate={handleTimeUpdate}
        onDurationChange={handleDurationChange}
        onVolumeChange={(e: any) => {
          if (typeof e === 'number') {
            setBaseUserVolume(e)
            return
          }

          const volume = e.target?.volume
          if (typeof volume === 'number') {
            setBaseUserVolume(volume)
          }
        }}
        onPlay={() => setPlaying(true)}
        onPause={() => setPlaying(false)}
        onEnded={() => {
          triggerEnded()
        }}
      />

      {/* Overlay */}
      {!isPlatformMode && (
        <div
          className={`absolute inset-0 bg-linear-to-t from-black/80 via-transparent to-transparent transition-opacity
          duration-400 ease-in-out ${showCustomControls ? 'opacity-100' : 'opacity-0'}`}
          onClick={handlePlayPause}
        >
          {/* Play button on the overlay */}
          <div
            className={`${playing === false ? 'opacity-100' : 'opacity-0'} absolute w-12 h-12 md:w-18 md:h-18
            bg-black/15 rounded-full top-1/2 left-1/2 -translate-y-1/2 -translate-x-1/2 flex items-center justify-center
            pl-1 transition-opacity duration-400 ease-in-out`}
          >
            <FilledPlay className="w-5 h-5 md:w-8 md:h-8 text-white/70" />
          </div>

          {/* Settings Menu Popup */}
          {showSettings && (
            <div
              ref={settingsContainerRef}
              className="absolute bottom-14 right-3 w-48 bg-black/30 backdrop-blur-md rounded-xl overflow-hidden z-20"
              onClick={e => e.stopPropagation()}
            >
              <div className="p-2 border-b border-white/5">
                <p className="text-neutral-400 text-xs font-bold px-3 py-1 uppercase tracking-wider">Speed</p>
              </div>
              <div className="py-1">
                {[0.25, 0.5, 0.75, 1, 1.25, 1.5, 1.75, 2].map(rate => (
                  <button
                    key={rate}
                    onClick={handleSetPlaybackRate}
                    className={`w-full flex items-center justify-between px-4 py-2 text-sm transition-colors
                    hover:bg-white/10 ${playbackRate === rate ? 'text-emerald-500 font-bold' : 'text-neutral-300'}`}
                    data-value={rate}
                  >
                    <span>{rate === 1 ? 'Normal' : `${rate}x`}</span>
                    {playbackRate === rate && (
                      <div className="w-1.5 h-1.5 rounded-full bg-emerald-500 shadow-[0_0_8px_#10b981]" />
                    )}
                  </button>
                ))}
              </div>
            </div>
          )}

          {/* Controls */}
          <div
            onClick={e => e.stopPropagation()}
            className={`absolute bottom-3 left-3 right-3 flex flex-col gap-2 transition-transform duration-500
            ${showCustomControls ? 'translate-y-0' : 'translate-y-10 pointer-events-none'}`}
          >
            {/* Progress Bar Container */}
            <div
              className="relative h-1.5 w-full bg-white/20 rounded-full group/bar cursor-pointer"
              onMouseMove={handleProgressMouseMove}
              onMouseLeave={handleProgressMouseLeave}
            >
              {/* Hint (Tooltip) */}
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
              <div
                className="absolute h-full bg-white/30 rounded-full transition-all"
                style={{ width: `${loaded * 100}%` }}
              />

              {/* Played line */}
              <div className="absolute h-full bg-emerald-500 rounded-full" style={{ width: `${played * 100}%` }} />

              {/* Played Circle */}
              <div
                className="absolute h-2.5 w-2.5 bg-emerald-500 rounded-full top-1/2 -translate-x-1/2 -translate-y-1/2
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
                onMouseDown={handleSeekMouseDown}
                onChange={handleSeekChange}
                onMouseUp={handleSeekMouseUp}
                className="absolute inset-0 w-full h-full opacity-0 cursor-pointer z-10"
              />
            </div>

            {/* Buttons etc */}
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-3">
                <button onClick={handlePlayPause} className="cursor-pointer">
                  {playing ? (
                    <PauseIcon
                      className="w-6 h-6 text-neutral-300 hover:text-emerald-500 duration-300 ease-out
                        transition-colors"
                    />
                  ) : (
                    <PlayIcon
                      className="w-6 h-6 text-neutral-300 hover:text-emerald-500 duration-300 ease-out
                        transition-colors"
                    />
                  )}
                </button>

                {/* Duration */}
                <div className="flex items-center gap-0.5">
                  <Duration seconds={duration * played} className="text-neutral-300 text-sm font-medium leading-5" />
                  <span className="text-neutral-300 text-sm font-medium leading-5">/</span>
                  <Duration seconds={duration} className="text-neutral-300 text-sm font-medium leading-5" />
                </div>
                {/* Volume */}
                <div className="flex items-center gap-2 group/volume">
                  <button onClick={toggleMute} className="cursor-pointer">
                    {baseUserVolume >= 0.5 && (
                      <MaxVolume
                        className="text-neutral-300 w-6 h-6 group-hover/volume:text-emerald-500 transition-colors"
                      />
                    )}
                    {baseUserVolume < 0.5 && baseUserVolume > 0 && (
                      <MinVolume
                        className="text-neutral-300 w-6 h-6 group-hover/volume:text-emerald-500 transition-colors"
                      />
                    )}
                    {baseUserVolume === 0 && (
                      <MutedVolume
                        className="text-neutral-300 w-6 h-6 group-hover/volume:text-emerald-500 transition-colors"
                      />
                    )}
                  </button>

                  {/* Volume Slider */}
                  <div className="relative w-22 h-1.5 bg-white/20 rounded-full group/slider">
                    {/* Progress Bar */}
                    <div
                      className="absolute h-full bg-emerald-500 rounded-full"
                      style={{ width: `${baseUserVolume * 100}%` }}
                    />

                    {/* Thumb */}
                    <div
                      className="absolute h-2.5 w-2.5 bg-emerald-500 rounded-full top-1/2 -translate-x-1/2
                        -translate-y-1/2 shadow-[0_0_10px_rgba(16,185,129,0.4)] pointer-events-none transition-transform
                        group-hover/slider:scale-125"
                      style={{ left: `calc(${baseUserVolume * 100}% + (${(0.5 - baseUserVolume) * 10}px))` }}
                    />

                    {/* Invisible Input for control */}
                    <input
                      id="volume"
                      type="range"
                      min={0}
                      max={1}
                      step="any"
                      value={baseUserVolume}
                      onChange={handleVolumeChange}
                      className="absolute inset-0 w-full h-full opacity-0 cursor-pointer z-10"
                    />
                  </div>
                </div>
              </div>
              <div className="flex items-center gap-3">
                <button type="button" onClick={toggleSettings} className={'cursor-pointer group/settings'}>
                  <SettingsIcon
                    className={`w-6 h-6 transition-colors duration-300
                    ${showSettings ? 'text-emerald-500 rotate-45' : 'text-neutral-300 group-hover/settings:text-emerald-500'}`}
                  />
                </button>
                <button type="button" onClick={toggleFullScreen} className="cursor-pointer group/fullscreen">
                  <FullscreenIcon
                    className="text-neutral-300 w-6 h-6 group-hover/fullscreen:text-emerald-500 transition-colors"
                  />
                </button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
