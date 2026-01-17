'use client'
import { useCallback, useRef, useState } from 'react'

import ReactPlayer from 'react-player'

import PauseIcon from '@/assets/icons/Player/pause.svg'
import PlayIcon from '@/assets/icons/Player/play.svg'

import Duration from './Duration'

type PlayerProps = {
  videoUrl: string
}

export default function CustomPlayer({ videoUrl }: PlayerProps) {
  const playerRef = useRef<HTMLVideoElement | null>(null)

  const setPlayerRef = useCallback((player: HTMLVideoElement) => {
    if (!player) return
    playerRef.current = player
  }, [])

  const initialState = {
    src: videoUrl,
    pip: false,
    playing: false,
    controls: false,
    light: false,
    volume: 1,
    muted: false,
    played: 0,
    loaded: 0,
    duration: 0,
    playbackRate: 1.0,
    loop: false,
    seeking: false,
    loadedSeconds: 0,
    playedSeconds: 0,
  }

  type PlayerState = Omit<typeof initialState, 'src'> & {
    src?: string
  }

  const [state, setState] = useState<PlayerState>(initialState)

  const { src, playing, controls, light, volume, muted, loop, played, loaded, duration, playbackRate, pip } = state

  const handlePlayPause = () => {
    setState(prevState => ({ ...prevState, playing: !prevState.playing }))
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
    setState(prevState => ({ ...prevState, seeking: false }))
    if (playerRef.current) {
      playerRef.current.currentTime = Number.parseFloat(inputTarget.value) * playerRef.current.duration
    }
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

  return (
    <div className="group relative aspect-video w-full bg-transparent rounded-2xl overflow-hidden">
      <ReactPlayer
        ref={setPlayerRef}
        playing={playing}
        controls={controls}
        width="100%"
        height="100%"
        src={videoUrl}
        onProgress={handleProgress}
        onTimeUpdate={handleTimeUpdate}
        onDurationChange={handleDurationChange}
      />

      {/* Overlay */}
      <div
        className="absolute inset-0 bg-linear-to-t from-black/80 via-transparent to-transparent opacity-0
          group-hover:opacity-100 transition-opacity"
        onClick={handlePlayPause}
      >
        {/* Controls */}
        <div onClick={e => e.stopPropagation()} className="absolute bottom-2 left-2 right-2 flex flex-col gap-2">
          {/* Progress Bar */}
          <div className="relative h-1.5 w-full bg-white/20 rounded-full group/bar">
            {/* Pre-Loaded */}
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
          <div className="flex items-center gap-3">
            <button onClick={handlePlayPause} className="cursor-pointer">
              {playing ? (
                <PauseIcon className="w-6 h-6 text-neutral-300" />
              ) : (
                <PlayIcon className="w-6 h-6 text-neutral-300" />
              )}
            </button>
            <div className="flex items-center gap-0.5">
              <Duration seconds={duration * played} className="text-neutral-300 text-sm font-medium leading-5" />
              <span className="text-neutral-300 text-sm font-medium leading-5">/</span>
              <Duration seconds={duration} className="text-neutral-300 text-sm font-medium leading-5" />
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
