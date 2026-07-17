import { useMutation } from '@tanstack/react-query'

import { initializePaymentAction } from '@/lib/actions/payment.actions'

export const useInitializePayment = () => {
  return useMutation({
    mutationFn: async () => {
      const response = await initializePaymentAction()

      if (!response.success || !response.data) {
        throw new Error(response.message || 'Failed to initialize payment')
      }

      return response.data
    },
  })
}
