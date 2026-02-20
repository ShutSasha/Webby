import { create } from 'zustand'
import { persist } from 'zustand/middleware'

type State = {
  isMobileNavOpen: boolean
  toggleMobileNav: () => void
}

export const useCommonStore = create(
  persist<State>(
    set => ({
      isMobileNavOpen: false,
      toggleMobileNav: () => set(state => ({ isMobileNavOpen: !state.isMobileNavOpen })),
    }),
    { name: 'common-storage' },
  ),
)
