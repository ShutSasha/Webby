import { create } from 'zustand'

import { ReactionSentPayload } from '@/lib/hooks/api/room/ws/useEntertainmentRoomListeners'

export type TabType = 'queue' | 'users' | 'settings' | 'chat'

export type ActiveReaction = ReactionSentPayload & { uid: string }

interface RoomState {
  tab: TabType
  setTab: (tab: TabType) => void

  socket: SocketIOClient.Socket | null
  setSocket: (socket: SocketIOClient.Socket | null) => void

  syncTriggerId: string | null
  syncTargetTimecode: number | null
  isSyncCooldown: boolean
  optimisticPendingId: string | null

  setSyncTriggerId: (id: string | null) => void
  setSyncTargetTimecode: (timecode: number | null) => void
  setSyncCooldown: (isOnCooldown: boolean) => void
  setOptimisticPendingId: (id: string | null) => void

  isVotesModalOpen: boolean
  setVotesModalOpen: (isOpen: boolean) => void

  isReactionsModalOpen: boolean
  setReactionsModalOpen: (isOpen: boolean) => void

  activeReactions: ActiveReaction[]
  addReaction: (reaction: ReactionSentPayload) => void
  removeReaction: (uid: string) => void
}

export const useRoomStore = create<RoomState>(set => ({
  tab: 'queue',
  setTab: tab => set({ tab: tab }),

  socket: null,
  setSocket: socket => set({ socket }),

  syncTriggerId: null,
  syncTargetTimecode: null,
  isSyncCooldown: false,
  optimisticPendingId: null,

  setSyncTriggerId: id => set({ syncTriggerId: id }),
  setSyncTargetTimecode: timecode => set({ syncTargetTimecode: timecode }),
  setSyncCooldown: isOnCooldown => set({ isSyncCooldown: isOnCooldown }),
  setOptimisticPendingId: id => set({ optimisticPendingId: id }),

  isVotesModalOpen: false,
  setVotesModalOpen: isOpen => set({ isVotesModalOpen: isOpen }),

  isReactionsModalOpen: false,
  setReactionsModalOpen: isOpen => set({ isReactionsModalOpen: isOpen }),

  activeReactions: [],
  addReaction: reaction =>
    set(state => ({
      activeReactions: [...state.activeReactions, { ...reaction, uid: `${Date.now()}-${Math.random()}` }],
    })),
  removeReaction: uid =>
    set(state => ({
      activeReactions: state.activeReactions.filter(r => r.uid !== uid),
    })),
}))
