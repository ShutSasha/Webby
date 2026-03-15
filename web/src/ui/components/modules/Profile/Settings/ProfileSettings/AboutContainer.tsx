'use client'

import { useToastStore } from '@/stores/toast-store'
import Button from '@/ui/components/shared/Button'

type Props = {
  about: string
}

export default function AboutContainer({ about }: Props) {
  // useActionState
  const addToast = useToastStore(state => state.addToast)

  const handleClick = () => {
    addToast('About field updated', 'info')
  }

  return (
    <>
      <textarea
        defaultValue={about}
        placeholder="Type something about yourself"
        className="w-full rounded-[20px] p-4 border border-border bg-transparent ring-0 outline-0 h-[200px] resize-none
          text-start align-top mb-4"
      />
      <div className="flex items-center justify-center">
        <Button viewType="confirm" className="rounded-xl font-medium" onClick={handleClick}>
          Save
        </Button>
      </div>
    </>
  )
}
