import { InfiniteData, useMutation, useQueryClient } from '@tanstack/react-query'

import { AdminUserRecord, changeUserRoleAction } from '@/lib/actions/admin.actions'
import { extractServerMessage, ServerActionError, unwrapServerAction } from '@/lib/utils/general.utils'
import { useToastStore } from '@/stores/toast-store'
import { PaginatedData } from '@/types/general.types'
import { Role } from '@/types/user.types'

type ChangeRoleVariables = {
  userId: string
  newRole: Role
}

export const useChangeUserRoleMutation = () => {
  const addToast = useToastStore(state => state.addToast)
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async ({ userId, newRole }: ChangeRoleVariables) => {
      return unwrapServerAction(await changeUserRoleAction(userId, newRole))
    },
    onSuccess: (_, { userId, newRole }) => {
      queryClient.setQueriesData(
        { queryKey: ['admin-users-search'] },
        (oldData: InfiniteData<PaginatedData<AdminUserRecord>> | undefined) => {
          if (!oldData) return oldData

          return {
            ...oldData,
            pages: oldData.pages.map(page => ({
              ...page,
              items: page.items.map(user => (user.userId === userId ? { ...user, role: newRole } : user)),
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
