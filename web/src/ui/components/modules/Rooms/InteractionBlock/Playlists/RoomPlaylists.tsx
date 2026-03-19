'use client'

import { useState, useMemo, useEffect } from 'react'

import SearchIcon from '@/assets/icons/ic_search.svg'

import PlaylistItem from './PlaylistItem'

interface Video {
  id: string
  title: string
  thumbnail: string
  isActive?: boolean
  isFolder?: boolean
}

const mockVideos: Video[] = [
  {
    id: '1',
    title: 'Fears to Fathom: Ironbark Lookout',
    thumbnail: 'https://i.ibb.co/PGL4ymBS/thumb-1920-415519.jpg',
    isActive: true,
  },
  {
    id: '2',
    title: 'Fears to Fathom episodes',
    thumbnail: 'https://i.ibb.co/PGL4ymBS/thumb-1920-415519.jpg',
    isFolder: true,
  },
  { id: '3', title: 'Fears to Fathom: Ironbark Lookout', thumbnail: 'https://i.ibb.co/PGL4ymBS/thumb-1920-415519.jpg' },
  { id: '4', title: 'Fears to Fathom: Ironbark Lookout', thumbnail: 'https://i.ibb.co/PGL4ymBS/thumb-1920-415519.jpg' },
]

export default function RoomPlaylists() {
  const [search, setSearch] = useState('')
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const timer = setTimeout(() => setLoading(false), 500)
    return () => clearTimeout(timer)
  }, [])

  const filteredPlaylist = useMemo(() => {
    return mockVideos.filter(v => v.title.toLowerCase().includes(search.toLowerCase()))
  }, [search])

  if (loading) return <p className="text-center text-neutral-500 py-4">Loading queue...</p>

  return (
    <div className="flex flex-col h-full min-h-0">
      <div className="relative group mb-2 shrink-0">
        <SearchIcon
          className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-neutral-500
            group-focus-within:text-emerald-500 transition-colors"
        />
        <input
          type="text"
          placeholder="Search video in queue"
          value={search}
          onChange={e => setSearch(e.target.value)}
          className="w-full bg-neutral-900 rounded-md py-2 pl-10 pr-4 text-sm outline-none border border-transparent
            focus:border-emerald-500/70 transition-all placeholder:text-neutral-500"
        />
      </div>

      <div className="flex-1 overflow-y-auto min-h-0 flex flex-col gap-2 pr-1 custom-scrollbar">
        {filteredPlaylist.map(video => (
          <PlaylistItem key={video.id} video={video} />
        ))}
      </div>
    </div>
  )
}
