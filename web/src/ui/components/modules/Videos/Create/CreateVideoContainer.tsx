'use client'

import { MouseEvent, useState } from 'react'

import Image from 'next/image'
import { useRouter } from 'next/navigation'

import { cancelVideoUploadAction } from '@/lib/actions/video.actions'
import { useCreateVideoMetadata, useUploadVideoFile } from '@/lib/hooks/api/video/useCreateVideo'
import { useIsClient } from '@/lib/hooks/useIsClient'
import { base64ToFile, fileToBase64 } from '@/lib/utils/file.utils'
import { serverLog } from '@/lib/utils/general.utils'
import { useToastStore } from '@/stores/toast-store'
import { useVideoDraftStore } from '@/stores/video-draft.store'
import CoreButton from '@/ui/components/shared/CoreButton'
import PageLoading from '@/ui/components/shared/PageLoading'
import Switch from '@/ui/components/shared/Switch'

const MAX_TAGS = 5

export default function CreateVideoContainer() {
  const router = useRouter()
  const addToast = useToastStore(state => state.addToast)
  const step = useVideoDraftStore(state => state.step)
  const uploadedVideoData = useVideoDraftStore(state => state.uploadedVideoData)
  const name = useVideoDraftStore(state => state.name)
  const description = useVideoDraftStore(state => state.description)
  const isPrivate = useVideoDraftStore(state => state.isPrivate)
  const previewBase64 = useVideoDraftStore(state => state.previewBase64)
  const setField = useVideoDraftStore(state => state.setField)
  const setDraft = useVideoDraftStore(state => state.setDraft)
  const clearDraft = useVideoDraftStore(state => state.clearDraft)
  const videoTags = useVideoDraftStore(state => state.videoTags || [])

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
      videoTags.forEach(tag => {
        formData.append('VideoTags', tag)
      })
    }

    const response = await createMetadata(formData)

    if (response.success) {
      clearDraft()
      router.push('/studio')
    }
  }

  const cancelUploadVideo = async () => {
    try {
      if (!uploadedVideoData) {
        return clearDraft()
      }

      const response = await cancelVideoUploadAction(uploadedVideoData.videoId)

      if (response.success) {
        return clearDraft()
      }
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

  if (!isClient) {
    return <PageLoading />
  }

  return (
    <div className="max-w-3xl w-full bg-neutral-900 rounded-[20px] border border-neutral-800 p-8 overflow-hidden">
      {step === 1 && (
        <div
          className="flex flex-col items-center justify-center py-20 border-2 border-dashed border-neutral-700
            rounded-xl bg-neutral-800/20"
        >
          {isUploading ? (
            <div className="flex flex-col items-center gap-4">
              <div className="size-10 border-4 border-emerald-500/20 border-t-emerald-500 rounded-full animate-spin" />
              <p className="text-neutral-400 font-medium animate-pulse">Uploading video to server...</p>
            </div>
          ) : (
            <div className="flex flex-col items-center gap-4">
              <div className="p-4 bg-neutral-800 rounded-full text-neutral-400">
                <svg className="w-8 h-8" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"
                  />
                </svg>
              </div>
              <div className="text-center">
                <p className="text-neutral-200 font-semibold mb-1">Select video file to upload</p>
                <p className="text-neutral-500 text-sm mb-6">MP4, WebM or MKV</p>
              </div>

              <label
                className="cursor-pointer bg-emerald-500 hover:bg-emerald-600 text-neutral-950 px-6 py-2.5 rounded-xl
                  font-semibold transition-colors"
              >
                Choose File
                <input type="file" accept="video/*" className="hidden" onChange={handleVideoSelect} />
              </label>
            </div>
          )}
        </div>
      )}

      {step === 2 && (
        <form onSubmit={handleSubmitMetadata} className="flex flex-col gap-6">
          <div className="flex items-center gap-2 pb-4 border-b border-neutral-800">
            <div className="size-2.5 rounded-full bg-emerald-500" />
            <span className="text-sm text-emerald-400 font-medium">File uploaded successfully. Fill in details.</span>
          </div>

          <div className="flex gap-6 flex-col md:flex-row">
            <div className="flex-1 flex flex-col gap-5">
              <label className="flex flex-col gap-2">
                <span className="text-sm font-medium text-neutral-300">
                  Title <span className="text-red-500">*</span>
                </span>
                <input
                  type="text"
                  required
                  value={name}
                  onChange={e => setField('name', e.target.value)}
                  className="bg-neutral-800 border border-neutral-700 rounded-lg px-4 py-2.5 text-neutral-100
                    focus:outline-none focus:border-emerald-500 transition-colors"
                  placeholder="Catchy title for your video"
                />
              </label>

              <label className="flex flex-col gap-2">
                <span className="text-sm font-medium text-neutral-300">
                  Description <span className="text-red-500">*</span>
                </span>
                <textarea
                  required
                  value={description}
                  onChange={e => setField('description', e.target.value)}
                  rows={5}
                  className="bg-neutral-800 border border-neutral-700 rounded-lg px-4 py-2.5 text-neutral-100
                    focus:outline-none focus:border-emerald-500 transition-colors resize-none"
                  placeholder="Tell viewers about your video"
                />
              </label>
              <label className="flex items-center gap-3 mt-2 cursor-pointer w-fit">
                <Switch isChecked={isPrivate} toggle={() => setField('isPrivate', !isPrivate)} />
                <span className="text-sm text-neutral-300">Make video Private</span>
              </label>

              <div className="flex flex-col gap-2 mt-2">
                <div className="flex items-center justify-between">
                  <span className="text-sm font-medium text-neutral-300">Tags</span>
                  <span className="text-xs text-neutral-500">
                    {videoTags.length} / {MAX_TAGS}
                  </span>
                </div>

                <div
                  className="flex flex-col gap-3 p-3 bg-neutral-800 border border-neutral-700 rounded-lg
                    transition-colors"
                >
                  {videoTags.length > 0 && (
                    <div className="flex flex-wrap gap-2">
                      {videoTags.map(tag => (
                        <span
                          key={tag}
                          className="flex items-center gap-1.5 px-2.5 py-1 rounded-md bg-neutral-700 text-xs font-medium
                            text-neutral-200"
                        >
                          #{tag}
                          <button
                            type="button"
                            onClick={e => removeTag(e, tag)}
                            className="text-neutral-400 hover:text-red-400 transition-colors focus:outline-none"
                            aria-label={`Remove tag ${tag}`}
                          >
                            &times;
                          </button>
                        </span>
                      ))}
                    </div>
                  )}

                  <input
                    type="text"
                    value={tagInput}
                    onChange={e => setTagInput(e.target.value)}
                    onKeyDown={handleTagKeyDown}
                    disabled={videoTags.length >= MAX_TAGS}
                    className="bg-transparent border-none text-sm text-neutral-100 focus:outline-none focus:ring-0
                      w-full placeholder:text-neutral-500 disabled:opacity-50 disabled:cursor-not-allowed"
                    placeholder={
                      videoTags.length >= MAX_TAGS ? 'Maximum tags reached' : 'Add a tag and press Enter or comma'
                    }
                  />
                </div>
              </div>
            </div>

            <div className="w-full md:w-[280px] flex flex-col gap-2">
              <span className="text-sm font-medium text-neutral-300">
                Thumbnail <span className="text-red-500">*</span>
              </span>
              <label
                className="relative aspect-video w-full rounded-xl border-2 border-dashed border-neutral-700
                  bg-neutral-800/50 overflow-hidden flex items-center justify-center cursor-pointer
                  hover:border-emerald-500/50 transition-colors group"
              >
                {previewBase64 ? (
                  <Image src={previewBase64} width={1280} height={720} alt="Preview" className="object-cover" />
                ) : (
                  <div className="text-center p-4">
                    <p className="text-sm text-neutral-400 group-hover:text-emerald-400 transition-colors">
                      Upload Thumbnail
                    </p>
                  </div>
                )}
                <input type="file" accept="image/*" required className="hidden" onChange={handlePreviewSelect} />
              </label>
              <p className="text-xs text-neutral-500 mt-1">Optimal size 1280x720 (16:9)</p>
            </div>
          </div>

          <div className="flex justify-end gap-3 pt-6 border-t border-neutral-800 mt-2">
            <CoreButton type="button" variant="secondary" onClick={cancelUploadVideo} disabled={isSaving}>
              Cancel
            </CoreButton>
            <CoreButton
              type="submit"
              variant="primary"
              className="w-32"
              isLoading={isSaving}
              disabled={!previewBase64 || !name || !description}
            >
              Publish
            </CoreButton>
          </div>
        </form>
      )}
    </div>
  )
}
