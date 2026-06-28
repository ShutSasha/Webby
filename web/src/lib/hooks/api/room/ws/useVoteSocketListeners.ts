import { useEffect } from 'react'

import { useQueryClient } from '@tanstack/react-query'

import { RoomVote } from '@/lib/actions/vote.actions'
import { useRoomStore } from '@/stores/room.store'

export type VotingStartedPayload = Pick<RoomVote, 'id' | 'voteText' | 'duration' | 'expiresAt' | 'choices'>
export type VotingLockedPayload = { votingId: string }
export type VotingResultsPayload = { votingId: string; rightChoice: string }

export const useVoteSocketListeners = (roomId: string | undefined) => {
  const socket = useRoomStore(state => state.socket)
  const queryClient = useQueryClient()

  useEffect(() => {
    if (!socket || !roomId) return

    const handleVotingStarted = (payload: VotingStartedPayload) => {
      if (!payload || !payload.id) return

      queryClient.setQueryData(['room-votes', roomId], (oldData: RoomVote[] | undefined) => {
        const newVote = { ...payload, isLocked: false }

        if (!Array.isArray(oldData)) return [newVote]
        return [newVote, ...oldData]
      })
    }

    const handleVotingLocked = ({ votingId }: VotingLockedPayload) => {
      if (!votingId) return

      queryClient.setQueryData(['room-votes', roomId], (oldData: RoomVote[] | undefined) => {
        if (!Array.isArray(oldData)) return oldData
        return oldData.map((vote: RoomVote) => (vote.id === votingId ? { ...vote, isLocked: true } : vote))
      })
    }

    const handleVotingResults = ({ votingId, rightChoice }: VotingResultsPayload) => {
      if (!votingId) return

      queryClient.setQueryData(['room-votes', roomId], (oldData: RoomVote[] | undefined) => {
        if (!Array.isArray(oldData)) return oldData

        return oldData.map((vote: RoomVote) => (vote.id === votingId ? { ...vote, isLocked: true, rightChoice } : vote))
      })
    }

    const handleNextVideoVotingStarted = (payload: { duration: number; expiresAt: string }) => {
      queryClient.setQueryData(['has-next-video-voting', roomId], {
        exists: true,
        duration: payload.duration,
        expiresAt: payload.expiresAt,
      })
    }

    const handleNextVideoVotingResults = (payload: { winnerId: string }) => {
      if (!payload) return

      queryClient.setQueryData(['has-next-video-voting', roomId], {
        exists: false,
        duration: 0,
        expiresAt: new Date().toISOString(),
      })
    }

    socket.on('VOTING_STARTED', handleVotingStarted)
    socket.on('VOTING_LOCKED', handleVotingLocked)
    socket.on('VOTING_RESULTS', handleVotingResults)
    socket.on('NEXT_VIDEO_VOTING_STARTED', handleNextVideoVotingStarted)
    socket.on('NEXT_VIDEO_VOTING_RESULTS', handleNextVideoVotingResults)

    return () => {
      socket.off('VOTING_STARTED', handleVotingStarted)
      socket.off('VOTING_LOCKED', handleVotingLocked)
      socket.off('VOTING_RESULTS', handleVotingResults)
      socket.off('NEXT_VIDEO_VOTING_STARTED', handleNextVideoVotingStarted)
      socket.off('NEXT_VIDEO_VOTING_RESULTS', handleNextVideoVotingResults)
    }
  }, [socket, roomId, queryClient])
}
