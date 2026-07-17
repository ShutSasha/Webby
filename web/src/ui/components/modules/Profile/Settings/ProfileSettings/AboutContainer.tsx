'use client'

import { useState, useTransition } from 'react'

import { updateAboutField } from '@/lib/actions/user.actions'
import { cn } from '@/lib/utils/general.utils'
import { useToastStore } from '@/stores/toast-store'

type Props = {
  userId: string
  about: string
}

export default function AboutContainer({ userId, about }: Props) {
  const [text, setText] = useState(about)
  const [isPending, startTransition] = useTransition()
  const addToast = useToastStore(state => state.addToast)

  const isChanged = text !== about

  const handleSave = () => {
    if (!isChanged) {
      addToast('No changes to save', 'info')
      return
    }

    startTransition(async () => {
      const result = await updateAboutField(userId, text)

      if (result?.success) {
        addToast('Profile information updated successfully!', 'success')
      }

      if (!result?.success && result?.errors) {
        for (const [, value] of Object.entries(result.errors)) {
          addToast(value, 'error')
        }
      }
    })
  }

  return (
    <div className="w-full flex flex-col">
      <textarea
        value={text}
        onChange={e => setText(e.target.value)}
        placeholder="Type something about yourself..."
        disabled={isPending}
        className="w-full rounded-2xl p-5 border border-border/60 bg-surface-secondary/30 h-40 resize-none text-start
          align-top mb-5 text-foreground placeholder:text-foreground-muted outline-none transition-all duration-300
          hover:border-border focus:border-emerald-500/50 focus:ring-4 focus:ring-emerald-500/10 focus:bg-surface
          disabled:opacity-50 disabled:cursor-not-allowed custom-scrollbar"
      />

      <div className="flex justify-center">
        <button
          disabled={isPending || !isChanged}
          onClick={handleSave}
          className={cn(
            'px-8 py-2.5 rounded-xl font-medium transition-all duration-300 min-w-40 border',
            isChanged
              ? `bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 hover:bg-emerald-500/20 border-emerald-500/30
                hover:border-emerald-500/50 shadow-sm`
              : 'bg-surface-secondary text-foreground-muted cursor-not-allowed opacity-70 border-border/50 shadow-none',
          )}
        >
          {isPending ? 'Saving...' : 'Save changes'}
        </button>
      </div>
    </div>
  )
}
