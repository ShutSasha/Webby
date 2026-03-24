'use client'

import { useState, useTransition } from 'react'

import { updateAboutField } from '@/app/api/user'
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

      if (!result.success && result?.errors) {
        for (const [, value] of Object.entries(result.errors)) {
          addToast(value, 'error')
        }
      }
    })
  }

  return (
    <div className="w-full">
      <textarea
        value={text}
        onChange={e => setText(e.target.value)}
        placeholder="Type something about yourself"
        disabled={isPending}
        className="w-full rounded-[20px] p-4 border border-border bg-transparent ring-0 outline-0 h-[200px] resize-none
          text-start align-top mb-4 transition-all focus:border-emerald-500 disabled:opacity-50"
      />

      <div className="flex items-center justify-center">
        <Button
          viewType={!isPending ? 'confirm' : 'loading'}
          className={'rounded-xl font-medium min-w-[120px]'}
          onClick={handleSave}
          disabled={isPending}
        >
          {isPending ? 'Saving...' : 'Save'}
        </Button>
      </div>
    </div>
  )
}
