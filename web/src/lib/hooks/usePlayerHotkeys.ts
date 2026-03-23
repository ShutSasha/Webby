import { useEffect } from 'react'

type UsePlayerHotkeysProps = {
  playerRef: React.MutableRefObject<any>
  hideTimeoutRef: React.MutableRefObject<NodeJS.Timeout | null>
  duration: number
  togglePlay: () => void
  setShowCustomControls: React.Dispatch<React.SetStateAction<boolean>>
  setState: React.Dispatch<React.SetStateAction<any>>
}

export const usePlayerHotkeys = ({
  playerRef,
  hideTimeoutRef,
  duration,
  togglePlay,
  setShowCustomControls,
  setState,
}: UsePlayerHotkeysProps) => {
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      const activeElement = document.activeElement
      if (
        activeElement?.tagName === 'INPUT' ||
        activeElement?.tagName === 'TEXTAREA' ||
        (activeElement as HTMLElement)?.isContentEditable
      ) {
        return
      }

      const player = playerRef.current as any
      if (!player) return

      const showControlsTemporarily = () => {
        setShowCustomControls(true)
        if (hideTimeoutRef.current) clearTimeout(hideTimeoutRef.current)
        hideTimeoutRef.current = setTimeout(() => setShowCustomControls(false), 3000)
      }

      switch (e.code) {
        case 'Space': {
          e.preventDefault()
          togglePlay()
          showControlsTemporarily()
          break
        }

        case 'ArrowLeft': {
          e.preventDefault()
          const currentTime = player.getCurrentTime ? player.getCurrentTime() : player.currentTime || 0
          const newTime = Math.max(0, currentTime - 5)

          if (player.seekTo) player.seekTo(newTime, 'seconds')
          else player.currentTime = newTime

          setState((prev: any) => ({
            ...prev,
            playedSeconds: newTime,
            played: duration > 0 ? newTime / duration : 0,
          }))
          showControlsTemporarily()
          break
        }

        case 'ArrowRight': {
          e.preventDefault()
          const currentTime = player.getCurrentTime ? player.getCurrentTime() : player.currentTime || 0
          const newTime = Math.min(duration, currentTime + 5)

          if (player.seekTo) player.seekTo(newTime, 'seconds')
          else player.currentTime = newTime

          setState((prev: any) => ({
            ...prev,
            playedSeconds: newTime,
            played: duration > 0 ? newTime / duration : 0,
          }))
          showControlsTemporarily()
          break
        }
      }
    }

    window.addEventListener('keydown', handleKeyDown)

    return () => {
      window.removeEventListener('keydown', handleKeyDown)
    }
  }, [duration, togglePlay, playerRef, hideTimeoutRef, setShowCustomControls, setState])
}
