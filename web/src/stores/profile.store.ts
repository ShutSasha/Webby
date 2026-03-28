import { create } from 'zustand'

interface ProfileState {
  isLoading: boolean
  setLoading: (loading: boolean) => void
}

export const useProfileStore = create<ProfileState>(set => ({
  isLoading: false,
  setLoading: loading => set({ isLoading: loading }),
}))
