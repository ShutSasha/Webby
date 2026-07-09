import React from 'react'

import Image from 'next/image'

import { MAX_TAGS } from '@/lib/hooks/useCreateVideoLogic'
import CoreButton from '@/ui/components/shared/CoreButton'
import Switch from '@/ui/components/shared/Switch'

import VideoTagsInput from './VideoTagsInput'

type Props = {
  name: string
  description: string
  isPrivate: boolean
  previewBase64: string | null
  videoTags: string[]
  tagInput: string
  isSaving: boolean
  isEditingMode?: boolean
  onNameChange: (val: string) => void
  onDescriptionChange: (val: string) => void
  onPrivateChange: (val: boolean) => void
  onPreviewSelect: (e: React.ChangeEvent<HTMLInputElement>) => void
  onTagInputChange: (val: string) => void
  onTagKeyDown: (e: React.KeyboardEvent<HTMLInputElement>) => void
  onTagRemove: (e: React.MouseEvent<HTMLButtonElement>, tag: string) => void
  onSubmit: (e: React.FormEvent) => void
  onCancel: () => void
}

export default function VideoMetadataForm({
  name,
  description,
  isPrivate,
  previewBase64,
  videoTags,
  tagInput,
  isSaving,
  isEditingMode = false,
  onNameChange,
  onDescriptionChange,
  onPrivateChange,
  onPreviewSelect,
  onTagInputChange,
  onTagKeyDown,
  onTagRemove,
  onSubmit,
  onCancel,
}: Props) {
  return (
    <form onSubmit={onSubmit} className="flex flex-col gap-6">
      {!isEditingMode && (
        <div className="flex items-center gap-2 pb-4 border-b border-border">
          <div className="size-2.5 rounded-full bg-emerald-500" />
          <span className="text-sm text-emerald-400 font-medium">File uploaded successfully. Fill in details.</span>
        </div>
      )}

      <div className="flex gap-6 flex-col md:flex-row">
        <div className="flex-1 flex flex-col gap-5">
          <label className="flex flex-col gap-2">
            <span className="text-sm font-medium text-foreground-subtle">
              Title <span className="text-red-500">*</span>
            </span>
            <input
              type="text"
              required
              value={name}
              onChange={e => onNameChange(e.target.value)}
              className="bg-background border border-neutral-700 rounded-lg px-4 py-2.5 text-foreground-secondary
                focus:outline-none focus:border-emerald-500 transition-colors"
              placeholder="Catchy title for your video"
            />
          </label>

          <label className="flex flex-col gap-2">
            <span className="text-sm font-medium text-foreground-subtle">
              Description <span className="text-red-500">*</span>
            </span>
            <textarea
              required
              value={description}
              onChange={e => onDescriptionChange(e.target.value)}
              rows={5}
              className="bg-background border border-neutral-700 rounded-lg px-4 py-2.5 text-foreground-secondary
                focus:outline-none focus:border-emerald-500 transition-colors resize-none"
              placeholder="Tell viewers about your video"
            />
          </label>

          <label className="flex items-center gap-3 mt-2 cursor-pointer w-fit">
            <Switch isChecked={isPrivate} toggle={() => onPrivateChange(!isPrivate)} />
            <span className="text-sm text-foreground-subtle">Make video Private</span>
          </label>

          <VideoTagsInput
            tags={videoTags}
            inputValue={tagInput}
            maxTags={MAX_TAGS}
            onInputChange={onTagInputChange}
            onKeyDown={onTagKeyDown}
            onRemoveTag={onTagRemove}
          />
        </div>

        <div className="w-full md:w-[280px] flex flex-col gap-2">
          <span className="text-sm font-medium text-foreground-subtle">
            Thumbnail <span className="text-red-500">*</span>
          </span>
          <label
            className="relative aspect-video w-full rounded-xl border-2 border-dashed border-neutral-700
              bg-background/50 overflow-hidden flex items-center justify-center cursor-pointer
              hover:border-emerald-500/50 transition-colors group"
          >
            {previewBase64 ? (
              <Image src={previewBase64} width={1280} height={720} alt="Preview" className="object-cover" />
            ) : (
              <div className="text-center p-4">
                <p className="text-sm text-foreground-muted group-hover:text-emerald-400 transition-colors">
                  Upload Thumbnail
                </p>
              </div>
            )}
            <input
              type="file"
              accept="image/jpeg, image/png, image/webp"
              required={!previewBase64}
              className="hidden"
              onChange={onPreviewSelect}
            />
          </label>
          <p className="text-xs text-foreground-faint mt-1">Optimal size 1280x720 (16:9). Max 1.5MB.</p>
        </div>
      </div>

      <div className="flex justify-end gap-3 pt-6 border-t border-border mt-2">
        <CoreButton type="button" variant="secondary" onClick={onCancel} disabled={isSaving}>
          Back to upload
        </CoreButton>
        <CoreButton
          type="submit"
          variant="primary"
          className="w-32"
          isLoading={isSaving}
          disabled={!previewBase64 || !name || !description}
        >
          {isEditingMode ? 'Save and publish' : 'Publish'}
        </CoreButton>
      </div>
    </form>
  )
}
