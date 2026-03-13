'use client'

import { useRef } from 'react'

import { useRouter } from 'next/navigation'
import { useSession } from 'next-auth/react'

import { uploadNewUserPhoto } from '@/app/api/user'
import { serverLog } from '@/lib/utils/utils'
import { useToastStore } from '@/stores/toast-store'
import Button from '@/ui/components/shared/Button'

export default function UploadAvatarContainer() {
  const { data: session, update } = useSession()
  const router = useRouter()
  const addToast = useToastStore(state => state.addToast)
  const fileInputRef = useRef<HTMLInputElement>(null)

  const handleFileChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (!file || !session?.user?.id) return

    if (file.size > 3 * 1024 * 1024) {
      addToast('File is too large (max 3MB)', 'error')
      return
    }

    try {
      const formData = new FormData()
      formData.append('file', file)

      const data = await uploadNewUserPhoto(session.user.id, formData)

      if (data?.avatarUrl) {
        await update({ image: data.avatarUrl })

        router.refresh()
        addToast('Avatar updated successfully', 'success')
      } else {
        addToast('Failed to upload image', 'error')
      }
    } catch (error) {
      serverLog('upload avatar', error)
      addToast('Upload error', 'error')
    } finally {
      if (fileInputRef.current) fileInputRef.current.value = ''
    }
  }

  return (
    <div className="flex items-center">
      <input type="file" ref={fileInputRef} onChange={handleFileChange} accept="image/*" className="hidden" />

      <Button viewType="confirm" className="rounded-xl font-medium" onClick={() => fileInputRef.current?.click()}>
        Upload new avatar
      </Button>
    </div>
  )
}
