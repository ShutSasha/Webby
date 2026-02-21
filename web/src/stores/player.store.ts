import { create } from 'zustand'
import { persist } from 'zustand/middleware'

type State = {
  baseVolume: number
  setBaseVolume: (volume: number) => void
}

export const usePlayerStore = create(
  persist<State>(
    set => ({
      baseVolume: 1,
      setBaseVolume: (volume: number) => set({ baseVolume: volume >= 0 || volume <= 1 ? volume : 0 }),
    }),
    { name: 'player-storage' },
  ),
)
