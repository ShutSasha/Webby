import { useMutation, useQueryClient } from '@tanstack/react-query'

import { setNotificationsReadStatusAction } from '@/lib/actions/notification.actions'
import { extractServerMessage } from '@/lib/utils/general.utils'
import { useToastStore } from '@/stores/toast-store'

export const useSetNotificationsReadStatus = () => {
  const queryClient = useQueryClient()
  const addToast = useToastStore(state => state.addToast)

  return useMutation({
    mutationFn: async (notificationIds: string[]) => await setNotificationsReadStatusAction(notificationIds),
    onSuccess: response => {
      if (response.success) {
        queryClient.invalidateQueries({ queryKey: ['user-notifications'] })
        queryClient.invalidateQueries({ queryKey: ['unread-notifications'] })
        queryClient.invalidateQueries({ queryKey: ['unread-notifications-count'] })
      } else {
        const msg = extractServerMessage(response.errors)
        addToast(msg || 'Error marking notifications as read', 'error')
      }
    },
    onError: () => {
      addToast('Critical error while updating status', 'error')
    },
  })
}