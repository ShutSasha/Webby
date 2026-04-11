'use client'

import ReactPlayer from 'react-player'

import { useCustomPlayerLogic } from '@/lib/hooks/useCustomPlayerLogic'
import { useIsClient } from '@/lib/hooks/useIsClient'
import { usePlayerControls } from '@/lib/hooks/usePlayerControls'
import { usePlayerHotkeys } from '@/lib/hooks/usePlayerHotkeys'
import { useVideoViewTracker } from '@/lib/hooks/useVideoViewTracker'

import PlayerBottomControls from './PlayerBottomControls'
import PlayerCenterButton from './PlayerCenterButton'
import PlayerLoader from './PlayerLoader'
import PlayerProgressBar from './PlayerProgressBar'
import PlayerSettingsMenu from './PlayerSettingsMenu'

type PlayerProps = {
  videoUrl: string
  videoId?: string
  isRoom?: boolean
}

export default function CustomPlayer({ videoUrl, videoId, isRoom = false }: PlayerProps) {
  const isMounted = useIsClient()
  const isPlatformMode = usePlayerControls(videoUrl)

  const trackViewProgress = useVideoViewTracker(videoId)

  const { refs, state, setState, uiState, actions } = useCustomPlayerLogic(videoUrl, isPlatformMode, trackViewProgress)

  usePlayerHotkeys({
    playerRef: refs.playerRef,
    hideTimeoutRef: refs.hideTimeoutRef,
    duration: state.duration,
    togglePlay: actions.togglePlay,
    setShowCustomControls: actions.setShowCustomControls,
    setState,
  })

  if (!isMounted)
    return (
      <div className="aspect-video bg-black w-full rounded-2xl flex items-center justify-center">
        <div className="w-12 h-12 border-4 border-emerald-500/20 border-t-emerald-500 rounded-full animate-spin" />
      </div>
    )

  return (
    <div
      ref={refs.playerContainerRef}
      onMouseMove={actions.handleMouseMove}
      onMouseLeave={actions.handleMouseLeave}
      className={`group relative w-full bg-transparent overflow-hidden transition-all
        ${uiState.isFullScreen ? 'w-screen h-screen rounded-0' : 'aspect-video rounded-2xl'}
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
        onVolumeChange={actions.handleReactPlayerVolumeChange}
        onReady={actions.handleReactPlayerReady}
        onPlay={actions.handleReactPlayerPlay}
        onPause={actions.handleReactPlayerPause}
        onEnded={actions.handleReactPlayerEnded}
        onError={actions.handleReactPlayerError}
      />

      <PlayerLoader isReady={state.isReady} buffering={state.buffering} />

      {/* Overlay */}
      {state.isReady && !isPlatformMode && (
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
              isRoom={isRoom}
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
