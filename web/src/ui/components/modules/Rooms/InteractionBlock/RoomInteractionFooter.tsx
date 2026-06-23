import { useParams } from 'next/navigation'

import CubeIcon from '@/assets/icons/Room/cube-points.svg'
import PollIcon from '@/assets/icons/Room/poll-icon.svg'
import ReactionSmileIcon from '@/assets/icons/Room/reaction-smile.svg'
import { useGetRoomMemberPointsQuery } from '@/lib/hooks/api/room/useGetRoomMemberPoints'
import { useGetRoomVotesQuery } from '@/lib/hooks/api/vote/useGetRoomVotes'
import { formatPoints } from '@/lib/utils/video.utils'
import { useRoomStore } from '@/stores/room.store'

export default function RoomInteractionFooter() {
  return (
    <div className="flex flex-row justify-between items-center">
      <div className="flex items-center gap-2">
        <UserPoints />
        <Reactions />
      </div>
      <Poll />
    </div>
  )
}

function UserPoints() {
  const params = useParams()
  const roomId = params?.id as string | undefined

  const { data, isLoading } = useGetRoomMemberPointsQuery(roomId)

  const points = data?.points ?? 0

  return (
    <div
      className="flex flex-row items-center bg-neutral-900 py-1.5 px-3 gap-2.5 rounded-md group hover:bg-emerald-500
        transition-colors duration-300 ease-in-out cursor-pointer"
    >
      <CubeIcon
        className="size-5 text-foreground-subtle group-hover:text-foreground-inverse-subtle transition-colors
          duration-300 ease-in-out"
      />
      <p
        className="group-hover:text-foreground-inverse-subtle text-foreground-subtle transition-colors duration-300
          ease-in-out leading-5 font-medium"
      >
        {isLoading ? <span className="animate-pulse">...</span> : formatPoints(points)}
      </p>
    </div>
  )
}

function Reactions() {
  return (
    <div
      className="flex flex-row items-center bg-neutral-900 py-1.5 px-3 gap-2.5 rounded-md group hover:bg-emerald-500
        transition-colors duration-300 ease-in-out cursor-pointer"
    >
      <ReactionSmileIcon
        className="size-5 text-foreground-subtle group-hover:text-foreground-inverse-subtle transition-colors
          duration-300 ease-in-out"
      />
    </div>
  )
}

function Poll() {
  const params = useParams()
  const roomId = params?.id as string | undefined
  const setVotesModalOpen = useRoomStore(state => state.setVotesModalOpen)

  const { data: votesResponse } = useGetRoomVotesQuery(roomId || '')
  const votes = votesResponse || []
  const hasActiveVote = votes.some(v => !v.isLocked)

  return (
    <div
      onClick={() => setVotesModalOpen(true)}
      className="relative flex flex-row items-center bg-neutral-900 py-1.5 px-3 gap-2.5 rounded-md group
        hover:bg-emerald-500 transition-colors duration-300 ease-in-out cursor-pointer"
    >
      <PollIcon
        className="size-5 text-foreground-subtle group-hover:text-foreground-inverse-subtle transition-colors
          duration-300 ease-in-out"
      />
      <p
        className="group-hover:text-foreground-inverse-subtle text-foreground-subtle transition-colors duration-300
          ease-in-out leading-5"
      >
        Votes
      </p>

      {hasActiveVote && (
        <span className="absolute -top-1 -right-1 flex size-3">
          <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
          <span className="relative inline-flex rounded-full size-3 bg-emerald-500 border-2 border-neutral-950"></span>
        </span>
      )}
    </div>
  )
}
