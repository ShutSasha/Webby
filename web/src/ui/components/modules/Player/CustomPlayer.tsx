'use client'
import { useEffect, useRef, useState } from 'react'

import ReactPlayer from 'react-player'

import { useIsClient } from '@/lib/hooks/useIsClient'
import { usePlayerControls } from '@/lib/hooks/usePlayerControls'
import { usePlayerHotkeys } from '@/lib/hooks/usePlayerHotkeys'
import { usePlayerPlayStore, usePlayerStore } from '@/stores/player.store'

import PlayerBottomControls from './PlayerBottomControls'
import PlayerCenterButton from './PlayerCenterButton'
import PlayerLoader from './PlayerLoader'
import PlayerProgressBar from './PlayerProgressBar'
import PlayerSettingsMenu from './PlayerSettingsMenu'

type PlayerProps = {
  videoUrl: string
}

export default function CustomPlayer({ videoUrl }: PlayerProps) {
  const playerRef = useRef<HTMLVideoElement>(null)
  const playerContainerRef = useRef<HTMLDivElement>(null)
  const isMounted = useIsClient()

  const settingsContainerRef = useRef<HTMLDivElement>(null)
  const [isFullScreen, setIsFullScreen] = useState(false)
  const [prevVolume, setPrevVolume] = useState(1)
  const hideTimeoutRef = useRef<NodeJS.Timeout | null>(null)
  const toggleThrottleRef = useRef<NodeJS.Timeout | null>(null)
  const [showCustomControls, setShowCustomControls] = useState(true)
  const isPlatformMode = usePlayerControls(videoUrl)

  const baseUserVolume = usePlayerStore(state => state.baseVolume)
  const setBaseUserVolume = usePlayerStore(state => state.setBaseVolume)

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
    buffering: true,
    isReady: false,
  }

  type PlayerState = typeof initialState

  const [state, setState] = useState<PlayerState>(initialState)

  const { light, muted, loop, played, loaded, duration, playbackRate, pip, showSettings, buffering, isReady } = state

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

    if (toggleThrottleRef.current) return

    togglePlay()

    toggleThrottleRef.current = setTimeout(() => {
      toggleThrottleRef.current = null
    }, 300)
  }

  const handleSetPlaybackRate = (rate: number) => {
    setState(prevState => ({ ...prevState, playbackRate: rate }))
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
    if (!player || state.seeking) return

    const buffered = player.buffered

    if (!buffered || buffered.length === 0) return

    const lastIndex = buffered.length - 1
    const loadedSeconds = buffered.end(lastIndex)

    setState(prevState => ({
      ...prevState,
      loadedSeconds,
      loaded: player.duration ? loadedSeconds / player.duration : 0,
    }))
  }

  const handleSeekMouseDown = () => {
    setState({ ...state, seeking: true })
  }

  const handleSeekChange = (event: React.SyntheticEvent<HTMLInputElement>) => {
    const inputTarget = event.target as HTMLInputElement
    setState(prevState => ({ ...prevState, played: Number.parseFloat(inputTarget.value) }))
  }

  const handleSeekMouseUp = (newTimeFraction: number) => {
    const player = playerRef.current
    if (!player) return
    const newTime = newTimeFraction * duration

    setState(prevState => ({
      ...prevState,
      seeking: false,
      played: newTimeFraction,
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
      onMouseLeave={() => playing && !showSettings && setShowCustomControls(false)}
      className={`group relative w-full bg-transparent overflow-hidden transition-all
        ${isFullScreen ? 'w-screen h-screen rounded-0' : 'aspect-video rounded-2xl'}
        ${!isPlatformMode && !showCustomControls ? 'cursor-none' : 'cursor-default'}`}
    >
      <ReactPlayer
        className="react-player"
        ref={playerRef}
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
        onReady={() => {
          const videoElement = playerRef.current

          if (videoElement) {
            const isActuallyLoaded = videoElement.readyState >= 3

            setState(prev => ({
              ...prev,
              buffering: !isActuallyLoaded,
              isReady: isActuallyLoaded,
            }))

            videoElement.onwaiting = () => setState(prev => ({ ...prev, buffering: true }))
            videoElement.onloadeddata = () => setState(prev => ({ ...prev, buffering: false, isReady: true }))
            videoElement.onplaying = () => setState(prev => ({ ...prev, buffering: false, isReady: true }))
            videoElement.oncanplay = () => setState(prev => ({ ...prev, buffering: false, isReady: true }))
          }
        }}
        onPlay={() => {
          if (playerRef.current?.paused) return

          setPlaying(true)
          setState(prev => ({ ...prev, buffering: false }))
        }}
        onPause={() => {
          if (!playerRef.current?.paused) return

          setPlaying(false)
        }}
        onEnded={() => {
          triggerEnded()
        }}
        onError={(e: any) => {
          if (e?.name === 'AbortError') {
            console.warn('Play interrupted safely. (AbortError)')
            setPlaying(false)
            return
          }

          if (e?.name === 'NotAllowedError') {
            console.warn('Autoplay prevented by browser. User must click play.')
            setPlaying(false)
            return
          }

          const target = e?.target as HTMLVideoElement | undefined
          if (target?.error) {
            console.error('Video Media Error. Code:', target.error.code, 'Message:', target.error.message)
            return
          }

          console.error('Unhandled ReactPlayer Error:', e)
        }}
      />

      <PlayerLoader isReady={isReady} buffering={buffering} />

      {/* Overlay */}
      {isReady && !isPlatformMode && (
        <div
          className={`absolute inset-0 bg-linear-to-t from-black/80 via-black/10 to-transparent transition-opacity
          duration-400 ease-in-out ${showCustomControls ? 'opacity-100' : 'opacity-0'}`}
          onClick={handlePlayPause}
        >
          {/* Play button on the overlay */}
          <PlayerCenterButton playing={playing} buffering={buffering} />

          {/* Settings Menu Popup */}
          <PlayerSettingsMenu
            showSettings={showSettings}
            playbackRate={playbackRate}
            onSetPlaybackRate={handleSetPlaybackRate}
            menuRef={settingsContainerRef}
          />

          {/* Controls */}
          <div
            onClick={e => e.stopPropagation()}
            className={`absolute z-2 bottom-3 left-3 right-3 flex flex-col gap-2 transition-transform duration-500
            ${showCustomControls ? 'translate-y-0' : 'translate-y-10 pointer-events-none'}`}
          >
            {/* Progress Bar Container */}
            <PlayerProgressBar
              played={played}
              loaded={loaded}
              duration={duration}
              onSeekMouseDown={handleSeekMouseDown}
              onSeekChange={handleSeekChange}
              onSeekMouseUp={handleSeekMouseUp}
            />

            {/* Buttons etc */}
            <PlayerBottomControls
              playing={playing}
              duration={duration}
              playedFraction={played}
              baseUserVolume={baseUserVolume}
              showSettings={showSettings}
              onPlayPause={handlePlayPause}
              onVolumeChange={handleVolumeChange}
              onToggleMute={toggleMute}
              onToggleSettings={toggleSettings}
              onToggleFullScreen={toggleFullScreen}
            />
          </div>
        </div>
      )}
    </div>
  )
}
