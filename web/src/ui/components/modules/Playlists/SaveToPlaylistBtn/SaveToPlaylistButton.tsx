'use client'

import { useCallback, useState } from 'react'

import { searchUserPlaylists, UserPlaylistDetails } from '@/app/api/playlists'
import PlusIcon from '@/assets/icons/ic_plus_create.svg'
import { useInfiniteSearch } from '@/lib/hooks/useInfiniteSearch'

import PlaylistItem from './PlaylistItem'
import Search from './Search'
import ActionButton from '../../../shared/ActionButton'
import Modal from '../../../shared/Modal'

type Props = {
  videoId: string
  userId: string
}

export default function SaveToPlaylistButton({ videoId, userId }: Props) {
  const [isOpen, setIsOpen] = useState(false)

  const fetchPlaylistsFn = useCallback(
    (query: string, page: number, pageSize: number) => {
      return searchUserPlaylists(userId, query, page, pageSize, videoId)
    },
    [userId, videoId],
  )

  const {
    items: playlists,
    loading,
    fetchingMore,
    hasMore,
    appliedQuery,
    lastElementRef,
    handleSearchChange,
    reset,
  } = useInfiniteSearch<UserPlaylistDetails>({
    fetchFn: fetchPlaylistsFn,
    pageSize: 15,
    enabled: isOpen && !!userId,
  })

  const handleClose = () => {
    setIsOpen(false)
    setTimeout(reset, 300)
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

            {!hasMore && playlists.length > 0 && (
              <p className="text-center text-xs text-neutral-600 py-4 italic">End of list</p>
            )}
          </div>
        </div>
      </Modal>
    </>
  )
}
