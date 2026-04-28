import { useMutation, useQueryClient } from '@tanstack/react-query'

import { deleteNotificationAction } from '@/lib/actions/notification.actions'
import { extractServerMessage } from '@/lib/utils/general.utils'
import { useToastStore } from '@/stores/toast-store'

export const useDeleteNotification = () => {
  const queryClient = useQueryClient()
  const addToast = useToastStore(state => state.addToast)

  return useMutation({
    mutationFn: async (notificationId: string) => await deleteNotificationAction(notificationId),
    onSuccess: response => {
      if (response.success) {

        queryClient.invalidateQueries({ queryKey: ['user-notifications'] })
        queryClient.invalidateQueries({ queryKey: ['unread-notifications'] })
        queryClient.invalidateQueries({ queryKey: ['unread-notifications-count'] })
      } else {
        const msg = extractServerMessage(response.errors)
        addToast(msg || 'Error deleting notification', 'error')
      }
    },
    onError: () => {
      addToast('Critical error while deleting notification', 'error')
    },
  })
}