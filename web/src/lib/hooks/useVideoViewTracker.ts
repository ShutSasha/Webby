import { useEffect, useRef, useCallback } from 'react'

import { incrementVideoView } from '@/lib/actions/video.actions'

export const useVideoViewTracker = (videoId: string | undefined) => {
  const hasCountedRef = useRef(false)

  const thresholdRef = useRef<number | null>(null)

  useEffect(() => {
    hasCountedRef.current = false
    thresholdRef.current = null
  }, [videoId])

  const trackProgress = useCallback(
    (playedSeconds: number, duration: number) => {
      if (!videoId || duration === 0 || hasCountedRef.current) return

      if (thresholdRef.current === null) {
        thresholdRef.current = Math.min(10, duration * 0.8)
      }

      if (playedSeconds >= thresholdRef.current) {
        hasCountedRef.current = true

        incrementVideoView(videoId).catch(error => {
          console.error('Failed to track view in background:', error)
        })
      }
    },
    [videoId],
  )

  return trackProgress
}
