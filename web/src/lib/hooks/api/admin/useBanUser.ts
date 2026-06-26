import { InfiniteData, useMutation, useQueryClient } from '@tanstack/react-query'

import { AdminUserRecord, banUserAction } from '@/lib/actions/admin.actions'
import { extractServerMessage, ServerActionError, unwrapServerAction } from '@/lib/utils/general.utils'
import { useToastStore } from '@/stores/toast-store'
import { PaginatedData } from '@/types/general.types'

export const useBanUserMutation = () => {
  const addToast = useToastStore(state => state.addToast)
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (userId: string) => {
      return unwrapServerAction(await banUserAction(userId))
    },
    onSuccess: (_, userId) => {
      queryClient.setQueriesData(
        { queryKey: ['admin-users-search'] },
        (oldData: InfiniteData<PaginatedData<AdminUserRecord>> | undefined) => {
          if (!oldData) return oldData

          return {
            ...oldData,
            pages: oldData.pages.map(page => ({
              ...page,
              items: page.items.map(user => (user.userId === userId ? { ...user, isBanned: true } : user)),
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
