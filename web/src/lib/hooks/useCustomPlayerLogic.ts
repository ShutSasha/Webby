import { useEffect, useRef, useState } from 'react'

import { usePlayerPlayStore, usePlayerStore } from '@/stores/player.store'

type PlayerState = {
  pip: boolean
  light: boolean
  muted: boolean
  played: number
  loaded: number
  duration: number
  playbackRate: number
  loop: boolean
  seeking: boolean
  loadedSeconds: number
  playedSeconds: number
  showSettings: boolean
  buffering: boolean
  isReady: boolean
}

export const useCustomPlayerLogic = (
  videoUrl: string,
  isPlatformMode: boolean,
  trackViewProgress?: (playedSeconds: number, duration: number) => void,
) => {
  const playerRef = useRef<HTMLVideoElement>(null)
  const playerContainerRef = useRef<HTMLDivElement>(null)
  const settingsContainerRef = useRef<HTMLDivElement>(null)

  const [isFullScreen, setIsFullScreen] = useState(false)
  const [prevVolume, setPrevVolume] = useState(1)
  const [showCustomControls, setShowCustomControls] = useState(true)

  const hideTimeoutRef = useRef<NodeJS.Timeout | null>(null)
  const toggleThrottleRef = useRef<NodeJS.Timeout | null>(null)

  const baseUserVolume = usePlayerStore(state => state.baseVolume)
  const setBaseUserVolume = usePlayerStore(state => state.setBaseVolume)

  const triggerEnded = usePlayerPlayStore(state => state.triggerEnded)
  const togglePlay = usePlayerPlayStore(state => state.togglePlay)
  const setPlaying = usePlayerPlayStore(state => state.setPlaying)
  const playing = usePlayerPlayStore(state => state.playing)

  const initialState: PlayerState = {
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

  const [state, setState] = useState<PlayerState>(initialState)

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

  const handlePlayPause = (e: React.MouseEvent) => {
    if (state.showSettings) {
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
    const newTime = newTimeFraction * state.duration

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
    if (!player || state.seeking) return
    if (!player.duration) return

    const currentSeconds = player.currentTime
    const currentDuration = player.duration

    setState(prevState => ({
      ...prevState,
      playedSeconds: currentSeconds,
      played: currentSeconds / currentDuration,
    }))

    if (trackViewProgress) {
      trackViewProgress(currentSeconds, currentDuration)
    }
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

    if (playing && !state.showSettings) {
      hideTimeoutRef.current = setTimeout(() => {
        setShowCustomControls(false)
      }, 3000)
    }
  }

  const handleMouseLeave = () => {
    if (playing && !state.showSettings) {
      setShowCustomControls(false)
    }
  }

  const handleReactPlayerVolumeChange = (e: any) => {
    if (typeof e === 'number') {
      setBaseUserVolume(e)
      return
    }

    const volume = e.target?.volume
    if (typeof volume === 'number') {
      setBaseUserVolume(volume)
    }
  }

  const handleReactPlayerReady = () => {
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
  }

  const handleReactPlayerPlay = () => {
    if (playerRef.current?.paused) return
    setPlaying(true)
    setState(prev => ({ ...prev, buffering: false }))
  }

  const handleReactPlayerPause = () => {
    if (!playerRef.current?.paused) return
    setPlaying(false)
  }

  const handleReactPlayerEnded = () => {
    triggerEnded()
  }

  const handleReactPlayerError = (e: any) => {
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
  }

  return {
    refs: {
      playerRef,
      playerContainerRef,
      settingsContainerRef,
      hideTimeoutRef,
    },
    state,
    setState,
    uiState: {
      isFullScreen,
      showCustomControls,
      playing,
      baseUserVolume,
    },
    actions: {
      setShowCustomControls,
      handlePlayPause,
      handleSetPlaybackRate,
      handleRateChange,
      toggleSettings,
      handleProgress,
      handleSeekMouseDown,
      handleSeekChange,
      handleSeekMouseUp,
      handleTimeUpdate,
      handleDurationChange,
      handleVolumeChange,
      toggleMute,
      toggleFullScreen,
      handleMouseMove,
      handleMouseLeave,
      handleReactPlayerVolumeChange,
      handleReactPlayerReady,
      handleReactPlayerPlay,
      handleReactPlayerPause,
      handleReactPlayerEnded,
      handleReactPlayerError,
      togglePlay,
    },
  }
}
