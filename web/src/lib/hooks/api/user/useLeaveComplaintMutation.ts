import { useMutation } from '@tanstack/react-query'

import { leaveComplaint } from '@/lib/actions/user.actions'
import { extractServerMessage, ServerActionError, unwrapServerAction } from '@/lib/utils/general.utils'
import { useToastStore } from '@/stores/toast-store'

type LeaveComplaintParams = {
  targetUserId: string
  reasonType: string
  targetType: 'Video' | 'User'
  additionalInfo: string
}

export const useLeaveComplaintMutation = (options?: { onSuccess?: () => void }) => {
  const addToast = useToastStore(state => state.addToast)

  return useMutation({
    mutationFn: async ({ targetUserId, reasonType, targetType, additionalInfo }: LeaveComplaintParams) => {
      return unwrapServerAction(await leaveComplaint(targetUserId, reasonType, targetType, additionalInfo))
    },
    onSuccess: () => {
      addToast('Complaint submitted successfully', 'success')

      if (options?.onSuccess) {
        options.onSuccess()
      }
    },
    onError: (error: ServerActionError) => {
      const errorMessage = extractServerMessage(error.errors) || error.message

      addToast(errorMessage, 'error')
    },
  })
}
