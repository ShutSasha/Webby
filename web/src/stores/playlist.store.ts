import { create } from 'zustand'

import { MediaType } from '@/types/general.types'

interface PlaylistState {
  activeResourceId: string | null
  activeMediaType: MediaType | null

  initResource: (id: string, type: MediaType) => void
  reset: () => void
  setActiveResourceId: (id: string | null) => void
  setActiveMediaType: (type: MediaType | null) => void
}

export const usePlaylistStore = create<PlaylistState>(set => ({
  activeResourceId: null,
  activeMediaType: null,

  initResource: (id, type) => set({ activeResourceId: id, activeMediaType: type }),
  reset: () => set({ activeResourceId: null, activeMediaType: null }),
  setActiveResourceId: id => set({ activeResourceId: id }),
  setActiveMediaType: type => set({ activeMediaType: type }),
}))
