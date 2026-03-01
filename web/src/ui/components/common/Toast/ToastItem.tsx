'use client'
import { useEffect, useRef, useState } from 'react'

import { clsx, type ClassValue } from 'clsx'
import { CheckCircle, XCircle, X } from 'lucide-react'
import { twMerge } from 'tailwind-merge'

import { ToastType, useToastStore } from '@/stores/toast-store'

interface ToastItemProps {
  id: string
  message: string
  type: ToastType
}

export default function ToastItem({ id, message, type }: ToastItemProps) {
  const removeToast = useToastStore(state => state.removeToast)
  const [remaining, setRemaining] = useState(4000)
  const timerRef = useRef<NodeJS.Timeout | null>(null)
  const startTimeRef = useRef(Date.now())

  const startTimer = () => {
    startTimeRef.current = Date.now()
    timerRef.current = setTimeout(() => {
      removeToast(id)
    }, remaining)
  }

  const pauseTimer = () => {
    if (timerRef.current) {
      clearTimeout(timerRef.current)

      setRemaining(prev => {
        return prev - (Date.now() - startTimeRef.current)
      })
    }
  }

  useEffect(() => {
    startTimer()
    return () => {
      if (timerRef.current) clearTimeout(timerRef.current)
    }
  }, [id])

  function cn(...inputs: ClassValue[]) {
    return twMerge(clsx(inputs))
  }

  return (
    <div
      onMouseEnter={pauseTimer}
      onMouseLeave={startTimer}
      onClick={() => removeToast(id)}
      className={cn(
        'flex items-center justify-between p-4 rounded-xl border cursor-pointer select-none',
        'duration-300 transition-all hover:scale-[1.02] active:scale-[0.98]',
        {
          'bg-linear-to-r from-red-950 via-neutral-900 to-neutral-900 border-red-500/50': type === 'error',
          'bg-linear-to-r from-emerald-950 via-neutral-900 to-neutral-900 border-emerald-500/50': type === 'success',
        },
      )}
    >
      <div className="flex items-center gap-3">
        {type === 'success' && <CheckCircle size={20} className="text-emerald-700" />}
        {type === 'error' && <XCircle size={20} className="text-red-500" />}
        <span className="text-sm font-medium text-neutral-300">{message}</span>
      </div>
      <X size={16} className="text-neutral-500 hover:text-neutral-300 transition-colors" />
    </div>
  )
}
