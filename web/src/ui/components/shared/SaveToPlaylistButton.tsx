'use client'

import { ChangeEvent, useCallback, useEffect, useRef, useState } from 'react'

import Image from 'next/image'
import { useDebouncedCallback } from 'use-debounce'

import { searchUserPlaylists, UserPlaylistDetails } from '@/app/api/playlists'
import PlusIcon from '@/assets/icons/ic_plus_create.svg'
import SearchIcon from '@/assets/icons/ic_search.svg'
import { BLUR_DATA_URLS } from '@/ui/images'

import ActionButton from './ActionButton'
import Modal from './Modal'
import Input from '../Input'

type Props = {
  userId: string
}

const PAGE_SIZE = 10

export default function SaveToPlaylistButton({ userId }: Props) {
  const [isOpen, setIsOpen] = useState(false)

  const [playlists, setPlaylists] = useState<UserPlaylistDetails[]>([])

  const [loading, setLoading] = useState(false)
  const [query, setQuery] = useState('')
  const [appliedQuery, setAppliedQuery] = useState('')

  const [page, setPage] = useState(1)
  const [fetchingMore, setFetchingMore] = useState(false)
  const [hasMore, setHasMore] = useState(true)

  const observer = useRef<IntersectionObserver | null>(null)

  const handleClose = () => {
    setIsOpen(false)
    setTimeout(() => {
      setQuery('')
      setAppliedQuery('')
      setPage(1)
      setPlaylists([])
    }, 300)
  }

  const lastElementRef = useCallback(
    (node: HTMLDivElement) => {
      if (loading || fetchingMore) return
      if (observer.current) observer.current.disconnect()

      observer.current = new IntersectionObserver(entries => {
        if (entries[0].isIntersecting && hasMore) {
          setPage(prevPage => prevPage + 1)
        }
      })

      if (node) observer.current.observe(node)
    },
    [loading, fetchingMore, hasMore],
  )

  const fetchPlaylists = useCallback(
    async (searchQuery: string, targetPage: number, isInitial: boolean) => {
      if (isInitial) setLoading(true)
      else setFetchingMore(true)

      try {
        const response = await searchUserPlaylists(userId, searchQuery, targetPage, PAGE_SIZE)

        if (response.success && response.data) {
          const newItems = response.data.items

          setPlaylists(prev => (isInitial ? newItems : [...prev, ...newItems]))
          setHasMore(newItems.length === PAGE_SIZE)
        }
      } catch (error) {
        console.error('FETCH_PLAYLISTS_ERROR', error)
      } finally {
        setLoading(false)
        setFetchingMore(false)
      }
    },
    [userId],
  )

  const debouncedSearch = useDebouncedCallback((value: string) => {
    setPage(1)
    setPlaylists([])
    setHasMore(true)
    setAppliedQuery(value)
    fetchPlaylists(value, 1, true)
  }, 400)

  useEffect(() => {
    if (isOpen && userId && playlists.length === 0 && appliedQuery === '') {
      fetchPlaylists('', 1, true)
    }
  }, [isOpen, userId, fetchPlaylists, playlists.length, appliedQuery])

  useEffect(() => {
    if (page > 1 && userId) {
      fetchPlaylists(appliedQuery, page, false)
    }
  }, [page, appliedQuery, fetchPlaylists, userId])

  const handleSearchChange = (e: ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value
    setQuery(value)
    debouncedSearch(value)
  }

  return (
    <>
      <ActionButton onClick={() => setIsOpen(true)} label="Add to playlist" btnClassName="self-end">
        <PlusIcon className="size-4 text-emerald-500" />
      </ActionButton>
      <Modal isOpen={isOpen} onClose={handleClose}>
        <div className="flex flex-col w-full gap-4">
          <div className="relative">
            <SearchIcon className="size-4 absolute top-1/2 -translate-y-1/2 left-3 text-neutral-700" />
            <Input
              className="w-full py-2 rounded-lg pl-9 text-neutral-300 placeholder:text-neutral-700"
              placeholder="Search your playlist or among your playlists idk"
              value={query}
              onChange={handleSearchChange}
            />
          </div>

          <div className="flex flex-col max-h-[400px] overflow-y-auto">
            {loading && playlists.length === 0 ? (
              <div className="flex justify-center py-10">
                <div className="size-6 border-2 border-emerald-500/20 border-t-emerald-500 rounded-full animate-spin" />
              </div>
            ) : !loading && playlists.length === 0 ? (
              <p className="text-center text-neutral-500 py-10">
                {appliedQuery.trim() !== '' ? 'No playlists found' : 'You have no playlists'}
              </p>
            ) : (
              playlists.map((playlist, index) => {
                const isLast = playlists.length === index + 1
                const item = (
                  <VideoItem
                    key={playlist.playlistId}
                    image={playlist.playlistCover}
                    name={playlist.name}
                    count={playlist.countOfVideos}
                  />
                )

                if (isLast) {
                  return (
                    <div ref={lastElementRef} key={`last-${playlist.playlistId}`}>
                      {item}
                    </div>
                  )
                }
                return item
              })
            )}

            {fetchingMore && (
              <div className="flex justify-center py-4">
                <div className="size-5 border-2 border-emerald-500/20 border-t-emerald-500 rounded-full animate-spin" />
              </div>
            )}
          </div>
        </div>
      </Modal>
    </>
  )
}

type VideoItemProps = {
  image: string
  name: string
  count: number
}

function VideoItem({ image, name, count }: VideoItemProps) {
  return (
    <div
      className="flex gap-4 items-center hover:bg-black/40 py-2 px-2.5 mr-1 transition-colors duration-300 ease-in-out
        rounded-xl cursor-pointer"
    >
      <Image
        src={image}
        width={100}
        height={100}
        alt=""
        className="aspect-square size-10 object-cover rounded-lg"
        loading="lazy"
        placeholder="blur"
        blurDataURL={BLUR_DATA_URLS['neutral800']}
      />
      <div className="flex flex-col">
        <p className="text-sm text-neutral-300 line-clamp-1" title={name}>
          {name}
        </p>
        <p className="text-sm text-neutral-500">{count > 1 ? 'videos' : 'video'}</p>
      </div>
    </div>
  )
}
