'use client'

import { ChangeEvent, FormEvent, useEffect, useRef, useState } from 'react'

import Link from 'next/link'

import EditPenIcon from '@/assets/icons/pencil.svg'
import MoreVertical from '@/assets/icons/shared/more-vertical.svg'
import { useDeletePlaylist } from '@/lib/hooks/api/playlist/useDeletePlaylist'
import { useUpdatePlaylist } from '@/lib/hooks/api/playlist/useUpdatePlaylist'
import { cn } from '@/lib/utils/general.utils'

import Button from '../../shared/Button'
import Input from '../../shared/Input'
import Modal from '../../shared/Modal'
import Switch from '../../shared/Switch'
import ImageBackground from '../Profile/ImageBackground'

type Props = {
  id: string
  src: string
  name: string
  isPrivate: boolean
  videoCount: number
}

export default function UserPlaylistItem(props: Props) {
  const { mutate, isPending } = useDeletePlaylist()
  const [isMenuOpen, setIsMenuOpen] = useState(false)
  const menuRef = useRef<HTMLDivElement>(null)

  const [isEditModalOpen, setIsEditModalOpen] = useState(false)
  const [editName, setEditName] = useState(props.name)
  const [editIsPrivate, setEditIsPrivate] = useState(props.isPrivate)

  const { mutate: updatePlaylist, isPending: isUpdating } = useUpdatePlaylist({
    onSuccess: () => {
      setIsEditModalOpen(false)
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

    if (isPending) return

    mutate(
      { playlistId: props.id },
      {
        onSuccess: () => {
          setIsMenuOpen(false)
        },
      },
    )
  }

  const handleEditPlaylist = (e: React.MouseEvent) => {
    e.preventDefault()
    e.stopPropagation()
    setIsMenuOpen(false)
    setEditName(props.name)
    setEditIsPrivate(props.isPrivate)
    setIsEditModalOpen(true)
  }

  const handleUpdateSubmit = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    e.stopPropagation()
    updatePlaylist({ playlistId: props.id, name: editName, isPrivate: editIsPrivate })
  }

  return (
    <>
      <div className="relative group cursor-pointer flex flex-col">
        <Link
          href={`/playlists/${props.id}`}
          className="absolute inset-0 z-10 rounded-xl"
          aria-label={`View playlist ${props.name}`}
        />

        <div className="relative">
          <ImageBackground src={props.src} />

          <div
            className="absolute bg-surface-strong/80 rounded-lg px-2 py-1 top-1/35 right-1/40 group-hover:top-1/20
              group-hover:right-1/25 transition-all duration-300 text-foreground-tertiary text-[12px]
              pointer-events-none"
          >
            {props.videoCount} videos
          </div>
        </div>

        <div className="flex justify-between items-start gap-2">
          <div className="flex flex-col overflow-hidden">
            <p className="text-sm font-medium text-foreground-secondary line-clamp-1">{props.name}</p>
            <p className="text-[12px] text-foreground-faint mt-1">
              {props.isPrivate ? 'Private' : 'Public'} &bull; Playlist
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
                  onClick={handleEditPlaylist}
                >
                  <span>Edit</span>
                </button>
                <button
                  className="w-full text-left px-4 py-2 text-sm text-red-400 hover:bg-red-500/10 transition-colors flex
                    items-center gap-3 cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed"
                  onClick={handleDelete}
                  disabled={isPending}
                >
                  <span>{isPending ? 'Deleting...' : 'Delete playlist'}</span>
                </button>
              </div>
            )}
          </div>
        </div>
      </div>

      <Modal isOpen={isEditModalOpen} onClose={() => setIsEditModalOpen(false)}>
        <form className="flex flex-col" onSubmit={handleUpdateSubmit}>
          <h3 className="text-foreground-subtle text-center mb-4 font-semibold text-xl">Edit playlist</h3>

          <div className="relative mb-4">
            <label htmlFor={`edit-playlist-${props.id}`} className="sr-only">
              Playlist Name
            </label>
            <Input
              id={`edit-playlist-${props.id}`}
              placeholder="Enter playlist name..."
              className="py-2.5 pl-10 rounded-xl w-full border-neutral-700"
              value={editName}
              onChange={(e: ChangeEvent<HTMLInputElement>) => setEditName(e.target.value)}
            />
            <EditPenIcon
              className="absolute left-3 top-1/2 -translate-y-1/2 size-4 text-foreground-ghost stroke-[1.5px]"
            />
          </div>

          <div className="flex items-center justify-between mb-6">
            <p className="text-sm text-foreground-subtle font-medium">Private:</p>

            <Switch isChecked={editIsPrivate} toggle={() => setEditIsPrivate(prev => !prev)} />
          </div>

          <div className="flex justify-center">
            <Button
              type="submit"
              disabled={isUpdating}
              viewType={isUpdating ? 'loading' : 'confirm'}
              className="rounded-xl px-8 py-2 font-semibold"
            >
              Save changes
            </Button>
          </div>
        </form>
      </Modal>
    </>
  )
}
