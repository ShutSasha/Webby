'use client'

import { useEffect, useRef, useState } from 'react'

import Image from 'next/image'
import Link from 'next/link'

import MoreVertical from '@/assets/icons/shared/more-vertical.svg'
import { cn } from '@/lib/utils/utils'
import { BLUR_DATA_URLS } from '@/ui/images'

export default function VideoCard() {
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

  return (
    <Link href={`/videos/22f510be-a044-47b2-9710-7660caba5192`} className="group flex flex-col gap-3">
      <div className="overflow-hidden rounded-lg relative">
        <Image
          src="https://i.ibb.co/d4fJc7FX/anime-moon-landscape.jpg"
          width={400}
          height={400}
          alt=""
          className="aspect-video rounded-lg group-hover:scale-115 transition-transform duration-600"
          placeholder="blur"
          blurDataURL={BLUR_DATA_URLS['neutral900']}
        />
        <div className="absolute inset-0 z-1 bg-black/15 transition-all duration-500 group-hover:bg-black/10" />
      </div>

      {/* description block */}
      <div className="flex flex-row justify-between">
        <div className="flex gap-2">
          <Image
            src="https://i.ibb.co/PsPPVfDL/thumb-1920-1311951.jpg"
            alt=""
            width={200}
            height={200}
            className="size-10 rounded-full object-cover"
            placeholder="blur"
            blurDataURL={BLUR_DATA_URLS['neutral900']}
          />
          <div className="flex flex-col">
            <p className="text-sm leading-4.5 text-neutral-300 line-clamp-1 font-bold">Title</p>
            <p className="text-sm leading-4.5 text-neutral-500">creator</p>
            <p className="text-[12px] leading-4 text-neutral-500">views</p>
          </div>
        </div>

        {/* More Actions Menu */}
        <div className="relative" ref={menuRef}>
          <button
            onClick={toggleMenu}
            className={cn(
              'p-1.5 rounded-full transition-all duration-300 cursor-pointer z-20',
              'hover:bg-neutral-500/20 active:bg-neutral-500/40',
              isMenuOpen ? 'bg-neutral-500/20 text-neutral-300' : 'text-neutral-400 hover:text-neutral-300',
            )}
          >
            <MoreVertical className="size-5 text-neutral-300" />
          </button>

          {isMenuOpen && (
            <div
              className="absolute right-0 top-full mt-2 w-48 bg-neutral-900 border border-neutral-800 rounded-xl
                shadow-2xl z-50 py-1.5 animate-in fade-in zoom-in-95 duration-200"
              onClick={e => e.preventDefault()}
            >
              <button
                className="w-full text-left px-4 py-2.5 text-sm text-neutral-300 hover:bg-neutral-800 transition-colors
                  flex items-center gap-3 cursor-pointer"
                onClick={e => {
                  e.preventDefault()
                  e.stopPropagation()

                  setIsMenuOpen(false)
                }}
              >
                <span>Add to playlist</span>
              </button>

              <button
                className="w-full text-left px-4 py-2.5 text-sm text-red-500 hover:bg-red-500/10 transition-colors flex
                  items-center gap-3 cursor-pointer"
                onClick={e => {
                  e.preventDefault()
                  e.stopPropagation()

                  setIsMenuOpen(false)
                }}
              >
                <span>Report</span>
              </button>
            </div>
          )}
        </div>
      </div>
    </Link>
  )
}
