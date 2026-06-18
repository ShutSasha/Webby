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
        <Button
          viewType="cancel"
          className="rounded-xl font-medium text-sm px-5 py-2"
          onClick={actions.triggerFileInput}
        >
          Change Avatar
        </Button>
      ) : (
        <div
          className="fixed inset-0 z-100 bg-black/80 backdrop-blur-sm flex flex-col items-center justify-center p-4
            animate-in fade-in duration-200"
        >
          <div
            className="bg-[#0A0A0A] p-8 rounded-3xl border border-neutral-800/50 shadow-2xl max-w-lg w-full flex
              flex-col items-center gap-6 animate-in zoom-in-95 duration-300"
          >
            <h3 className="text-neutral-100 text-xl font-bold">Adjust your avatar</h3>

            <div className="bg-neutral-900/50 rounded-2xl overflow-hidden border border-neutral-800/50 p-2">
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
                  width={500}
                  height={500}
                  alt="Crop"
                  onLoad={actions.onImageLoad}
                  className="object-contain max-h-[400px] w-full"
                />
              </ReactCrop>
            </div>

            <div className="flex gap-3 w-full mt-2">
              <Button onClick={actions.cancelUpload} className="flex-1 rounded-xl" viewType="cancel">
                Cancel
              </Button>
              <Button
                onClick={actions.handleUpload}
                disabled={state.isPending}
                className="flex-1 rounded-xl"
                viewType={state.isPending ? 'loading' : 'confirm'}
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
