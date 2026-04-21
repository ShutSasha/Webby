'use client'

import { useState } from 'react'

import { AnimatePresence, motion } from 'framer-motion'
import { useDebounce } from 'use-debounce'

import SearchIcon from '@/assets/icons/ic_search.svg'
import XIcon from '@/assets/icons/shared/x.svg'
import { useSearchPlaylistsQuery } from '@/lib/hooks/api/playlist/useSearchPlaylists'
import { usePublicRoomsQuery } from '@/lib/hooks/api/room/usePublicRoomsQuery'
import { useSearchStreamsQuery } from '@/lib/hooks/api/stream/useSearchStreams'
import { useSearchVideosQuery } from '@/lib/hooks/api/video/useSearchVideos'
import { useInfiniteScroll } from '@/lib/hooks/useInfiniteScroll'
import EmptyState from '@/ui/components/shared/EmptyState'
import Modal from '@/ui/components/shared/Modal'

import GlobalSearchCard, { GlobalSearchCardSkeleton } from './GlobalSearchCard'
import GlobalSearchTabs, { SearchTab } from './GlobalSearchTabs'

type Props = {
  isOpen: boolean
  onClose: () => void
}

export default function GlobalSearchModal({ isOpen, onClose }: Props) {
  const [query, setQuery] = useState('')
  const [activeTab, setActiveTab] = useState<SearchTab>('Videos')

  const [appliedQuery] = useDebounce(query, 300)

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      modalClasses="max-w-[740px] w-full bg-[#0A0A0A] border-neutral-800 shadow-2xl p-0 overflow-hidden"
    >
      <div className="flex flex-col w-full h-[85vh] max-h-[800px]">
        <div className="p-4 border-b border-neutral-800/50 shrink-0">
          <div
            className="relative flex items-center w-full bg-neutral-900 border border-neutral-800 rounded-xl px-4 py-3
              transition-colors focus-within:border-emerald-500/50"
          >
            <SearchIcon className="w-5 h-5 text-neutral-500 shrink-0" />
            <input
              type="text"
              placeholder="Search users, rooms, videos, or streams..."
              className="w-full bg-transparent border-none outline-none text-neutral-200 placeholder:text-neutral-500
                ml-3 text-sm"
              value={query}
              onChange={e => setQuery(e.target.value)}
              autoFocus
            />
            {query && (
              <button
                onClick={() => setQuery('')}
                className="p-1 hover:bg-neutral-800 rounded-full transition-colors shrink-0 ml-2"
              >
                <XIcon className="w-4 h-4 text-neutral-400 hover:text-neutral-200" />
              </button>
            )}
          </div>
        </div>

        <div className="shrink-0 pt-2">
          <GlobalSearchTabs activeTab={activeTab} onChange={setActiveTab} />
        </div>

        <div
          className="flex-1 min-h-0 overflow-y-auto pb-6 relative [&::-webkit-scrollbar]:w-1.5
            [&::-webkit-scrollbar-track]:bg-transparent [&::-webkit-scrollbar-thumb]:bg-neutral-800
            [&::-webkit-scrollbar-thumb]:rounded-full hover:[&::-webkit-scrollbar-thumb]:bg-neutral-700"
        >
          <AnimatePresence mode="wait">
            <motion.div
              key={activeTab}
              initial={{ opacity: 0, y: 10 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -10 }}
              transition={{ duration: 0.2 }}
              className="flex flex-col gap-1"
            >
              {activeTab === 'Videos' && <VideosTab query={appliedQuery} />}
              {activeTab === 'Rooms' && <RoomsTab query={appliedQuery} />}
              {activeTab === 'Playlists' && <PlaylistsTab query={appliedQuery} />}
              {activeTab === 'Streams' && <StreamsTab query={appliedQuery} />}
              {activeTab === 'Users' && <UsersTab query={appliedQuery} />}
            </motion.div>
          </AnimatePresence>
        </div>
      </div>
    </Modal>
  )
}

function VideosTab({ query }: { query: string }) {
  const [webbyExpanded, setWebbyExpanded] = useState(false)
  const [ytExpanded, setYtExpanded] = useState(false)

  const {
    data: webbyData,
    isLoading: webbyLoading,
    hasNextPage: hasWebbyNext,
    fetchNextPage: fetchWebby,
    isFetchingNextPage: isFetchingWebby,
  } = useSearchVideosQuery(query, 'Webby')

  const {
    data: ytData,
    isLoading: ytLoading,
    hasNextPage: hasYTNext,
    fetchNextPage: fetchYT,
    isFetchingNextPage: isFetchingYT,
  } = useSearchVideosQuery(query, 'YouTube')

  const rawWebby = webbyData?.pages.flatMap(p => p.data?.items || []) || []
  const webbyVideos = Array.from(new Map(rawWebby.map(v => [v.videoId, v])).values())

  const rawYT = ytData?.pages.flatMap(p => p.data?.items || []) || []
  const ytVideos = Array.from(new Map(rawYT.map(v => [v.videoId, v])).values())

  const displayWebby = webbyExpanded ? webbyVideos : webbyVideos.slice(0, 3)
  const displayYT = ytExpanded ? ytVideos : ytVideos.slice(0, 3)

  if (webbyLoading || ytLoading) return <Skeletons count={6} />
  if (webbyVideos.length === 0 && ytVideos.length === 0) return <EmptyState title="No videos found" />

  return (
    <div className="flex flex-col gap-6 px-1">
      {webbyVideos.length > 0 && (
        <div className="flex flex-col">
          <h3 className="text-[11px] font-bold text-neutral-500 uppercase tracking-wider mb-3 px-2">Webby Videos</h3>
          <div className="flex flex-col gap-1">
            {displayWebby.map(video => (
              <GlobalSearchCard
                key={video.videoId}
                id={video.videoId}
                title={video.name}
                subtitle={`Webby • ${video.views.toLocaleString()} views • by ${video.user.username}`}
                thumbnail={video.previewUrl}
                type="Video"
              />
            ))}
          </div>

          {!webbyExpanded && webbyVideos.length > 3 ? (
            <button
              onClick={() => setWebbyExpanded(true)}
              className="text-sm text-emerald-500 hover:text-emerald-400 font-medium py-2.5 mt-2 text-center w-full
                hover:bg-emerald-500/10 rounded-xl transition-colors"
            >
              Show all ({webbyVideos.length})
            </button>
          ) : webbyExpanded && hasWebbyNext ? (
            <button
              onClick={() => fetchWebby()}
              disabled={isFetchingWebby}
              className="text-sm text-emerald-500 hover:text-emerald-400 font-medium py-2.5 mt-2 text-center w-full
                hover:bg-emerald-500/10 rounded-xl transition-colors disabled:opacity-50"
            >
              {isFetchingWebby ? 'Loading...' : 'Show more'}
            </button>
          ) : null}
        </div>
      )}

      {ytVideos.length > 0 && (
        <div className="flex flex-col">
          <h3 className="text-[11px] font-bold text-neutral-500 uppercase tracking-wider mb-3 px-2">YouTube</h3>
          <div className="flex flex-col gap-1">
            {displayYT.map(video => (
              <GlobalSearchCard
                key={video.videoId}
                id={video.videoId}
                title={video.name}
                subtitle={`YouTube • ${video.views.toLocaleString()} views • by ${video.user.username}`}
                thumbnail={video.previewUrl}
                type="Video"
              />
            ))}
          </div>

          {!ytExpanded && ytVideos.length > 3 ? (
            <button
              onClick={() => setYtExpanded(true)}
              className="text-sm text-emerald-500 hover:text-emerald-400 font-medium py-2.5 mt-2 text-center w-full
                hover:bg-emerald-500/10 rounded-xl transition-colors"
            >
              Show all ({ytVideos.length})
            </button>
          ) : ytExpanded && hasYTNext ? (
            <button
              onClick={() => fetchYT()}
              disabled={isFetchingYT}
              className="text-sm text-emerald-500 hover:text-emerald-400 font-medium py-2.5 mt-2 text-center w-full
                hover:bg-emerald-500/10 rounded-xl transition-colors disabled:opacity-50"
            >
              {isFetchingYT ? 'Loading...' : 'Show more'}
            </button>
          ) : null}
        </div>
      )}
    </div>
  )
}

function RoomsTab({ query }: { query: string }) {
  const { data, isLoading, isFetchingNextPage, hasNextPage, fetchNextPage } = usePublicRoomsQuery(query)
  const rawRooms = data?.pages.flatMap(p => p.data?.items || []) || []
  const rooms = Array.from(new Map(rawRooms.map(r => [r.id, r])).values())
  const lastElementRef = useInfiniteScroll({ isLoading, isFetchingNextPage, hasNextPage, fetchNextPage })

  if (isLoading && rooms.length === 0) return <Skeletons count={12} />
  if (rooms.length === 0) return <EmptyState title="No rooms found" />

  return (
    <div className="flex flex-col gap-1 px-1">
      {rooms.map((room, index) => {
        const isLast = rooms.length === index + 1
        const card = (
          <GlobalSearchCard
            key={room.id}
            id={room.id}
            title={room.name}
            subtitle={`${room.categoryName} • Hosted by ${room.hostUsername}`}
            thumbnail={room.thumbnail}
            type="Room"
          />
        )
        return isLast ? (
          <div key={`last-${room.id}`} ref={lastElementRef}>
            {card}
          </div>
        ) : (
          card
        )
      })}
      {isFetchingNextPage && <LoadingSpinner />}
    </div>
  )
}

function PlaylistsTab({ query }: { query: string }) {
  const { data, isLoading, isFetchingNextPage, hasNextPage, fetchNextPage } = useSearchPlaylistsQuery(query)
  const rawPlaylists = data?.pages.flatMap(p => p.data?.items || []) || []
  const playlists = Array.from(new Map(rawPlaylists.map(pl => [pl.playlistId, pl])).values())

  const lastElementRef = useInfiniteScroll({ isLoading, isFetchingNextPage, hasNextPage, fetchNextPage })

  if (isLoading && playlists.length === 0) return <Skeletons count={12} />
  if (playlists.length === 0) return <EmptyState title="No playlists found" />

  return (
    <div className="flex flex-col gap-1 px-1">
      {playlists.map((pl, index) => {
        const isLast = playlists.length === index + 1
        const card = (
          <GlobalSearchCard
            key={pl.playlistId}
            id={pl.playlistId}
            title={pl.name}
            subtitle={`${pl.isPrivate ? 'Private' : 'Public'} • ${pl.countOfVideos} videos`}
            thumbnail={pl.playlistCover}
            type="Playlist"
          />
        )
        return isLast ? (
          <div key={`last-${pl.playlistId}`} ref={lastElementRef}>
            {card}
          </div>
        ) : (
          card
        )
      })}
      {isFetchingNextPage && <LoadingSpinner />}
    </div>
  )
}

function StreamsTab({ query }: { query: string }) {
  const { data, isLoading, isFetchingNextPage, hasNextPage, fetchNextPage } = useSearchStreamsQuery(query)
  const rawStreams = data?.pages.flatMap(p => p.data?.items || []) || []
  const streams = Array.from(new Map(rawStreams.map(s => [s.streamId, s])).values())
  const lastElementRef = useInfiniteScroll({ isLoading, isFetchingNextPage, hasNextPage, fetchNextPage })

  if (isLoading && streams.length === 0) return <Skeletons count={12} />
  if (streams.length === 0) return <EmptyState title="No streams found" />

  return (
    <div className="flex flex-col gap-1 px-1">
      {streams.map((stream, index) => {
        const isLast = streams.length === index + 1
        const card = (
          <GlobalSearchCard
            key={stream.streamId}
            id={stream.streamId}
            title={stream.name}
            subtitle={`Twitch • ${stream.viewers.toLocaleString()} viewers • ${stream.user.username}`}
            thumbnail={stream.previewUrl}
            type="Stream"
          />
        )
        return isLast ? (
          <div key={`last-${stream.streamId}`} ref={lastElementRef}>
            {card}
          </div>
        ) : (
          card
        )
      })}
      {isFetchingNextPage && <LoadingSpinner />}
    </div>
  )
}

function UsersTab({ query }: { query: string }) {
  const mockUsers = [
    {
      id: '1',
      username: 'guuuntersteam',
      role: 'Admin',
      avatar:
        'https://webby-watch-platform-bucket.s3.eu-north-1.amazonaws.com/user_data/019cbcd6-163d-7868-bb60-7af9665e27b6/18:15:59thumbnail.jpg',
    },
    {
      id: '2',
      username: 'qwerty123123',
      role: 'User',
      avatar:
        'https://webby-watch-platform-bucket.s3.eu-north-1.amazonaws.com/user_data/019cde7f-add6-7c2f-8e3a-eff25292b5c3/00:22:01avatar.png',
    },
  ]

  const filtered = mockUsers.filter(u => u.username.toLowerCase().includes(query.toLowerCase()))

  if (filtered.length === 0) return <EmptyState title="No users found" />

  return (
    <div className="flex flex-col gap-1 px-1">
      {filtered.map(user => (
        <GlobalSearchCard
          key={user.id}
          id={user.id}
          title={user.username}
          subtitle={user.role}
          thumbnail={user.avatar}
          type="User"
          showAddButton={false}
        />
      ))}
    </div>
  )
}

function Skeletons({ count, isUser = false }: { count: number; isUser?: boolean }) {
  return (
    <div className="flex flex-col gap-1 px-1">
      {[...Array(count)].map((_, i) => (
        <GlobalSearchCardSkeleton key={i} isUser={isUser} />
      ))}
    </div>
  )
}

function LoadingSpinner() {
  return (
    <div className="w-full flex justify-center py-6">
      <div className="size-6 border-4 border-emerald-500/20 border-t-emerald-500 rounded-full animate-spin" />
    </div>
  )
}
