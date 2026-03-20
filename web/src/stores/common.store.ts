import { create } from 'zustand'
import { persist } from 'zustand/middleware'

type ForgotPasswordStep = 1 | 2 | 3

type State = {
  isMobileNavOpen: boolean
  toggleMobileNav: () => void

  forgotPasswordStep: ForgotPasswordStep
  forgotPasswordEmail: string
  setForgotPasswordStep: (step: ForgotPasswordStep) => void
  setForgotPasswordEmail: (email: string) => void
  resetForgotPassword: () => void
}

export const useCommonStore = create(
  persist<State>(
    set => ({
      isMobileNavOpen: false,
      toggleMobileNav: () => set(state => ({ isMobileNavOpen: !state.isMobileNavOpen })),
      forgotPasswordStep: 1,
      forgotPasswordEmail: '',
      setForgotPasswordStep: step => set({ forgotPasswordStep: step }),
      setForgotPasswordEmail: email => set({ forgotPasswordEmail: email }),
      resetForgotPassword: () => set({ forgotPasswordStep: 1, forgotPasswordEmail: '' }),
    }),
    { name: 'common-storage' },
  ),
)
