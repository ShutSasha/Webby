'use client'

import { ChangeEvent, FormEvent, useMemo, useRef, useState } from 'react'

import EditPenIcon from '@/assets/icons/shared/edit-pen.svg'
import { useCategoriesQuery } from '@/lib/hooks/api/category/useCategoriesQuery'
import { useCreateRoom } from '@/lib/hooks/api/room/useCreateRoom'
import { useToastStore } from '@/stores/toast-store'
import Button from '@/ui/components/shared/Button'
import Input from '@/ui/components/shared/Input'
import MediaButton from '@/ui/components/shared/MediaButton'
import Switch from '@/ui/components/shared/Switch'

export default function CreateRoomButton() {
  const addToast = useToastStore(state => state.addToast)
  const [isModalOpen, setIsModalOpen] = useState<boolean>(false)
  const fileInputRef = useRef<HTMLInputElement>(null)

  const [name, setName] = useState<string>('')
  const [categoryName, setCategoryName] = useState<string>('')
  const [isPrivate, setIsPrivate] = useState<boolean>(false)
  const [thumbnail, setThumbnail] = useState<File | null>(null)

  const { data: categoriesData, isLoading: isLoadingCategories } = useCategoriesQuery('')

  const categories = useMemo(() => {
    return categoriesData?.pages.flatMap(page => page?.data?.items || []) || []
  }, [categoriesData])

  const togglePrivate = () => setIsPrivate(prev => !prev)

  const { mutate: createRoom, isPending } = useCreateRoom({
    onSuccess: () => {
      setIsModalOpen(false)
      setName('')
      setCategoryName('')
      setIsPrivate(false)
      setThumbnail(null)
      if (fileInputRef.current) {
        fileInputRef.current.value = ''
      }
    },
  })

  const handleSubmit = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault()

    if (!name || name.trim().length < 2) {
      addToast('Room name must contain at least 2 characters', 'error')
      return
    }
    if (!categoryName) {
      addToast('Room category must be selected', 'error')
      return
    }

    const formData = new FormData()
    formData.append('name', name)
    formData.append('categoryName', categoryName)
    formData.append('isPrivate', isPrivate ? 'true' : 'false')

    if (thumbnail) {
      formData.append('thumbnail', thumbnail)
    }

    createRoom(formData)
  }

  return (
    <MediaButton isOpen={isModalOpen} setIsOpen={setIsModalOpen} actionLabel="Create Room">
      <form className="flex flex-col" onSubmit={handleSubmit}>
        <h3 className="text-neutral-300 text-center mb-6 font-semibold text-xl">Create a new room</h3>

        <div className="relative mb-4">
          <label htmlFor="room-name" className="sr-only">
            Room Name
          </label>
          <Input
            id="room-name"
            placeholder="Enter room name..."
            className="py-2.5 pl-10 rounded-xl w-full border-neutral-700"
            value={name}
            onChange={(e: ChangeEvent<HTMLInputElement>) => setName(e.target.value)}
            maxLength={50}
          />
          <EditPenIcon className="absolute left-3 top-1/2 -translate-y-1/2 size-4 text-neutral-700" />
        </div>

        <div className="mb-4">
          {/* TODO: change default html select tag to custom component  */}
          <select
            value={categoryName}
            onChange={e => setCategoryName(e.target.value)}
            disabled={isLoadingCategories}
            className="py-2.5 px-4 rounded-xl w-full border border-neutral-700 bg-neutral-900 text-neutral-200
              outline-none focus:border-emerald-500 transition-colors appearance-none cursor-pointer disabled:opacity-50
              disabled:cursor-not-allowed"
          >
            <option value="" disabled>
              Select category...
            </option>

            {categories.map((category: string) => {
              return (
                <option key={category} value={category}>
                  {category}
                </option>
              )
            })}
          </select>
        </div>

        <div className="flex flex-col gap-1 mb-6">
          <label className="text-sm text-neutral-300 font-medium ml-1">Thumbnail (Optional)</label>
          <input
            type="file"
            accept="image/*"
            ref={fileInputRef}
            onChange={e => setThumbnail(e.target.files?.[0] || null)}
            className="block w-full text-sm text-neutral-400 file:mr-4 file:py-2.5 file:px-4 file:rounded-xl
              file:border-0 file:text-sm file:font-semibold file:bg-emerald-500/10 file:text-emerald-500
              hover:file:bg-emerald-500/20 transition-colors cursor-pointer"
          />
        </div>

        <div className="flex items-center justify-between mb-8 px-1">
          <p className="text-sm text-neutral-300 font-medium">Private room:</p>
          <Switch isChecked={isPrivate} toggle={togglePrivate} />
        </div>

        <div className="flex justify-center">
          <Button
            type="submit"
            disabled={isPending}
            viewType={isPending ? 'loading' : 'confirm'}
            className="rounded-xl px-10 py-2.5 font-semibold w-full md:w-auto"
          >
            Create Room
          </Button>
        </div>
      </form>
    </MediaButton>
  )
}
