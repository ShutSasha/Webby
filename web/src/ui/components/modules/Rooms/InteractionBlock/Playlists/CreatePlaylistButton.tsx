'use client'
import { ChangeEvent, FormEvent, useState, useTransition } from 'react'

import { useRouter } from 'next/navigation'

import { createUserPlaylist } from '@/app/api/playlists'
import EditPenIcon from '@/assets/icons/shared/edit-pen.svg'
import { extractServerMessage, serverLog } from '@/lib/utils/utils'
import { useToastStore } from '@/stores/toast-store'
import MediaButton from '@/ui/components/common/MediaButton'
import Input from '@/ui/components/Input'
import Button from '@/ui/components/shared/Button'
import Switch from '@/ui/components/shared/Switch'

export default function CreatePlaylistButton() {
  const addToast = useToastStore(state => state.addToast)
  const [isModalOpen, setIsModalOpen] = useState<boolean>(false)
  const [isPrivate, setIsPrivate] = useState<boolean>(false)
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [name, setName] = useState<string>('')
  const router = useRouter()

  const [isPending, startTransition] = useTransition()

  const togglePrivate = () => setIsPrivate(prev => !prev)
  const handleInput = (e: ChangeEvent<HTMLInputElement>) => setName(e.target.value)

  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault()

    const trimmedName = name.trim()
    if (!trimmedName) {
      addToast(`Playlist's name must be not empty`, 'info')
      return
    }

    try {
      setIsSubmitting(true)
      const response = await createUserPlaylist(trimmedName, isPrivate)

      if (response.success) {
        addToast('User playlist has been created successfully', 'success')

        setIsModalOpen(false)
        setName('')
        setIsPrivate(false)

        startTransition(() => {
          router.refresh()
        })
      } else {
        const msg = extractServerMessage(response.errors)
        addToast(msg ?? `Something went wrong while creating ${name} playlist`, 'error')
      }
    } catch (error: unknown) {
      serverLog('Handle submit for create user playlist', error, true)
    } finally {
      setIsSubmitting(false)
    }
  }

  const isSubmitDisabled = isSubmitting || isPending

  return (
    <MediaButton actionLabel="Create a playlist" isOpen={isModalOpen} setIsOpen={setIsModalOpen}>
      <form className="flex flex-col" onSubmit={handleSubmit}>
        <h3 className="text-neutral-300 text-center mb-4 font-semibold text-xl">Create a new playlist</h3>

        <div className="relative mb-4">
          <label htmlFor="playlist-name" className="sr-only">
            Playlist Name
          </label>
          <Input
            id="playlist-name"
            placeholder="Enter playlist name..."
            className="py-2.5 pl-10 rounded-xl w-full border-neutral-700"
            value={name}
            onChange={handleInput}
          />
          <EditPenIcon className="absolute left-3 top-1/2 -translate-y-1/2 size-4 text-neutral-700" />
        </div>

        <div className="flex items-center justify-between mb-6">
          <p className="text-sm text-neutral-300 font-medium">Private:</p>

          <Switch isChecked={isPrivate} toggle={togglePrivate} />
        </div>

        <div className="flex justify-center">
          <Button
            type="submit"
            disabled={isSubmitDisabled}
            viewType={isSubmitDisabled ? 'loading' : 'confirm'}
            className="rounded-xl px-8 py-2 font-semibold"
          >
            Create
          </Button>
        </div>
      </form>
    </MediaButton>
  )
}
