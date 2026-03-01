'use client'
import { useToastStore } from '@/stores/toast-store'

import ToastItem from './ToastItem'

export default function ToastContainer() {
  const toasts = useToastStore(state => state.toasts)

  return (
    <div className="fixed bottom-8 left-1/2 -translate-x-1/2 z-100 flex flex-col gap-2 w-full max-w-md px-4">
      {toasts.map(toast => (
        <ToastItem key={toast.id} {...toast} />
      ))}
    </div>
  )
}
