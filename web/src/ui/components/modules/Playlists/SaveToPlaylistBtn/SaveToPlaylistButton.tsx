'use client'

import { useCallback, useEffect, useRef, useState } from 'react'

import { useDebouncedCallback } from 'use-debounce'

import { searchUserPlaylists, UserPlaylistDetails } from '@/app/api/playlists'
import PlusIcon from '@/assets/icons/ic_plus_create.svg'

import PlaylistItem from './PlaylistItem'
import Search from './Search'
import ActionButton from '../../../shared/ActionButton'
import Modal from '../../../shared/Modal'

type Props = {
  videoId: string
  userId: string
}

const PAGE_SIZE = 10

export default function SaveToPlaylistButton({ videoId, userId }: Props) {
  const [isOpen, setIsOpen] = useState(false)

  const [playlists, setPlaylists] = useState<UserPlaylistDetails[]>([])

  const [loading, setLoading] = useState(false)
  const [appliedQuery, setAppliedQuery] = useState('')

  const [page, setPage] = useState(1)
  const [fetchingMore, setFetchingMore] = useState(false)
  const [hasMore, setHasMore] = useState(true)

  const observer = useRef<IntersectionObserver | null>(null)

  const handleClose = () => {
    setIsOpen(false)
    setTimeout(() => {
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
        const response = await searchUserPlaylists(userId, searchQuery, targetPage, PAGE_SIZE, videoId)

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
    [userId, videoId],
  )

  const debouncedSearch = useDebouncedCallback((value: string) => {
    setPage(1)
    setPlaylists([])
    setHasMore(true)
    setAppliedQuery(value)
    fetchPlaylists(value, 1, true)
  }, 400)

  useEffect(() => {
    if (isOpen && playlists.length === 0 && appliedQuery === '') {
      fetchPlaylists('', 1, true)
    }
  }, [isOpen, fetchPlaylists, playlists.length, appliedQuery])

  useEffect(() => {
    if (page > 1) {
      fetchPlaylists(appliedQuery, page, false)
    }
  }, [page, appliedQuery, fetchPlaylists])

  const handleSearchChange = (value: string) => {
    debouncedSearch(value)
  }

  return (
    <>
      <ActionButton onClick={() => setIsOpen(true)} label="Add to playlist" btnClassName="self-end">
        <PlusIcon className="size-4 text-emerald-500" />
      </ActionButton>
      <Modal isOpen={isOpen} onClose={handleClose}>
        <div className="flex flex-col w-full gap-4">
          <Search handleSearchChange={handleSearchChange} />

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
                  <PlaylistItem
                    key={playlist.playlistId}
                    playlistId={playlist.playlistId}
                    image={playlist.playlistCover}
                    name={playlist.name}
                    count={playlist.countOfVideos}
                    isVideoAdded={playlist.isVideoAdded}
                    videoId={videoId}
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
