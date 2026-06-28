'use client'

import { ChangeEvent, FormEvent, useEffect, useMemo, useRef, useState } from 'react'

import Image from 'next/image'
import Link from 'next/link'

import EditPenIcon from '@/assets/icons/pencil.svg'
import MoreVertical from '@/assets/icons/shared/more-vertical.svg'
import { useCategoriesQuery } from '@/lib/hooks/api/category/useCategoriesQuery'
import { useDeleteRoom } from '@/lib/hooks/api/room/useDeleteRoom'
import { useUpdateRoom } from '@/lib/hooks/api/room/useUpdateRoom'
import { cn } from '@/lib/utils/general.utils'
import { UserRoom } from '@/types/room.types'
import Button from '@/ui/components/shared/Button'
import Input from '@/ui/components/shared/Input'
import Modal from '@/ui/components/shared/Modal'
import Switch from '@/ui/components/shared/Switch'
import { BLUR_DATA_URLS } from '@/ui/images'

type Props = {
  room: UserRoom
}

export default function UserRoomItem({ room }: Props) {
  const [isMenuOpen, setIsMenuOpen] = useState(false)
  const menuRef = useRef<HTMLDivElement>(null)
  const fileInputRef = useRef<HTMLInputElement>(null)

  const [isEditModalOpen, setIsEditModalOpen] = useState(false)
  const [editName, setEditName] = useState(room.name)
  const [editCategoryName, setEditCategoryName] = useState(room.categoryName)
  const [editIsPrivate, setEditIsPrivate] = useState(room.isPrivate)
  const [editThumbnail, setEditThumbnail] = useState<File | null>(null)

  const { data: categoriesData, isLoading: isLoadingCategories } = useCategoriesQuery('')
  const categories = useMemo(() => {
    return categoriesData?.pages.flatMap(page => page?.data?.items || []) || []
  }, [categoriesData])

  const { mutate: deleteRoom, isPending: isDeleting } = useDeleteRoom({
    onSuccess: () => setIsMenuOpen(false),
  })

  const { mutate: updateRoom, isPending: isUpdating } = useUpdateRoom({
    onSuccess: () => {
      setIsEditModalOpen(false)
      setEditThumbnail(null)
      if (fileInputRef.current) fileInputRef.current.value = ''
    },
  })

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (menuRef.current && !menuRef.current.contains(event.target as Node)) {
        setIsMenuOpen(false)
      }
    }
    if (isMenuOpen) document.addEventListener('mousedown', handleClickOutside)
    return () => document.removeEventListener('mousedown', handleClickOutside)
  }, [isMenuOpen])

  const toggleMenu = (e: React.MouseEvent) => {
    e.preventDefault()
    e.stopPropagation()
    setIsMenuOpen(!isMenuOpen)
  }

  const handleDelete = (e: React.MouseEvent) => {
    e.preventDefault()
    e.stopPropagation()
    if (!isDeleting) deleteRoom(room.id)
  }

  const handleOpenEdit = (e: React.MouseEvent) => {
    e.preventDefault()
    e.stopPropagation()
    setIsMenuOpen(false)
    setEditName(room.name)
    setEditCategoryName(room.categoryName)
    setEditIsPrivate(room.isPrivate)
    setIsEditModalOpen(true)
  }

  const handleUpdateSubmit = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    e.stopPropagation()

    if (!editName || editName.trim().length < 2 || !editCategoryName) return

    const formData = new FormData()
    formData.append('name', editName)
    formData.append('categoryName', editCategoryName)
    formData.append('isPrivate', editIsPrivate ? 'true' : 'false')

    if (editThumbnail) {
      formData.append('thumbnail', editThumbnail)
    }

    updateRoom({ roomId: room.id, formData })
  }

  return (
    <>
      <div className="flex flex-col gap-3 group cursor-pointer relative">
        <Link
          href={`/rooms/${room.id}`}
          className="absolute inset-0 z-10 rounded-xl"
          aria-label={`Enter ${room.name}`}
        />

        <div className="relative aspect-video w-full overflow-hidden rounded-xl bg-background">
          <Image
            src={room.thumbnail}
            alt="Room preview"
            width={640}
            height={480}
            loading="lazy"
            className="object-cover w-full h-full z-0 transition-transform duration-500 ease-out group-hover:scale-105"
            placeholder="blur"
            blurDataURL={BLUR_DATA_URLS['neutral800']}
          />

          <div className="absolute top-2 left-2 z-10 bg-black/60 backdrop-blur-md px-2 py-1 rounded-md
            pointer-events-none">
            <p className="uppercase font-bold text-[10px] tracking-wider text-neutral-100 line-clamp-1 max-w-[100px]">
              {room.categoryName}
            </p>
          </div>
        </div>

        <div className="flex justify-between items-start gap-2 px-1">
          <div className="flex flex-col overflow-hidden">
            <h3
              className="text-foreground-secondary text-sm font-semibold leading-snug line-clamp-2 transition-colors
                duration-200"
            >
              {room.name}
            </h3>
            <p className="text-[12px] text-foreground-faint mt-1">
              {room.isPrivate ? 'Private' : 'Public'} &bull; Room
            </p>
          </div>

          <div className="relative shrink-0 z-20" ref={menuRef}>
            <button
              onClick={toggleMenu}
              className={cn(
                'p-1.5 -mr-1.5 -mt-1 rounded-full transition-all duration-300 cursor-pointer z-20',
                'hover:bg-surface-faint/20 active:bg-surface-faint/40',
                isMenuOpen
                  ? 'bg-surface-faint/20 text-foreground-subtle'
                  : 'text-foreground-muted opacity-0 group-hover:opacity-100 md:opacity-100',
              )}
            >
              <MoreVertical className="size-5 text-foreground-subtle" />
            </button>

            {isMenuOpen && (
              <div
                className="absolute right-0 top-full mt-2 w-48 bg-background border border-neutral-700/60 shadow-xl
                  shadow-black/50 z-50 py-1.5 rounded-xl animate-in fade-in zoom-in-95 duration-200"
                onClick={e => e.preventDefault()}
              >
                <button
                  className="w-full text-left px-4 py-2 text-sm text-foreground-tertiary hover:bg-surface-tertiary/50
                    transition-colors flex items-center gap-3 cursor-pointer"
                  onClick={handleOpenEdit}
                >
                  <span>Edit</span>
                </button>
                <button
                  className="w-full text-left px-4 py-2 text-sm text-red-400 hover:bg-red-500/10 transition-colors flex
                    items-center gap-3 cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed"
                  onClick={handleDelete}
                  disabled={isDeleting}
                >
                  <span>{isDeleting ? 'Deleting...' : 'Delete room'}</span>
                </button>
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Edit Modal */}
      <Modal isOpen={isEditModalOpen} onClose={() => setIsEditModalOpen(false)}>
        <form className="flex flex-col" onSubmit={handleUpdateSubmit}>
          <h3 className="text-foreground-subtle text-center mb-6 font-semibold text-xl">Edit room</h3>

          <div className="relative mb-4">
            <Input
              placeholder="Enter room name..."
              className="py-2.5 pl-10 rounded-xl w-full border-neutral-700"
              value={editName}
              onChange={(e: ChangeEvent<HTMLInputElement>) => setEditName(e.target.value)}
              maxLength={50}
            />
            <EditPenIcon
              className="absolute left-3 top-1/2 -translate-y-1/2 size-4 text-foreground-ghost stroke-[1.5px]"
            />
          </div>

          <div className="mb-4">
            <select
              value={editCategoryName}
              onChange={e => setEditCategoryName(e.target.value)}
              disabled={isLoadingCategories}
              className="py-2.5 px-4 rounded-xl w-full border border-neutral-700 bg-surface text-foreground-tertiary
                outline-none focus:border-emerald-500 transition-colors appearance-none cursor-pointer
                disabled:opacity-50"
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
            <label className="text-sm text-foreground-subtle font-medium ml-1">New Thumbnail (Optional)</label>
            <input
              type="file"
              ref={fileInputRef}
              accept="image/*"
              onChange={e => setEditThumbnail(e.target.files?.[0] || null)}
              className="block w-full text-sm text-foreground-muted file:mr-4 file:py-2.5 file:px-4 file:rounded-xl
                file:border-0 file:text-sm file:font-semibold file:bg-emerald-500/10 file:text-emerald-500
                hover:file:bg-emerald-500/20 transition-colors cursor-pointer"
            />
          </div>

          <div className="flex items-center justify-between mb-8 px-1">
            <p className="text-sm text-foreground-subtle font-medium">Private room:</p>
            <Switch isChecked={editIsPrivate} toggle={() => setEditIsPrivate(prev => !prev)} />
          </div>

          <div className="flex justify-center">
            <Button
              type="submit"
              disabled={isUpdating}
              viewType={isUpdating ? 'loading' : 'confirm'}
              className="rounded-xl px-10 py-2.5 font-semibold w-full md:w-auto"
            >
              Save changes
            </Button>
          </div>
        </form>
      </Modal>
    </>
  )
}
