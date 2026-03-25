'use client'

import { ChangeEvent, useState } from 'react'

import SearchIcon from '@/assets/icons/ic_search.svg'

type Props = {
  loading: boolean
  handleSearchChange: (value: string) => void
}

export default function Search({ loading, handleSearchChange }: Props) {
  const [query, setQuery] = useState<string>('')

  const handleInput = (e: ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value
    setQuery(value)
    handleSearchChange(value)
  }

  return (
    <div className="relative mb-3">
      <SearchIcon className="absolute top-1/2 -translate-y-1/2 left-4 h-5 w-5 text-neutral-600" aria-hidden="true" />
      <input
        id="search"
        type="text"
        autoComplete="off"
        className="focus:border-emerald-500 focus:ring-emerald-500 ring-[0.3px] ring-transparent block
          placeholder:text-neutral-600 focus:outline-none pl-12 py-3 rounded-2xl font-medium text-sm border
          border-border w-full bg-neutral-900/50 text-neutral-200 transition-all"
        placeholder="Search a video in playlist"
        value={query}
        onChange={handleInput}
      />
      {loading && query && (
        <div className="absolute right-4 top-1/2 -translate-y-1/2">
          <div className="size-4 border-2 border-emerald-500/20 border-t-emerald-500 rounded-full animate-spin" />
        </div>
      )}
    </div>
  )
}
