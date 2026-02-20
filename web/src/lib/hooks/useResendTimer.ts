import { useEffect, useState } from 'react'

export function useResendTimer(email: string | null, initialSeconds = 60) {
  const [timeLeft, setTimeLeft] = useState(0)
  const KEY = `verify_timer_${email}`

  useEffect(() => {
    const savedExpiry = localStorage.getItem(KEY)
    if (savedExpiry) {
      const remaining = Math.floor((parseInt(savedExpiry) - Date.now()) / 1000)
      if (remaining > 0) {
        setTimeLeft(remaining)
      } else {
        localStorage.removeItem(KEY)
      }
    }
  }, [email, KEY])

  useEffect(() => {
    if (timeLeft <= 0) return
    const timer = setInterval(() => {
      setTimeLeft(prev => {
        if (prev <= 1) {
          localStorage.removeItem(KEY)
          return 0
        }
        return prev - 1
      })
    }, 1000)
    return () => clearInterval(timer)
  }, [timeLeft, KEY])

  const start = () => {
    localStorage.setItem(KEY, (Date.now() + initialSeconds * 1000).toString())
    setTimeLeft(initialSeconds)
  }

  return { timeLeft, start }
}
