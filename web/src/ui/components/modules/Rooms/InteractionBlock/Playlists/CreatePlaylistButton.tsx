'use client'
import { ChangeEvent, MouseEvent, useState } from 'react'

import { useRouter } from 'next/navigation'

import { createUserPlaylist } from '@/app/api/playlists'
import EditPenIcon from '@/assets/icons/shared/edit-pen.svg'
import { extractServerMessage, serverLog } from '@/lib/utils/utils'
import { useToastStore } from '@/stores/toast-store'
import MediaButton from '@/ui/components/common/MediaButton'
import Input from '@/ui/components/Input'
import Button from '@/ui/components/shared/Button'

export default function CreatePlaylistButton() {
  const addToast = useToastStore(state => state.addToast)
  const [isModalOpen, setIsModalOpen] = useState<boolean>(false)
  const [isPrivate, setIsPrivate] = useState<boolean>(false)
  const [loading, setLoading] = useState<boolean>(false)
  const [name, setName] = useState<string>('')
  const router = useRouter()

  const togglePrivate = () => setIsPrivate(prev => !prev)
  const handleInput = (e: ChangeEvent<HTMLInputElement>) => setName(e.target.value)

  const handleSubmit = async (e: MouseEvent<HTMLButtonElement>) => {
    e.preventDefault()
    try {
      setLoading(true)
      const response = await createUserPlaylist(name, isPrivate)

      if (response.success) {
        addToast('User playlist has been created successfully', 'success')

        setIsModalOpen(false)
        setName('')
        setIsPrivate(false)
        router.refresh()
      } else {
        const msg = extractServerMessage(response.errors)
        addToast(msg ?? `Something went wrong while creating ${name} playlist`, 'error')
      }
    } catch (error: unknown) {
      serverLog('Handle submit for create user playlist', error, true)
    } finally {
      setLoading(false)
    }
  }

  return (
    <MediaButton actionLabel="Create a playlist" isOpen={isModalOpen} setIsOpen={setIsModalOpen}>
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

        <button
          type="button"
          role="switch"
          aria-checked={isPrivate}
          onClick={togglePrivate}
          className={` relative inline-flex p-1 w-16 cursor-pointer items-center rounded-full transition-colors
            duration-200 ease-in-out focus:outline-none ${isPrivate ? 'bg-emerald-500' : 'bg-neutral-800'}
            hover:bg-opacity-80 `}
        >
          <span
            className={` inline-block size-5 transform rounded-full shadow-lg transition duration-200 ease-in-out
              ${isPrivate ? 'translate-x-9 bg-neutral-900' : 'translate-x-0 bg-neutral-700'} `}
          />
        </button>
      </div>

      <div className="flex justify-center">
        <Button
          viewType={loading ? 'loading' : 'confirm'}
          className="rounded-xl px-8 py-2 font-semibold"
          onClick={handleSubmit}
        >
          Create
        </Button>
      </div>
    </MediaButton>
  )
}
