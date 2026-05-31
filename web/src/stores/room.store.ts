import { create } from 'zustand'

export type TabType = 'playlist' | 'users' | 'settings' | 'chat'

interface RoomState {
  tab: TabType
  setTab: (tab: TabType) => void

  syncTriggerId: string | null
  syncTargetTimecode: number | null
  isSyncCooldown: boolean

  setSyncTriggerId: (id: string | null) => void
  setSyncTargetTimecode: (timecode: number | null) => void
  setSyncCooldown: (isOnCooldown: boolean) => void
}

export const useRoomStore = create<RoomState>(set => ({
  tab: 'chat',
  setTab: tab => set({ tab: tab }),

  syncTriggerId: null,
  syncTargetTimecode: null,
  isSyncCooldown: false,

  setSyncTriggerId: id => set({ syncTriggerId: id }),
  setSyncTargetTimecode: timecode => set({ syncTargetTimecode: timecode }),
  setSyncCooldown: isOnCooldown => set({ isSyncCooldown: isOnCooldown }),
}))
