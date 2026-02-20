'use client'

import { motion } from 'framer-motion'
import { Route } from 'next'
import { usePathname, useRouter, useSearchParams } from 'next/navigation'

type Genre = {
  label: string
  param: string
}

const GENRES: Genre[] = [
  { label: 'cinema', param: 'cinema' },
  { label: 'music', param: 'music' },
  { label: 'games', param: 'games' },
  { label: 'education', param: 'education' },
  { label: 'all', param: 'all' },
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
    <div
      className="flex gap-2 py-1.5 px-2 self-end rounded-2xl overflow-hidden bg-neutral-800 order-2 lg:order-1 relative"
    >
      {GENRES.map(genre => {
        const isActive = currentGenre === genre.param

        return (
          <div
            key={genre.param}
            className={`relative cursor-pointer px-3 py-2 rounded-xl font-semibold md:px-5 md:py-2 z-10 block
            md:font-bold uppercase text-[14px] leading-4.5 transition-colors duration-300
            ${isActive ? 'text-neutral-900' : 'text-neutral-600 hover:hover:bg-[#F5F5F5]/5'}`}
            onClick={() => handleGenre(genre.param)}
          >
            {isActive && (
              <motion.div
                layoutId="active-category-pill"
                className="absolute inset-0 bg-emerald-500 rounded-xl -z-10"
                transition={{ type: 'spring', bounce: 0.2, duration: 0.45 }}
              />
            )}
            {genre.label}
          </div>
        )
      })}
    </div>
  )
}
