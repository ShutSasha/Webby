import { create } from 'zustand'

export type TabType = 'queue' | 'users' | 'settings' | 'chat'

interface RoomState {
  tab: TabType
  setTab: (tab: TabType) => void

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
}

export const useRoomStore = create<RoomState>(set => ({
  tab: 'queue',
  setTab: tab => set({ tab: tab }),

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
}))
