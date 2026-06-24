'use client'

import ReactPlayer from 'react-player'

import { useCustomPlayerLogic } from '@/lib/hooks/useCustomPlayerLogic'
import { usePlayerControls } from '@/lib/hooks/usePlayerControls'
import { usePlayerHotkeys } from '@/lib/hooks/usePlayerHotkeys'
import { useSafeVideoTransition } from '@/lib/hooks/useSafeVideoTransition'
import { useVideoViewTracker } from '@/lib/hooks/useVideoViewTracker'

import PlayerBottomControls from './PlayerBottomControls'
import PlayerCenterButton from './PlayerCenterButton'
import PlayerLoader from './PlayerLoader'
import PlayerProgressBar from './PlayerProgressBar'
import PlayerSettingsMenu from './PlayerSettingsMenu'

type PlayerProps = {
  videoUrl: string
  videoId?: string
  roomId?: string
}

export default function CustomPlayer({ videoUrl, videoId, roomId }: PlayerProps) {
  const { isActuallyReady } = useSafeVideoTransition({
    videoUrl,
    delayMs: 800,
  })

  const isPlatformModeHook = usePlayerControls(videoUrl)
  const trackViewProgress = useVideoViewTracker(videoId)

  const isTwitch = videoUrl.includes('twitch.tv')
  const isPlatformMode = isPlatformModeHook || isTwitch

  const { refs, state, setState, uiState, actions } = useCustomPlayerLogic(
    videoUrl,
    isPlatformMode,
    trackViewProgress,
    roomId,
  )

  usePlayerHotkeys({
    playerRef: refs.playerRef,
    hideTimeoutRef: refs.hideTimeoutRef,
    duration: state.duration,
    togglePlay: actions.togglePlay,
    setShowCustomControls: actions.setShowCustomControls,
    setState,
  })

  if (!isActuallyReady) {
    return (
      <div
        className="aspect-video bg-neutral-100 dark:bg-surface-strong/80 w-full rounded-2xl flex items-center
          justify-center"
      >
        <div className="w-16 h-16 border-6 border-emerald-500/20 border-t-emerald-500 rounded-full animate-spin" />
      </div>
    )
  }

  return (
    <div
      ref={refs.playerContainerRef}
      onMouseMove={actions.handleMouseMove}
      onMouseLeave={actions.handleMouseLeave}
      className={`group relative w-full bg-black transition-all
        ${uiState.isFullScreen ? 'w-screen h-screen' : 'aspect-video'} ${!isTwitch ? 'rounded-2xl overflow-hidden' : ''}
        ${!isPlatformMode && !uiState.showCustomControls ? 'cursor-none' : 'cursor-default'}`}
    >
      <ReactPlayer
        className="react-player"
        ref={refs.playerRef}
        playing={uiState.playing}
        controls={isPlatformMode}
        width="100%"
        height="100%"
        light={state.light}
        muted={state.muted}
        loop={state.loop}
        volume={uiState.baseUserVolume}
        playbackRate={state.playbackRate}
        pip={state.pip}
        src={videoUrl}
        onRateChange={actions.handleRateChange}
        onProgress={actions.handleProgress}
        onTimeUpdate={actions.handleTimeUpdate}
        onDurationChange={actions.handleDurationChange}
        onVolumeChange={e => actions.handlePlayerVolumeChange(e, isPlatformMode)}
        onReady={actions.handleReactPlayerReady}
        onPlay={actions.handleReactPlayerPlay}
        onPause={actions.handleReactPlayerPause}
        onEnded={actions.handleReactPlayerEnded}
        onError={actions.handleReactPlayerError}
      />

      {state.error && (
        <div className="absolute inset-0 z-50 flex flex-col items-center justify-center bg-neutral-900/80 text-center
          px-4">
          <svg className="w-12 h-12 text-red-500 mb-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth="2"
              d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
            />
          </svg>
          <p className="text-neutral-100 font-semibold text-lg">Playback Error</p>
          <p className="text-foreground-muted text-sm mt-1">{state.error}</p>
        </div>
      )}

      {!state.error && !isPlatformMode && <PlayerLoader isReady={state.isReady} buffering={state.buffering} />}

      {/* Overlay */}
      {!state.error && state.isReady && !isPlatformMode && (
        <div
          className={`absolute inset-0 bg-linear-to-t from-black/80 via-black/10 to-transparent transition-opacity
          duration-400 ease-in-out ${uiState.showCustomControls ? 'opacity-100' : 'opacity-0'}`}
          onClick={actions.handlePlayPause}
        >
          {/* Play button on the overlay */}
          <PlayerCenterButton playing={uiState.playing} buffering={state.buffering} />

          {/* Settings Menu Popup */}
          <PlayerSettingsMenu
            showSettings={state.showSettings}
            playbackRate={state.playbackRate}
            onSetPlaybackRate={actions.handleSetPlaybackRate}
            menuRef={refs.settingsContainerRef}
          />

          {/* Controls */}
          <div
            onClick={e => e.stopPropagation()}
            className={`absolute z-2 bottom-2 left-2 right-2 md:bottom-3 md:left-3 md:right-3 flex flex-col gap-2
            transition-transform duration-500
            ${uiState.showCustomControls ? 'translate-y-0' : 'translate-y-10 pointer-events-none'}`}
          >
            {/* Progress Bar Container */}
            <PlayerProgressBar
              played={state.played}
              loaded={state.loaded}
              duration={state.duration}
              onSeekMouseDown={actions.handleSeekMouseDown}
              onSeekChange={actions.handleSeekChange}
              onSeekMouseUp={actions.handleSeekMouseUp}
            />

            {/* Buttons etc */}
            <PlayerBottomControls
              playing={uiState.playing}
              duration={state.duration}
              playedFraction={state.played}
              baseUserVolume={uiState.baseUserVolume}
              showSettings={state.showSettings}
              onPlayPause={actions.handlePlayPause}
              onVolumeChange={actions.handleVolumeChange}
              onToggleMute={actions.toggleMute}
              onToggleSettings={actions.toggleSettings}
              onToggleFullScreen={actions.toggleFullScreen}
            />
          </div>
        </div>
      )}
    </div>
  )
}
