'use client'

import { motion } from 'framer-motion'
import { Route } from 'next'
import { usePathname, useRouter, useSearchParams } from 'next/navigation'

import { cn } from '@/lib/utils/utils'

type Genre = {
  label: string
  param: string
}

const GENRES: Genre[] = [
  { label: 'all', param: 'all' },
  { label: 'cinema', param: 'cinema' },
  { label: 'music', param: 'music' },
  { label: 'games', param: 'games' },
  { label: 'education', param: 'education' },
]

export default function GenreToggle() {
  const searchParams = useSearchParams()
  const pathname = usePathname()
  const { replace } = useRouter()

  const currentGenre = searchParams.get('genre') || 'all'

  const handleGenre = (genre: string) => {
    const params = new URLSearchParams(searchParams)
    params.set('page', '1')
    params.set('genre', genre)

    const url = `${pathname}?${params.toString()}`
    replace(url as Route, { scroll: false })
  }

  return (
    <div className="flex rounded-xl bg-neutral-900/60 p-1 border border-neutral-800 w-fit relative">
      {GENRES.map(genre => {
        const isActive = currentGenre === genre.param

        return (
          <div
            key={genre.param}
            onClick={() => handleGenre(genre.param)}
            className={cn(
              `relative px-4 py-1.5 md:px-5 md:py-2 z-10 block font-semibold md:font-bold uppercase text-[13px]
              tracking-wide cursor-pointer text-nowrap`,
              'transition-colors duration-300 rounded-lg',
              isActive ? 'text-neutral-100' : 'text-neutral-500 hover:text-neutral-300',
            )}
          >
            {isActive && (
              <motion.div
                layoutId="active-genre-pill"
                className="absolute inset-0 bg-neutral-800 rounded-lg -z-10 shadow-sm border border-neutral-700/50"
                transition={{ type: 'spring', stiffness: 400, damping: 30 }}
              />
            )}
            {genre.label}
          </div>
        )
      })}
    </div>
  )
}
