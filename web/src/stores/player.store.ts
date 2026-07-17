import { create } from 'zustand'
import { persist } from 'zustand/middleware'

type State = {
  prevVolume: number
  baseVolume: number
  setBaseVolume: (volume: number) => void
  setPrevVolume: (volume: number) => void
}

export const usePlayerStore = create(
  persist<State>(
    set => ({
      prevVolume: 1,
      baseVolume: 1,
      setBaseVolume: (volume: number) => set({ baseVolume: volume >= 0 && volume <= 1 ? volume : 0 }),
      setPrevVolume: (volume: number) => set({ prevVolume: volume >= 0 && volume <= 1 ? volume : 0 }),
    }),
    { name: 'player-storage' },
  ),
)

type PlayerPlayState = {
  playing: boolean
  setPlaying: (play: boolean) => void
  togglePlay: () => void
  endedSignal: number
  triggerEnded: () => void
}

export const usePlayerPlayStore = create<PlayerPlayState>(set => ({
  playing: false,
  setPlaying: play => set({ playing: play }),
  togglePlay: () => set(state => ({ playing: !state.playing })),
  endedSignal: 0,
  triggerEnded: () => set(state => ({ endedSignal: state.endedSignal + 1 })),
}))
