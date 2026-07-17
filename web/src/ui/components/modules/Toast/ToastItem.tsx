'use client'
import { useEffect, useRef, useState } from 'react'

import { motion } from 'framer-motion'

import CheckCircle from '@/assets/icons/shared/check-circle.svg'
import XCircleIcon from '@/assets/icons/shared/circle-xmark.svg'
import InfoCircle from '@/assets/icons/shared/info-circle.svg'
import X from '@/assets/icons/shared/x.svg'
import { cn } from '@/lib/utils/general.utils'
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
    timerRef.current = setTimeout(() => removeToast(id), remaining)
  }

  const pauseTimer = () => {
    if (timerRef.current) {
      clearTimeout(timerRef.current)
      setRemaining(prev => prev - (Date.now() - startTimeRef.current))
    }
  }

  useEffect(() => {
    startTimer()
    return () => {
      if (timerRef.current) clearTimeout(timerRef.current)
    }
  }, [id])

  return (
    <motion.div
      layout
      initial={{ opacity: 0, y: 100 }}
      animate={{ opacity: 1, y: 0 }}
      exit={{
        opacity: 0,
        y: -40,
        scale: 1,
        transition: { duration: 0.3 },
      }}
      transition={{
        default: {
          type: 'tween',
          duration: 0.3,
          ease: 'circOut',
        },
      }}
      onMouseEnter={pauseTimer}
      onMouseLeave={startTimer}
      onClick={() => removeToast(id)}
      className={cn(
        'flex items-center justify-between p-4 rounded-xl border cursor-pointer select-none shadow-sm dark:shadow-none',
        'hover:scale-[1.01] active:scale-[0.98] transition-all duration-200',
        {
          [`bg-linear-to-r from-emerald-50 via-white to-white border-emerald-200 dark:from-emerald-950
          dark:via-neutral-900 dark:to-neutral-900 dark:border-emerald-500/50`]: type === 'success',
          [`bg-linear-to-r from-purple-50 via-white to-white border-purple-200 dark:from-purple-900 dark:via-neutral-900
          dark:to-neutral-900 dark:border-purple-500/50`]: type === 'info',
          [`bg-linear-to-r from-red-50 via-white to-white border-red-200 dark:from-red-950 dark:via-neutral-900
          dark:to-neutral-900 dark:border-red-500/50`]: type === 'error',
        },
      )}
    >
      <div className="flex items-center gap-3 pr-3">
        {type === 'success' && <CheckCircle className="size-5 text-emerald-600 dark:text-emerald-500" />}
        {type === 'info' && <InfoCircle className="size-5 text-purple-600 dark:text-purple-400" />}
        {type === 'error' && <XCircleIcon className="size-5 text-red-600 dark:text-red-500" />}

        <span className="text-sm font-medium text-neutral-800 dark:text-neutral-200">{message}</span>
      </div>

      <X
        className="size-4 text-neutral-400 hover:text-neutral-600 dark:text-neutral-500 dark:hover:text-neutral-300
          transition-colors"
      />
    </motion.div>
  )
}
