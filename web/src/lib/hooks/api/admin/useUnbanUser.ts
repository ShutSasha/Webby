import { InfiniteData, useMutation, useQueryClient } from '@tanstack/react-query'

import { AdminUserRecord, unbanUserAction } from '@/lib/actions/admin.actions'
import { extractServerMessage, ServerActionError, unwrapServerAction } from '@/lib/utils/general.utils'
import { useToastStore } from '@/stores/toast-store'
import { PaginatedData } from '@/types/general.types'

export const useUnbanUserMutation = () => {
  const addToast = useToastStore(state => state.addToast)
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (userId: string) => {
      return unwrapServerAction(await unbanUserAction(userId))
    },
    onSuccess: (_, userId) => {
      addToast('User successfully unbanned', 'success')

      queryClient.setQueriesData(
        { queryKey: ['admin-users-search'] },
        (oldData: InfiniteData<PaginatedData<AdminUserRecord>> | undefined) => {
          if (!oldData) return oldData

          return {
            ...oldData,
            pages: oldData.pages.map(page => ({
              ...page,
              items: page.items.map(user => (user.userId === userId ? { ...user, isBanned: false } : user)),
            })),
          }
        },
      )
    },
    onError: (error: ServerActionError) => {
      const errorMessage = extractServerMessage(error.errors) || error.message
      addToast(errorMessage, 'error')
    },
  })
}
