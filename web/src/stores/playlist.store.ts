import { create } from 'zustand'

interface PlaylistState {
  activeVideoId: string | null
  activeMediaType: string | null
  initPlaylist: (videoId: string) => void
  reset: () => void
  setActiveVideoId: (id: string | null) => void
  setActiveMediaType: (type: string | null) => void
}

export const usePlaylistStore = create<PlaylistState>(set => ({
  activeVideoId: null,
  activeMediaType: null,

  initPlaylist: (videoId: string) => set({ activeVideoId: videoId, activeMediaType: null }),
  reset: () => set({ activeVideoId: null, activeMediaType: null }),
  setActiveVideoId: id => set({ activeVideoId: id }),
  setActiveMediaType: type => set({ activeMediaType: type }),
}))
