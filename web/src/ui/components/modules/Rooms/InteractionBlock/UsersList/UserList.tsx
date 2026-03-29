'use client'

import { useState, useMemo } from 'react'

import SearchIcon from '@/assets/icons/ic_search.svg'

import UserItem from './UserItem'

interface User {
  id: number
  name: string
  avatar: string
}

export default function UserList({ initialUsers }: { initialUsers: User[] }) {
  const [search, setSearch] = useState('')

  const filteredUsers = useMemo(() => {
    return initialUsers.filter(u => u.name.toLowerCase().includes(search.toLowerCase()))
  }, [search, initialUsers])

  return (
    <>
      {/* Search Bar */}
      <div className="relative group mb-2">
        <SearchIcon
          className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-neutral-500
            group-focus-within:text-emerald-500 transition-colors stroke-[1.5px]"
        />
        <input
          type="text"
          placeholder="Search a user"
          value={search}
          onChange={e => setSearch(e.target.value)}
          className="w-full bg-neutral-900 rounded-md py-2 pl-10 pr-4 text-sm outline-none border border-transparent
            focus:border-emerald-500/70 transition-all placeholder:text-neutral-500"
        />
      </div>

      <div
        className="flex-1 overflow-y-auto min-h-0 flex flex-col gap-1 pr-1 [&::-webkit-scrollbar]:w-1.5
          [&::-webkit-scrollbar-track]:bg-transparent [&::-webkit-scrollbar-thumb]:bg-neutral-900
          [&::-webkit-scrollbar-thumb]:border-0 [&::-webkit-scrollbar-thumb]:rounded-full
          hover:[&::-webkit-scrollbar-thumb]:bg-neutral-800"
      >
        {filteredUsers.map(user => (
          <UserItem key={user.id} user={user} />
        ))}
      </div>
    </>
  )
}
