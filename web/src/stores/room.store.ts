import { create } from 'zustand'

export type TabType = 'playlist' | 'users' | 'settings' | 'chat'

interface RoomState {
  tab: TabType
  setTab: (tab: TabType) => void
}

export const useRoomStore = create<RoomState>(set => ({
  tab: 'chat',
  setTab: tab => set({ tab: tab }),
}))
