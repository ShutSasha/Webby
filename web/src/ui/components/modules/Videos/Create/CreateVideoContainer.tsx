'use client'

import { useCreateVideoLogic } from '@/lib/hooks/useCreateVideoLogic'
import PageLoading from '@/ui/components/shared/PageLoading'

import VideoMetadataForm from './VideoMetadataForm'
import VideoUploadStep from './VideoUploadStep'

export default function CreateVideoContainer() {
  const { state, actions } = useCreateVideoLogic()

  if (!state.isClient) {
    return <PageLoading />
  }

  return (
    <div className="max-w-3xl w-full bg-neutral-900 rounded-[20px] border border-neutral-800 p-8 overflow-hidden">
      {state.step === 1 && (
        <VideoUploadStep isUploading={state.isUploading} onVideoSelect={actions.handleVideoSelect} />
      )}

      {state.step === 2 && (
        <VideoMetadataForm
          name={state.name}
          onNameChange={value => actions.setField('name', value)}
          description={state.description}
          onDescriptionChange={value => actions.setField('description', value)}
          isPrivate={state.isPrivate}
          onPrivateChange={value => actions.setField('isPrivate', value)}
          previewBase64={state.previewBase64}
          onPreviewSelect={actions.handlePreviewSelect}
          videoTags={state.videoTags}
          tagInput={state.tagInput}
          onTagInputChange={actions.setTagInput}
          onTagKeyDown={actions.handleTagKeyDown}
          onTagRemove={actions.removeTag}
          onSubmit={actions.handleSubmitMetadata}
          onCancel={actions.cancelUploadVideo}
          isSaving={state.isSaving}
        />
      )}
    </div>
  )
}
