'use client'

import 'react-image-crop/dist/ReactCrop.css'

import Image from 'next/image'
import { ReactCrop } from 'react-image-crop'

import { useAvatarUpload } from '@/lib/hooks/useAvatarUpload'
import Button from '@/ui/components/shared/Button'

export default function UploadAvatarContainer() {
  const { refs, state, actions } = useAvatarUpload()

  return (
    <div className="flex flex-col items-center gap-4">
      <input type="file" ref={refs.fileInputRef} onChange={actions.onSelectFile} accept="image/*" className="hidden" />

      {!state.imgSrc ? (
        <Button viewType="confirm" className="rounded-xl font-medium" onClick={actions.triggerFileInput}>
          Upload new avatar
        </Button>
      ) : (
        <div className="fixed inset-0 z-100 bg-black/80 flex flex-col items-center justify-center p-4">
          <div className="bg-neutral-900 p-6 rounded-2xl max-w-lg w-full flex flex-col items-center gap-4">
            <h3 className="text-white text-lg font-bold">Adjust your avatar</h3>

            <ReactCrop
              crop={state.crop}
              onChange={c => actions.setCrop(c)}
              onComplete={c => actions.setCompletedCrop(c)}
              aspect={1}
              circularCrop
            >
              <Image
                ref={refs.imgRef}
                src={state.imgSrc}
                width={700}
                height={500}
                alt="Crop"
                onLoad={actions.onImageLoad}
                className="object-contain w-full h-full"
              />
            </ReactCrop>

            <div className="flex gap-3 w-full">
              <Button onClick={actions.cancelUpload} className="flex-1 bg-neutral-800 text-white" viewType="cancel">
                Cancel
              </Button>
              <Button
                onClick={actions.handleUpload}
                disabled={state.isPending}
                className={`flex-1 bg-emerald-500 text-neutral-900
                  ${state.isPending ? 'cursor-not-allowed' : 'cursor-pointer'}`}
                viewType="confirm"
              >
                {state.isPending ? 'Saving...' : 'Save Avatar'}
              </Button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
