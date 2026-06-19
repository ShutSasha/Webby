import { useMutation, useQueryClient } from '@tanstack/react-query'

import { addRoomMembersAction } from '@/lib/actions/room.actions'
import { extractServerMessage, ServerActionError, unwrapServerAction } from '@/lib/utils/general.utils'
import { useToastStore } from '@/stores/toast-store'

export const useAddRoomMembersMutation = () => {
  const addToast = useToastStore(state => state.addToast)
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async ({ roomId, userIds }: { roomId: string; userIds: string[] }) => {
      return unwrapServerAction(await addRoomMembersAction(roomId, userIds))
    },
    onSuccess: (_, variables) => {
      addToast('Invitations sent successfully', 'success')
      queryClient.invalidateQueries({ queryKey: ['room-members', variables.roomId] })
    },
    onError: (error: ServerActionError) => {
      const errorMessage = extractServerMessage(error.errors) || error.message

      addToast(errorMessage, 'error')
    },
  })
}
