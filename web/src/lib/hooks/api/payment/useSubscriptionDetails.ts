import { useQuery } from '@tanstack/react-query'

import { getSubscriptionDetailsAction } from '@/lib/actions/payment.actions'

export const useSubscriptionDetailsQuery = (paymentId: string | null) => {
  return useQuery({
    queryKey: ['subscription-details', paymentId],
    queryFn: async () => {
      if (!paymentId) throw new Error('Payment ID is required')

      const response = await getSubscriptionDetailsAction(paymentId)

      if (!response.success || !response.data) {
        throw new Error(response.message || 'Failed to fetch subscription details')
      }

      return response.data
    },
    enabled: !!paymentId,
    retry: 1,
  })
}
