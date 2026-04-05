'use client'

import { useEffect, useRef, useState } from 'react'

import Link from 'next/link'

import MoreVertical from '@/assets/icons/shared/more-vertical.svg'
import { useDeletePlaylist } from '@/lib/hooks/api/playlist/useDeletePlaylist'
import { cn } from '@/lib/utils/general.utils'

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

  return (
    <Link href={`/playlists/${props.id}`} className="relative group cursor-pointer flex flex-col">
      <div className="relative">
        <ImageBackground src={props.src} />

        <div
          className="absolute bg-black/80 rounded-lg px-2 py-1 top-1/35 right-1/40 group-hover:top-1/20
            group-hover:right-1/25 transition-all duration-300 text-neutral-200 text-[12px]"
        >
          {props.videoCount} videos
        </div>
      </div>

      <div className="flex justify-between items-start gap-2">
        <div className="flex flex-col overflow-hidden">
          <p className="text-sm font-medium text-neutral-100 line-clamp-1">{props.name}</p>
          <p className="text-[12px] text-neutral-500">{props.isPrivate ? 'Private' : 'Public'} &bull; Playlist</p>
        </div>

        <div className="relative shrink-0" ref={menuRef}>
          <button
            onClick={toggleMenu}
            className={cn(
              'p-1.5 -mr-1.5 -mt-1 rounded-full transition-all duration-300 cursor-pointer z-20',
              'hover:bg-neutral-500/20 active:bg-neutral-500/40',
              isMenuOpen
                ? 'bg-neutral-500/20 text-neutral-300'
                : 'text-neutral-400 opacity-0 group-hover:opacity-100 md:opacity-100',
            )}
          >
            <MoreVertical className="size-5 text-neutral-300" />
          </button>

          {isMenuOpen && (
            <div
              className="absolute right-0 top-full mt-2 w-48 bg-neutral-800 border border-neutral-700/60 shadow-xl
                shadow-black/50 z-50 py-1.5 rounded-xl animate-in fade-in zoom-in-95 duration-200"
              onClick={e => e.preventDefault()}
            >
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
    </Link>
  )
}
