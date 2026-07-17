import { useMutation, useQueryClient } from '@tanstack/react-query'
import { useRouter } from 'next/navigation'
import { useSession } from 'next-auth/react'

import { uploadNewUserPhoto } from '@/lib/actions/user.actions'
import { extractServerMessage, ServerActionError, unwrapServerAction } from '@/lib/utils/general.utils'
import { useToastStore } from '@/stores/toast-store'

type UploadAvatarParams = {
  id: string
  formData: FormData
}

export const useUploadAvatarMutation = (options?: { onSuccess?: () => void; onError?: () => void }) => {
  const queryClient = useQueryClient()
  const addToast = useToastStore(state => state.addToast)
  const { update } = useSession()
  const router = useRouter()

  return useMutation({
    mutationFn: async ({ id, formData }: UploadAvatarParams) => {
      return unwrapServerAction(await uploadNewUserPhoto(id, formData))
    },
    onSuccess: async (data, variables) => {
      if (data?.avatarUrl) {
        await update({ image: data.avatarUrl })
        addToast('Avatar updated!', 'success')
        router.refresh()

        queryClient.invalidateQueries({ queryKey: ['user', variables.id] })
      }

      if (options?.onSuccess) {
        options.onSuccess()
      }
    },
    onError: (error: ServerActionError) => {
      const errorMessage = extractServerMessage(error.errors) || error.message
      addToast(errorMessage, 'error')

      if (options?.onError) {
        options.onError()
      }
    },
  })
}
