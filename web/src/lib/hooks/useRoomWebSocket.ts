import { useEffect, useRef } from 'react'

import { InfiniteData, useQueryClient } from '@tanstack/react-query'
import io from 'socket.io-client'

import { getWsTokenAction, MemberPoints } from '@/lib/actions/room.actions'
import { useRoomStore } from '@/stores/room.store'
import { BaseServerResponse, PaginatedData } from '@/types/general.types'

import { ChatMessage } from '../actions/chat.actions'
import { RoomVote } from '../actions/vote.actions'

export type VotingStartedPayload = Pick<RoomVote, 'id' | 'voteText' | 'duration' | 'expiresAt' | 'choices'>
export type VotingLockedPayload = { votingId: string }
export type VotingResultsPayload = { votingId: string; rightChoice: string }

export const useRoomWebSocket = (
  roomId: string | undefined,
  chatId: string | undefined,
  userId: string | undefined,
) => {
  const socketRef = useRef<SocketIOClient.Socket | null>(null)
  const queryClient = useQueryClient()

  useEffect(() => {
    if (!chatId || !roomId || !userId) return

    let isMounted = true

    const connectSocket = async () => {
      const tokenRes = await getWsTokenAction()

      if (!tokenRes.success || !tokenRes.data || !isMounted) return

      socketRef.current = io(process.env.NEXT_PUBLIC_WS_URL || 'http://localhost:5000', {
        query: {
          token: tokenRes.data,
          chat_id: chatId,
        },
        transports: ['websocket'],
        upgrade: false,
        reconnection: false,
      })

      const socket = socketRef.current

      socket.on('connect_error', (error: unknown) => {
        console.error('[WS ROOM] Connection Error Detailed:', error)
      })

      socket.on('disconnect', (reason: string) => {
        console.warn('[WS ROOM] Disconnected. Reason:', reason)
      })

      socket.on('REPORT_TIMECODE', (payload: { synchronizeId: string }) => {
        if (payload?.synchronizeId) {
          useRoomStore.getState().setSyncTriggerId(payload.synchronizeId)
        }
      })

      socket.on('SYNCHRONIZE', (payload: { timecode: number }) => {
        if (payload?.timecode !== undefined) {
          useRoomStore.getState().setSyncTargetTimecode(payload.timecode)
        }
      })

      socket.on('QUEUE_UPDATED', () => {
        queryClient.invalidateQueries({ queryKey: ['room-queue', roomId] })
      })

      socket.on('NEW_MESSAGE', (payload: ChatMessage) => {
        if (!payload) return

        type ChatQueryData = InfiniteData<BaseServerResponse<PaginatedData<ChatMessage>>>

        queryClient.setQueryData<ChatQueryData>(['chat-messages', chatId], oldData => {
          if (!oldData || !oldData.pages || oldData.pages.length === 0) {
            return oldData
          }

          const allItems = oldData.pages.flatMap(p => p.data?.items || [])
          if (allItems.some(msg => msg.id === payload.id)) {
            return oldData
          }

          const newPages = [...oldData.pages]
          const firstPage = { ...newPages[0] }

          if (firstPage.data && firstPage.data.items) {
            const tempMsgIndex = firstPage.data.items.findIndex(
              msg =>
                msg.id.startsWith('temp-') && msg.content === payload.content && msg.sender.id === payload.sender.id,
            )

            if (tempMsgIndex !== -1) {
              const newItems = [...firstPage.data.items]
              newItems[tempMsgIndex] = payload
              firstPage.data = {
                ...firstPage.data,
                items: newItems,
              }
            } else {
              firstPage.data = {
                ...firstPage.data,
                items: [payload, ...firstPage.data.items],
              }
            }

            newPages[0] = firstPage
          }

          return { ...oldData, pages: newPages }
        })
      })

      socket.on('ROOM_POINTS_UPDATED', (payload: { added_points: number; totals: Record<string, number> }) => {
        if (!payload || !payload.totals) return

        const myNewPoints = payload.totals[userId]

        if (myNewPoints !== undefined) {
          queryClient.setQueryData(['room-member-points', roomId], (oldData: MemberPoints | undefined) => {
            if (!oldData) return { points: myNewPoints }

            return {
              ...oldData,
              points: myNewPoints,
            }
          })
        }
      })

      socket.on('VOTING_STARTED', (payload: VotingStartedPayload) => {
        if (!payload || !payload.id) return

        queryClient.setQueryData(['room-votes', roomId], (oldData: RoomVote[] | undefined) => {
          const newVote = { ...payload, isLocked: false }

          if (!Array.isArray(oldData)) return [newVote]
          return [newVote, ...oldData]
        })
      })

      socket.on('VOTING_LOCKED', ({ votingId }: VotingLockedPayload) => {
        if (!votingId) return

        queryClient.setQueryData(['room-votes', roomId], (oldData: RoomVote[] | undefined) => {
          if (!Array.isArray(oldData)) return oldData
          return oldData.map((vote: RoomVote) => (vote.id === votingId ? { ...vote, isLocked: true } : vote))
        })
      })

      socket.on('VOTING_RESULTS', ({ votingId, rightChoice }: VotingResultsPayload) => {
        if (!votingId) return

        queryClient.setQueryData(['room-votes', roomId], (oldData: RoomVote[] | undefined) => {
          if (!Array.isArray(oldData)) return oldData

          return oldData.map((vote: RoomVote) =>
            vote.id === votingId ? { ...vote, isLocked: true, rightChoice } : vote,
          )
        })
      })

      socket.on('NEXT_VIDEO_VOTING_STARTED', (payload: { duration: number; expiresAt: string }) => {
        queryClient.setQueryData(['has-next-video-voting', roomId], {
          exists: true,
          duration: payload.duration,
          expiresAt: payload.expiresAt,
        })
      })

      socket.on('NEXT_VIDEO_VOTING_RESULTS', (payload: { winnerId: string }) => {
        if (!payload) return

        queryClient.setQueryData(['has-next-video-voting', roomId], {
          exists: false,
          duration: 0,
          expiresAt: new Date().toISOString(),
        })
      })
    }

    connectSocket()

    return () => {
      isMounted = false
      if (socketRef.current) {
        socketRef.current.disconnect()
      }
    }
  }, [chatId, roomId, userId, queryClient])

  return socketRef.current
}
