import { MouseEvent, useState } from 'react'

import { useRouter } from 'next/navigation'

import { useUpdateVideoMetadata } from '@/lib/hooks/api/video/useEditVideo'
import { fileToBase64 } from '@/lib/utils/file.utils'
import { useToastStore } from '@/stores/toast-store'
import { Video } from '@/types/video.types'

export const MAX_TAGS = 10

export const useEditVideoLogic = (initialVideo: Video) => {
  const router = useRouter()
  const addToast = useToastStore(state => state.addToast)

  const [name, setName] = useState(initialVideo.name || '')
  const [description, setDescription] = useState(initialVideo.description || '')
  const [isPrivate, setIsPrivate] = useState(initialVideo.isPrivate || false)

  const [previewBase64, setPreviewBase64] = useState<string | null>(initialVideo.previewUrl || null)
  const [previewFile, setPreviewFile] = useState<File | null>(null)

  const [videoTags, setVideoTags] = useState<string[]>(initialVideo.videoTags || initialVideo.videotags || [])
  const [tagInput, setTagInput] = useState('')

  const { mutateAsync: updateMetadata, isPending: isSaving } = useUpdateVideoMetadata()

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
    setPreviewBase64(base64)
    setPreviewFile(file)
  }

  const handleSubmitMetadata = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!previewBase64) return addToast('Please upload a thumbnail for your video.', 'error')
    if (!name) return addToast('Please enter a title for your video.', 'error')
    if (!description) return addToast('Please enter a description for your video.', 'error')

    const formData = new FormData()
    formData.append('VideoId', initialVideo.videoId)
    formData.append('Name', name)
    formData.append('Description', description)
    formData.append('IsPrivate', isPrivate.toString())

    if (previewFile) {
      formData.append('PreviewFile', previewFile)
    }

    if (videoTags.length > 0) {
      videoTags.forEach(tag => formData.append('VideoTags', tag))
    }

    const response = await updateMetadata(formData)

    if (response.success) {
      router.push('/studio')
    }
  }

  const cancelEdit = () => {
    router.back()
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

      setVideoTags([...videoTags, newTag])
      setTagInput('')
    }
  }

  const removeTag = (e: MouseEvent<HTMLButtonElement>, tagToRemove: string) => {
    e.preventDefault()
    setVideoTags(videoTags.filter(tag => tag !== tagToRemove))
  }

  return {
    state: {
      isSaving,
      name,
      description,
      isPrivate,
      previewBase64,
      videoTags,
      tagInput,
    },
    actions: {
      setName,
      setDescription,
      setIsPrivate,
      setTagInput,
      handlePreviewSelect,
      handleSubmitMetadata,
      cancelEdit,
      handleTagKeyDown,
      removeTag,
    },
  }
}
