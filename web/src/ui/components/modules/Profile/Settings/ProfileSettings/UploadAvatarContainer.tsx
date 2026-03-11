'use client'

import { useRouter } from 'next/navigation'
import { useSession } from 'next-auth/react'

import { serverLog } from '@/lib/utils/utils'
import { useToastStore } from '@/stores/toast-store'
import Button from '@/ui/components/shared/Button'

export default function UploadAvatarContainer() {
  const { update } = useSession()
  const router = useRouter()
  const addToast = useToastStore(state => state.addToast)

  // https://i.ibb.co/JRZXnfZ8/thumb-1920-614743.png
  // https://i.ibb.co/PsPPVfDL/thumb-1920-1311951.jpg
  const handleUpdateAvatar = async () => {
    try {
      await update({
        image: 'https://i.ibb.co/JRZXnfZ8/thumb-1920-614743.png',
      })
      router.refresh()
      addToast('New photo added', 'success')
    } catch (error) {
      serverLog('upload avatar', error)
    }
  }

  return (
    <Button viewType="confirm" className="rounded-xl font-medium" onClick={handleUpdateAvatar}>
      Upload new avatar
    </Button>
  )
}
