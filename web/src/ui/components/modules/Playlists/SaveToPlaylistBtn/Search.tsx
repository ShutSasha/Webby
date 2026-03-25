'use client'
import { ChangeEvent, useState } from 'react'

import SearchIcon from '@/assets/icons/ic_search.svg'
import Input from '@/ui/components/Input'

type Props = {
  handleSearchChange: (value: string) => void
}

export default function Search({ handleSearchChange }: Props) {
  const [query, setQuery] = useState('')

  const handleInput = (e: ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value
    setQuery(value)
    handleSearchChange(value)
  }

  return (
    <div className="relative">
      <SearchIcon className="size-4 absolute top-1/2 -translate-y-1/2 left-3 text-neutral-700" />
      <Input
        className="w-full py-2 rounded-lg pl-9 text-neutral-300 placeholder:text-neutral-700"
        placeholder="Search your playlist or among your playlists idk"
        value={query}
        onChange={handleInput}
      />
    </div>
  )
}
