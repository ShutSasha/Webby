import { InfiniteData, useMutation, useQueryClient } from '@tanstack/react-query'

import { Complaint, acceptComplaintAction } from '@/lib/actions/admin.actions'
import { extractServerMessage, ServerActionError, unwrapServerAction } from '@/lib/utils/general.utils'
import { removePaginatedCacheItem } from '@/lib/utils/query.utils'
import { useToastStore } from '@/stores/toast-store'
import { PaginatedData } from '@/types/general.types'

type AcceptComplaintVariables = {
  complaintId: string
}

export const useAcceptComplaintMutation = () => {
  const addToast = useToastStore(state => state.addToast)
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async ({ complaintId }: AcceptComplaintVariables) => {
      return unwrapServerAction(await acceptComplaintAction(complaintId))
    },

    onMutate: async ({ complaintId }) => {
      await queryClient.cancelQueries({ queryKey: ['complaints'] })

      const previousComplaints = queryClient.getQueryData(['complaints'])

      const updatePaginationCache = (oldData: InfiniteData<PaginatedData<Complaint>> | undefined) => {
        return removePaginatedCacheItem(oldData, complaintId)
      }

      queryClient.setQueriesData({ queryKey: ['complaints'] }, updatePaginationCache)

      return { previousComplaints }
    },
    onError: (error: ServerActionError, _variables, context) => {
      if (context?.previousComplaints) {
        queryClient.setQueryData(['complaints'], context.previousComplaints)
      }
      const errorMessage = extractServerMessage(error.errors) || error.message
      addToast(errorMessage, 'error')
    },
  })
}
