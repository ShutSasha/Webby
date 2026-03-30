'use client'
import { ChangeEvent, FormEvent, useState } from 'react'

import EditPenIcon from '@/assets/icons/shared/edit-pen.svg'
import { useCreatePlaylist } from '@/lib/hooks/api/playlist/useCreatePlaylist'
import Button from '@/ui/components/shared/Button'
import Input from '@/ui/components/shared/Input'
import MediaButton from '@/ui/components/shared/MediaButton'
import Switch from '@/ui/components/shared/Switch'

export default function CreatePlaylistButton() {
  const [isModalOpen, setIsModalOpen] = useState<boolean>(false)
  const [isPrivate, setIsPrivate] = useState<boolean>(false)
  const [name, setName] = useState<string>('')

  const togglePrivate = () => setIsPrivate(prev => !prev)
  const handleInput = (e: ChangeEvent<HTMLInputElement>) => setName(e.target.value)

  const { handleCreate, isLoading } = useCreatePlaylist({
    onSuccess: () => {
      setIsModalOpen(false)
      setName('')
      setIsPrivate(false)
    },
  })

  const handleSubmit = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    handleCreate(name, isPrivate)
  }

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
            disabled={isLoading}
            viewType={isLoading ? 'loading' : 'confirm'}
            className="rounded-xl px-8 py-2 font-semibold"
          >
            Create
          </Button>
        </div>
      </form>
    </MediaButton>
  )
}
