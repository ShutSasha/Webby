import { MouseEvent, useState } from 'react'

import { useRouter } from 'next/navigation'

import { cancelVideoUploadAction } from '@/lib/actions/video.actions'
import { useCreateVideoMetadata, useUploadVideoFile } from '@/lib/hooks/api/video/useCreateVideo'
import { useIsClient } from '@/lib/hooks/useIsClient'
import { base64ToFile, fileToBase64 } from '@/lib/utils/file.utils'
import { serverLog } from '@/lib/utils/general.utils'
import { useToastStore } from '@/stores/toast-store'
import { useVideoDraftStore } from '@/stores/video-draft.store'

export const MAX_TAGS = 5

export const useCreateVideoLogic = () => {
  const router = useRouter()
  const addToast = useToastStore(state => state.addToast)

  const step = useVideoDraftStore(state => state.step)
  const uploadedVideoData = useVideoDraftStore(state => state.uploadedVideoData)
  const name = useVideoDraftStore(state => state.name)
  const description = useVideoDraftStore(state => state.description)
  const isPrivate = useVideoDraftStore(state => state.isPrivate)
  const previewBase64 = useVideoDraftStore(state => state.previewBase64)
  const videoTags = useVideoDraftStore(state => state.videoTags || [])

  const setField = useVideoDraftStore(state => state.setField)
  const setDraft = useVideoDraftStore(state => state.setDraft)
  const clearDraft = useVideoDraftStore(state => state.clearDraft)

  const { mutateAsync: uploadFile, isPending: isUploading } = useUploadVideoFile()
  const { mutateAsync: createMetadata, isPending: isSaving } = useCreateVideoMetadata()

  const isClient = useIsClient()
  const [tagInput, setTagInput] = useState('')

  const handleVideoSelect = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (!file) return addToast('Please select a valid video file.', 'error')

    const formData = new FormData()
    formData.append('VideoFile', file)

    const response = await uploadFile(formData)

    if (response.success && response.data) {
      setDraft({
        uploadedVideoData: response.data,
        name: file.name.split('.')[0],
        step: 2,
      })
    }
  }

  const handlePreviewSelect = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (!file) return addToast('Please select a valid image file.', 'error')

    const allowedTypes = ['image/jpeg', 'image/png', 'image/webp']
    if (!allowedTypes.includes(file.type)) {
      addToast('Unsupported file format. Please select a JPEG, PNG, or WEBP image.', 'error')
      e.target.value = ''
      return
    }

    if (file.size > 2 * 1024 * 1024) {
      addToast('Thumbnail image is too large. The maximum allowed size is 2MB.', 'error')
      return
    }

    const base64 = await fileToBase64(file)
    setField('previewBase64', base64)
  }

  const handleSubmitMetadata = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!uploadedVideoData) return addToast('Video data is missing. Please try uploading again.', 'error')
    if (!previewBase64) return addToast('Please upload a thumbnail for your video.', 'error')
    if (!name) return addToast('Please enter a title for your video.', 'error')
    if (!description) return addToast('Please enter a description for your video.', 'error')

    const formData = new FormData()
    formData.append('VideoId', uploadedVideoData.videoId)
    formData.append('Name', name)
    formData.append('Description', description)
    formData.append('IsPrivate', isPrivate.toString())

    const previewFile = base64ToFile(previewBase64, 'thumbnail.png')
    formData.append('PreviewFile', previewFile)

    if (videoTags.length > 0) {
      videoTags.forEach(tag => formData.append('VideoTags', tag))
    }

    const response = await createMetadata(formData)

    if (response.success) {
      clearDraft()
      router.push('/studio')
    }
  }

  const cancelUploadVideo = async () => {
    try {
      if (!uploadedVideoData) return clearDraft()
      const response = await cancelVideoUploadAction(uploadedVideoData.videoId)
      if (response.success) return clearDraft()
    } catch (error) {
      serverLog('FAILED_CANCEL_UPLOAD_VIDEO', error, true)
    }
  }

  const handleTagKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter' || e.key === ',') {
      e.preventDefault()
      const newTag = tagInput.trim().toLowerCase()
      if (!newTag) return

      if (videoTags.length >= MAX_TAGS) {
        addToast(`You can only add up to ${MAX_TAGS} tags.`, 'error')
        return
      }
      if (videoTags.includes(newTag)) {
        setTagInput('')
        return
      }

      setField('videoTags', [...videoTags, newTag])
      setTagInput('')
    }
  }

  const removeTag = (e: MouseEvent<HTMLButtonElement>, tagToRemove: string) => {
    e.preventDefault()
    setField(
      'videoTags',
      videoTags.filter(tag => tag !== tagToRemove),
    )
  }

  return {
    state: {
      isClient,
      step,
      isUploading,
      isSaving,
      name,
      description,
      isPrivate,
      previewBase64,
      videoTags,
      tagInput,
    },
    actions: {
      setField,
      setTagInput,
      handleVideoSelect,
      handlePreviewSelect,
      handleSubmitMetadata,
      cancelUploadVideo,
      handleTagKeyDown,
      removeTag,
    },
  }
}
