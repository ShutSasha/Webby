'use client'

import { useState } from 'react'

import Image from 'next/image'
import { useRouter } from 'next/navigation'

import { useCreateVideoMetadata, useUploadVideoFile } from '@/lib/hooks/api/video/useCreateVideo'
import { UploadVideoResponse } from '@/types/video.types'
import CoreButton from '@/ui/components/shared/CoreButton'
import Switch from '@/ui/components/shared/Switch'

export default function CreateVideoContainer() {
  const router = useRouter()

  const { mutateAsync: uploadFile, isPending: isUploading } = useUploadVideoFile()
  const { mutateAsync: createMetadata, isPending: isSaving } = useCreateVideoMetadata()

  const [step, setStep] = useState<1 | 2>(1)
  const [uploadedVideoData, setUploadedVideoData] = useState<UploadVideoResponse | null>(null)

  const [name, setName] = useState('')
  const [description, setDescription] = useState('')
  const [isPrivate, setIsPrivate] = useState(false)
  const [previewFile, setPreviewFile] = useState<File | null>(null)
  const [previewUrl, setPreviewUrl] = useState<string | null>(null)

  const handleVideoSelect = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (!file) return

    const formData = new FormData()
    formData.append('VideoFile', file)

    const response = await uploadFile(formData)

    if (response.success && response.data) {
      setUploadedVideoData(response.data)
      setName(response.data.name || file.name.split('.')[0])
      setStep(2)
    }
  }

  const handlePreviewSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (file) {
      setPreviewFile(file)
      setPreviewUrl(URL.createObjectURL(file))
    }
  }

  const handleSubmitMetadata = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!uploadedVideoData || !previewFile || !name || !description) return

    const formData = new FormData()

    formData.append('VideoId', uploadedVideoData.videoId)
    formData.append('Name', name)
    formData.append('Description', description)
    formData.append('PreviewFile', previewFile)
    formData.append('IsPrivate', isPrivate.toString())
    // formData.append('PlaylistId', '...')
    // formData.append('VideoTags', '...')

    const response = await createMetadata(formData)

    if (response.success) {
      router.push('/studio')
      setStep(1)
    }
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
                  onChange={e => setName(e.target.value)}
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
                  onChange={e => setDescription(e.target.value)}
                  rows={5}
                  className="bg-neutral-800 border border-neutral-700 rounded-lg px-4 py-2.5 text-neutral-100
                    focus:outline-none focus:border-emerald-500 transition-colors resize-none"
                  placeholder="Tell viewers about your video"
                />
              </label>
              <label className="flex items-center gap-3 mt-2 cursor-pointer w-fit">
                <Switch isChecked={isPrivate} toggle={() => setIsPrivate(prev => !prev)} />
                <span className="text-sm text-neutral-300">Make video Private</span>
              </label>
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
                {previewUrl ? (
                  <Image src={previewUrl} alt="Preview" fill className="object-cover" />
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
            <CoreButton type="button" variant="secondary" onClick={() => setStep(1)} disabled={isSaving}>
              Cancel
            </CoreButton>
            <CoreButton
              type="submit"
              variant="primary"
              className="w-32"
              isLoading={isSaving}
              disabled={!previewFile || !name || !description}
            >
              Publish
            </CoreButton>
          </div>
        </form>
      )}
    </div>
  )
}
