'use client'

import { useState, useTransition } from 'react'

import { updateAboutField } from '@/lib/actions/user.actions'
import { useToastStore } from '@/stores/toast-store'
import Button from '@/ui/components/shared/Button'

type Props = {
  userId: string
  about: string
}

export default function AboutContainer({ userId, about }: Props) {
  const [text, setText] = useState(about)
  const [isPending, startTransition] = useTransition()
  const addToast = useToastStore(state => state.addToast)

  const handleSave = () => {
    if (text === about) {
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
        className="w-full rounded-2xl p-5 border border-neutral-800/80 bg-neutral-950/50 h-40 resize-none text-start
          align-top mb-5 text-neutral-200 placeholder:text-neutral-600 outline-none transition-all duration-300
          hover:border-neutral-600 focus:border-emerald-500/50 focus:bg-neutral-900/80 disabled:opacity-50
          disabled:cursor-not-allowed"
      />

      <div className="flex justify-end">
        <Button
          viewType={!isPending ? 'confirm' : 'loading'}
          className="rounded-xl font-medium min-w-[140px]"
          onClick={handleSave}
          disabled={isPending || text === about}
        >
          {isPending ? 'Saving...' : 'Save changes'}
        </Button>
      </div>
    </div>
  )
}
