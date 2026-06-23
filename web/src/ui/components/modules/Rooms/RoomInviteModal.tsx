'use client'

import { useMemo, useState } from 'react'

import { useDebounce } from 'use-debounce'

import SearchIcon from '@/assets/icons/ic_search.svg'
import { useAddRoomMembersMutation } from '@/lib/hooks/api/room/useAddRoomMembers'
import { useSearchUsersQuery } from '@/lib/hooks/api/user/useSearchUsers'
import { useInfiniteScroll } from '@/lib/hooks/useInfiniteScroll'
import { cn } from '@/lib/utils/general.utils'
import { useToastStore } from '@/stores/toast-store'
import Modal from '@/ui/components/shared/Modal'
import SafeImage from '@/ui/components/shared/SafeImage'
import { BLUR_DATA_URLS } from '@/ui/images'

type Props = {
  roomId: string
  isOpen: boolean
  onClose: () => void
}

export default function RoomInviteModal({ roomId, isOpen, onClose }: Props) {
  const [search, setSearch] = useState('')
  const [debouncedSearch] = useDebounce(search, 500)
  const [selectedUserIds, setSelectedUserIds] = useState<Set<string>>(new Set())

  const addToast = useToastStore(state => state.addToast)
  const { mutate: addMembers, isPending } = useAddRoomMembersMutation()

  const { data, isLoading, isFetchingNextPage, hasNextPage, fetchNextPage } = useSearchUsersQuery(debouncedSearch)

  const users = useMemo(() => {
    return data?.pages.flatMap(page => page?.data?.items || []) || []
  }, [data])

  const lastElementRef = useInfiniteScroll({ isLoading, isFetchingNextPage, hasNextPage, fetchNextPage })

  const toggleUser = (userId: string) => {
    setSelectedUserIds(prev => {
      const newSet = new Set(prev)
      if (newSet.has(userId)) newSet.delete(userId)
      else newSet.add(userId)
      return newSet
    })
  }

  const handleCopyLink = () => {
    navigator.clipboard.writeText(window.location.href)
    addToast('Room link copied to clipboard', 'success')
  }

  const handleInvite = () => {
    if (selectedUserIds.size === 0 || isPending) return

    addMembers({ roomId, userIds: Array.from(selectedUserIds) })
  }

  return (
    <Modal isOpen={isOpen} onClose={onClose} modalClasses="max-w-[400px] w-full bg-[#0A0A0A] p-5 rounded-2xl">
      <div className="flex flex-col gap-4">
        {/* Search Bar & Copy Link */}
        <div className="flex items-center gap-2">
          <div
            className="flex-1 flex items-center bg-neutral-900 border border-neutral-800 rounded-xl px-3 py-2.5
              transition-colors focus-within:border-emerald-500/50"
          >
            <SearchIcon className="w-4 h-4 text-foreground-faint shrink-0" />
            <input
              type="text"
              placeholder="Search a friend"
              value={search}
              onChange={e => setSearch(e.target.value)}
              className="w-full bg-transparent border-none outline-none text-foreground-tertiary
                placeholder:text-foreground-faint ml-2 text-sm"
            />
          </div>
          <button
            onClick={handleCopyLink}
            className="shrink-0 p-2.5 bg-neutral-900 border border-neutral-800 hover:bg-neutral-800 rounded-xl
              transition-colors text-foreground-muted hover:text-foreground-tertiary"
            title="Copy link"
          >
            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1"
              />
            </svg>
          </button>
        </div>

        <div className="grid grid-cols-4 content-start gap-y-5 gap-x-2 py-2 h-80 overflow-y-auto">
          {isLoading && users.length === 0 ? (
            Array.from({ length: 12 }).map((_, i) => <UserSkeleton key={i} />)
          ) : users.length === 0 ? (
            <p className="col-span-4 text-center text-foreground-faint py-10 text-sm">No users found</p>
          ) : (
            users.map((user, index) => {
              const isLast = users.length === index + 1
              const isSelected = selectedUserIds.has(user.userId)

              const userCard = (
                <div
                  key={user.userId}
                  onClick={() => toggleUser(user.userId)}
                  className="flex flex-col items-center gap-2 cursor-pointer group"
                >
                  <div className="relative">
                    <SafeImage
                      src={user.avatarUrl}
                      fallbackType="user"
                      alt={user.username}
                      width={64}
                      height={64}
                      className={cn(
                        'size-16 object-cover rounded-full transition-all duration-300',
                        isSelected
                          ? 'opacity-100 ring-2 ring-emerald-500 ring-offset-2 ring-offset-[#0A0A0A]'
                          : 'opacity-70 group-hover:opacity-100',
                      )}
                      placeholder="blur"
                      blurDataURL={BLUR_DATA_URLS['neutral900']}
                    />
                    {isSelected && (
                      <div
                        className="absolute bottom-0 right-0 size-5 bg-emerald-500 rounded-full flex items-center
                          justify-center border-2 border-[#0A0A0A] animate-in zoom-in duration-200"
                      >
                        <svg
                          className="size-3 text-foreground-strong"
                          fill="none"
                          viewBox="0 0 24 24"
                          stroke="currentColor"
                        >
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={3} d="M5 13l4 4L19 7" />
                        </svg>
                      </div>
                    )}
                  </div>
                  <span
                    className="text-xs text-foreground-subtle truncate w-full text-center transition-colors
                      group-hover:text-foreground-strong"
                  >
                    {user.username}
                  </span>
                </div>
              )

              if (isLast) {
                return (
                  <div key={`last-${user.userId}`} ref={lastElementRef}>
                    {userCard}
                  </div>
                )
              }
              return userCard
            })
          )}
          {isFetchingNextPage && (
            <div className="col-span-4 flex justify-center py-2">
              <div className="size-5 border-4 border-emerald-500/20 border-t-emerald-500 rounded-full animate-spin" />
            </div>
          )}
        </div>

        {/* Action Button */}
        <button
          onClick={handleInvite}
          disabled={selectedUserIds.size === 0 || isPending}
          className={cn(
            'w-full py-3 rounded-xl font-medium transition-all duration-300',
            selectedUserIds.size > 0 && !isPending
              ? 'bg-neutral-800 hover:bg-neutral-700 text-foreground-strong cursor-pointer'
              : 'bg-neutral-900 text-foreground-faint cursor-not-allowed',
          )}
        >
          {isPending ? 'Sending...' : `Send an invite ${selectedUserIds.size > 0 ? `(${selectedUserIds.size})` : ''}`}
        </button>
      </div>
    </Modal>
  )
}

function UserSkeleton() {
  return (
    <div className="flex flex-col items-center gap-2 animate-pulse">
      <div className="size-16 rounded-full bg-neutral-800/80 shrink-0" />
      <div className="h-2.5 w-14 bg-neutral-800/80 rounded mt-1" />
    </div>
  )
}
