'use client'

import { useEditVideoLogic } from '@/lib/hooks/useEditVideoLogic'
import { Video } from '@/types/video.types'

import VideoMetadataForm from '../Create/VideoMetadataForm'

type Props = {
  initialVideo: Video
}

export default function EditVideoContainer({ initialVideo }: Props) {
  const { state, actions } = useEditVideoLogic(initialVideo)

  return (
    <div className="max-w-3xl w-full bg-surface rounded-[20px] border border-border p-8 overflow-hidden">
      <VideoMetadataForm
        isEditingMode={true}
        name={state.name}
        onNameChange={actions.setName}
        description={state.description}
        onDescriptionChange={actions.setDescription}
        isPrivate={state.isPrivate}
        onPrivateChange={actions.setIsPrivate}
        previewBase64={state.previewBase64}
        onPreviewSelect={actions.handlePreviewSelect}
        videoTags={state.videoTags}
        tagInput={state.tagInput}
        onTagInputChange={actions.setTagInput}
        onTagKeyDown={actions.handleTagKeyDown}
        onTagRemove={actions.removeTag}
        onSubmit={actions.handleSubmitMetadata}
        onCancel={actions.cancelEdit}
        isSaving={state.isSaving}
      />
    </div>
  )
}
