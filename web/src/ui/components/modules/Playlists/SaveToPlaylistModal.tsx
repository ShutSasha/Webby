'use client'

import { useEffect, useState } from 'react'

import { useDebouncedCallback } from 'use-debounce'

import { useBulkTogglePlaylistMediaMutation } from '@/lib/hooks/api/playlist/useBulkTogglePlaylistMedia'
import { useSearchUserPlaylistsQuery } from '@/lib/hooks/api/playlist/useSearchUserPlaylists'
import { useInfiniteScroll } from '@/lib/hooks/useInfiniteScroll'
import { useToastStore } from '@/stores/toast-store'
import { MediaType } from '@/types/general.types'

import PlaylistItem from './SaveToPlaylistBtn/PlaylistItem'
import { PlaylistItemSkeleton } from './SaveToPlaylistBtn/PlaylistItemSkeleton'
import Search from './SaveToPlaylistBtn/Search'
import Modal from '../../shared/Modal'

type ModalProps = {
  isOpen: boolean
  onClose: () => void
  videoId: string
  userId: string | undefined
  mediaType: MediaType | null
}

export default function SaveToPlaylistModal({ isOpen, onClose, videoId, userId, mediaType }: ModalProps) {
  const [query, setQuery] = useState('')

  const [localSelections, setLocalSelections] = useState<Map<string, boolean>>(new Map())

  const addToast = useToastStore(state => state.addToast)

  const { mutate: saveBulkChanges, isPending } = useBulkTogglePlaylistMediaMutation(userId as string)

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

  const playlists = data?.pages.flatMap(page => page?.items || []) || []

  const lastElementRef = useInfiniteScroll({
    isLoading,
    isFetchingNextPage,
    hasNextPage,
    fetchNextPage,
  })

  useEffect(() => {
    if (!isOpen) {
      setLocalSelections(new Map())
      setQuery('')
    }
  }, [isOpen])

  const handleToggle = (playlistId: string, serverState: boolean) => {
    if (isPending) return

    setLocalSelections(prev => {
      const newMap = new Map(prev)

      const currentState = newMap.has(playlistId) ? newMap.get(playlistId)! : serverState
      const nextState = !currentState

      if (nextState === serverState) {
        newMap.delete(playlistId)
      } else {
        newMap.set(playlistId, nextState)
      }

      return newMap
    })
  }

  const handleSave = () => {
    if (localSelections.size === 0) {
      return
    }

    saveBulkChanges(
      {
        playlistIds: Array.from(localSelections.keys()),
        mediaId: videoId,
        mediaType: mediaType || 'Video',
      },
      {
        onSuccess: () => {
          addToast('Playlists updated successfully', 'success')
          setLocalSelections(new Map())
        },
        onError: error => {
          addToast(error.message || 'Failed to update playlists', 'error')
        },
      },
    )
  }

  return (
    <Modal isOpen={isOpen} onClose={onClose}>
      <div className="flex flex-col w-full gap-4">
        <Search handleSearchChange={debouncedSearch} />

        <div className="flex flex-col h-70 overflow-y-auto scrollbar-hide">
          {isLoading && playlists.length === 0 ? (
            <div className="flex flex-col gap-0.5 mt-1">
              {Array.from({ length: 5 }).map((_, i) => (
                <PlaylistItemSkeleton key={i} />
              ))}
            </div>
          ) : !isLoading && playlists.length === 0 ? (
            <p className="text-center text-foreground0 py-10">
              {query.trim() !== '' ? 'No playlists found' : 'You have no playlists'}
            </p>
          ) : (
            playlists.map((playlist, index) => {
              const isLast = playlists.length === index + 1

              const isActuallyAdded = localSelections.has(playlist.playlistId)
                ? localSelections.get(playlist.playlistId)!
                : playlist.isVideoAdded

              const countDiff = isActuallyAdded === playlist.isVideoAdded ? 0 : isActuallyAdded ? 1 : -1
              const displayCount = playlist.countOfVideos + countDiff
              const itemStatus = isActuallyAdded ? 'added' : playlist.isVideoAdded ? 'removed' : 'none'

              const item = (
                <PlaylistItem
                  key={playlist.playlistId}
                  playlistId={playlist.playlistId}
                  image={playlist.playlistCover}
                  name={playlist.name}
                  count={displayCount}
                  status={itemStatus}
                  disabled={isPending}
                  onToggle={() => handleToggle(playlist.playlistId, playlist.isVideoAdded)}
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
        </div>

        <div className="pt-3 mt-1 border-t border-neutral-800/50">
          <button
            onClick={handleSave}
            disabled={localSelections.size === 0 || isPending}
            className={`w-full py-3 rounded-xl font-semibold transition-all duration-300 ${
              localSelections.size > 0 && !isPending
                ? `bg-emerald-500 hover:bg-emerald-400 text-foreground-inverse shadow-[0_4px_12px_rgba(16,185,129,0.25)]
                  cursor-pointer`
                : 'bg-neutral-800/60 text-foreground0 cursor-not-allowed'
              }`}
          >
            {isPending ? 'Saving...' : 'Save'}
          </button>
        </div>
      </div>
    </Modal>
  )
}
