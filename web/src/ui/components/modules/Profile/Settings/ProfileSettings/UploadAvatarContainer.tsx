'use client'

import { useRef, useState, useTransition } from 'react'

import 'react-image-crop/dist/ReactCrop.css'

import { useRouter } from 'next/navigation'
import { useSession } from 'next-auth/react'
import ReactCrop, { centerCrop, makeAspectCrop, type Crop, type PixelCrop } from 'react-image-crop'

import { uploadNewUserPhoto } from '@/lib/actions/user.actions'
import { serverLog } from '@/lib/utils/utils'
import { useProfileStore } from '@/stores/profile.store'
import { useToastStore } from '@/stores/toast-store'
import Button from '@/ui/components/shared/Button'

export default function UploadAvatarContainer() {
  const setLoading = useProfileStore(state => state.setLoading)
  const { data: session, update } = useSession()
  const router = useRouter()
  const addToast = useToastStore(state => state.addToast)
  const [isPending, startTransition] = useTransition()

  const [imgSrc, setImgSrc] = useState('')
  const [crop, setCrop] = useState<Crop>()
  const [completedCrop, setCompletedCrop] = useState<PixelCrop>()
  const imgRef = useRef<HTMLImageElement>(null)
  const fileInputRef = useRef<HTMLInputElement>(null)

  const MAX_FILE_SIZE = 5 * 1024 * 1024

  const onSelectFile = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files.length > 0) {
      const file = e.target.files[0]

      if (file.size > MAX_FILE_SIZE) {
        addToast('File is too large. Maximum size is 5MB', 'error')

        if (fileInputRef.current) {
          fileInputRef.current.value = ''
        }
        return
      }

      const reader = new FileReader()
      reader.addEventListener('load', () => setImgSrc(reader.result?.toString() || ''))
      reader.readAsDataURL(file)
    }
  }

  const onImageLoad = (e: React.SyntheticEvent<HTMLImageElement>) => {
    const { width, height } = e.currentTarget

    const cropSizeInPixels = Math.min(width, height) * 0.9

    const initialCrop = centerCrop(
      makeAspectCrop(
        {
          unit: 'px',
          width: cropSizeInPixels,
        },
        1,
        width,
        height,
      ),
      width,
      height,
    )

    setCrop(initialCrop)
    setCompletedCrop(initialCrop)
  }

  const handleUpload = async () => {
    if (!completedCrop || !imgRef.current || !session?.user?.id) return

    const canvas = document.createElement('canvas')
    const scaleX = imgRef.current.naturalWidth / imgRef.current.width
    const scaleY = imgRef.current.naturalHeight / imgRef.current.height
    canvas.width = Math.floor(completedCrop.width * scaleX)
    canvas.height = Math.floor(completedCrop.height * scaleY)
    const ctx = canvas.getContext('2d')

    if (ctx) {
      ctx.drawImage(
        imgRef.current,
        completedCrop.x * scaleX,
        completedCrop.y * scaleY,
        completedCrop.width * scaleX,
        completedCrop.height * scaleY,
        0,
        0,
        canvas.width,
        canvas.height,
      )

      canvas.toBlob(async blob => {
        if (!blob) return

        startTransition(async () => {
          try {
            setLoading(true)
            const file = new File([blob], 'avatar.png', { type: 'image/png' })
            const formData = new FormData()
            formData.append('file', file)

            const data = await uploadNewUserPhoto(session.user.id, formData)
            if (data?.avatarUrl) {
              await update({ image: data.avatarUrl })
              router.refresh()
              addToast('Avatar updated!', 'success')
              setImgSrc('')
            }
          } catch (error) {
            serverLog('upload error', error)
            addToast('Upload failed', 'error')
          } finally {
            setLoading(false)
          }
        })
      }, 'image/png')
    }
  }

  return (
    <div className="flex flex-col items-center gap-4">
      <input type="file" ref={fileInputRef} onChange={onSelectFile} accept="image/*" className="hidden" />

      {!imgSrc ? (
        <Button
          viewType="confirm"
          className="rounded-xl font-medium"
          onClick={() => {
            if (fileInputRef.current) {
              fileInputRef.current.value = ''
              fileInputRef.current.click()
            }
          }}
        >
          Upload new avatar
        </Button>
      ) : (
        <div className="fixed inset-0 z-100 bg-black/80 flex flex-col items-center justify-center p-4">
          <div className="bg-neutral-900 p-6 rounded-2xl max-w-lg w-full flex flex-col items-center gap-4">
            <h3 className="text-white text-lg font-bold">Adjust your avatar</h3>

            <ReactCrop
              crop={crop}
              onChange={c => setCrop(c)}
              onComplete={c => setCompletedCrop(c)}
              aspect={1}
              circularCrop
            >
              <img ref={imgRef} src={imgSrc} alt="Crop" onLoad={onImageLoad} className="max-h-[60vh] object-contain" />
            </ReactCrop>

            <div className="flex gap-3 w-full">
              <Button onClick={() => setImgSrc('')} className="flex-1 bg-neutral-800 text-white" viewType="cancel">
                Cancel
              </Button>
              <Button
                onClick={handleUpload}
                disabled={isPending}
                className={`flex-1 bg-emerald-500 text-neutral-900
                  ${isPending ? 'cursor-not-allowed' : 'cursor-pointer'}`}
                viewType="confirm"
              >
                {isPending ? 'Saving...' : 'Save Avatar'}
              </Button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
