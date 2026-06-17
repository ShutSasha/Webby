import { useState, useEffect } from 'react'

export const useCountdown = (expiresAt?: string, duration?: number) => {
  const [timeLeftMs, setTimeLeftMs] = useState(0)
  const [progress, setProgress] = useState(0)
  const [isExpired, setIsExpired] = useState(false)

  useEffect(() => {
    if (!expiresAt || !duration) {
      setTimeLeftMs(0)
      setProgress(0)
      setIsExpired(true)
      return
    }

    const targetTime = new Date(expiresAt).getTime()
    const totalDurationMs = duration * 1000

    const updateTimer = () => {
      const now = Date.now()
      const remaining = Math.max(0, targetTime - now)

      setTimeLeftMs(remaining)

      setProgress(Math.max(0, Math.min(100, (remaining / totalDurationMs) * 100)))

      if (remaining === 0) {
        setIsExpired(true)
      } else {
        setIsExpired(false)
      }

      return remaining
    }

    updateTimer()

    let animationFrameId: number

    const tick = () => {
      const remaining = updateTimer()
      if (remaining > 0) {
        animationFrameId = requestAnimationFrame(tick)
      }
    }

    animationFrameId = requestAnimationFrame(tick)

    return () => cancelAnimationFrame(animationFrameId)
  }, [expiresAt, duration])

  const secondsLeft = Math.ceil(timeLeftMs / 1000)

  return { timeLeftMs, secondsLeft, progress, isExpired }
}
