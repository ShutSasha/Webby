import { create } from 'zustand'
import { persist } from 'zustand/middleware'

import { UploadVideoResponse } from '@/types/video.types'

interface VideoDraftState {
  step: 1 | 2
  uploadedVideoData: UploadVideoResponse | null
  name: string
  description: string
  isPrivate: boolean
  previewBase64: string | null
  videoTags: string[]
}

interface VideoDraftActions {
  setField: <K extends keyof VideoDraftState>(field: K, value: VideoDraftState[K]) => void
  setDraft: (data: Partial<VideoDraftState>) => void
  clearDraft: () => void
}

const initialState: VideoDraftState = {
  step: 1,
  uploadedVideoData: null,
  name: '',
  description: '',
  isPrivate: false,
  previewBase64: null,
  videoTags: [],
}

export const useVideoDraftStore = create<VideoDraftState & VideoDraftActions>()(
  persist(
    set => ({
      ...initialState,
      setField: (field, value) => set(state => ({ ...state, [field]: value })),
      setDraft: data => set(state => ({ ...state, ...data })),
      clearDraft: () => set(initialState),
    }),
    {
      name: 'webby-video-draft',
    },
  ),
)
