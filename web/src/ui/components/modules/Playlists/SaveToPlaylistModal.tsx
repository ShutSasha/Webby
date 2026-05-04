'use client'

import { useState } from 'react'

import { useDebouncedCallback } from 'use-debounce'

import { useSearchUserPlaylistsQuery } from '@/lib/hooks/api/playlist/useSearchUserPlaylists'
import { useInfiniteScroll } from '@/lib/hooks/useInfiniteScroll'
import { VideoSource } from '@/types/video.types'

import PlaylistItem from './SaveToPlaylistBtn/PlaylistItem'
import Search from './SaveToPlaylistBtn/Search'
import Modal from '../../shared/Modal'

type ModalProps = {
  isOpen: boolean
  onClose: () => void
  videoId: string
  userId: string | undefined
  videoSource: VideoSource | undefined
}

export default function SaveToPlaylistModal({ isOpen, onClose, videoId, userId, videoSource }: ModalProps) {
  const [query, setQuery] = useState('')

  const debouncedSearch = useDebouncedCallback((value: string) => {
    setQuery(value)
  }, 400)

  const { data, isLoading, isFetchingNextPage, hasNextPage, fetchNextPage } = useSearchUserPlaylistsQuery(
    userId,
    query,
    videoId,
    isOpen,
    true,
  )

  const playlists = data?.pages.flatMap(page => page?.data?.items || []) || []

  const lastElementRef = useInfiniteScroll({
    isLoading,
    isFetchingNextPage,
    hasNextPage,
    fetchNextPage,
  })

  return (
    <Modal isOpen={isOpen} onClose={onClose}>
      <div className="flex flex-col w-full gap-4">
        <Search handleSearchChange={debouncedSearch} />

        <div className="flex flex-col max-h-[400px] overflow-y-auto">
          {isLoading && playlists.length === 0 ? (
            <div className="flex justify-center py-10">
              <div className="size-6 border-2 border-emerald-500/20 border-t-emerald-500 rounded-full animate-spin" />
            </div>
          ) : !isLoading && playlists.length === 0 ? (
            <p className="text-center text-neutral-500 py-10">
              {query.trim() !== '' ? 'No playlists found' : 'You have no playlists'}
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
                  userId={userId as string}
                  videoSource={videoSource}
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

          {isFetchingNextPage && (
            <div className="flex justify-center py-4">
              <div className="size-5 border-2 border-emerald-500/20 border-t-emerald-500 rounded-full animate-spin" />
            </div>
          )}

          {!hasNextPage && playlists.length > 0 && (
            <p className="text-center text-xs text-neutral-600 py-4 italic">End of list</p>
          )}
        </div>
      </div>
    </Modal>
  )
}
