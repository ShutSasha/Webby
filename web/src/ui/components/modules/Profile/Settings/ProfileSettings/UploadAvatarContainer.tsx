'use client'

import { useToastStore } from '@/stores/toast-store'
import Button from '@/ui/components/shared/Button'

export default function UploadAvatarContainer() {
  const addToast = useToastStore(state => state.addToast)
  const handleUpdateAvatar = () => {
    addToast('New photo added', 'success')
  }

  return (
    <Button viewType="confirm" className="rounded-xl font-medium" onClick={handleUpdateAvatar}>
      Upload new avatar
    </Button>
  )
}
