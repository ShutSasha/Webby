import { useMutation, useQueryClient } from '@tanstack/react-query'

import { removeRoomMemberAction } from '@/lib/actions/room.actions'

export const useRemoveRoomMemberMutation = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async ({ roomId, memberId }: { roomId: string; memberId: string }) => {
      return await removeRoomMemberAction(roomId, memberId)
    },
    onSuccess: (res, variables) => {
      if (res.success) {
        queryClient.invalidateQueries({ queryKey: ['room-members', variables.roomId] })
      }
    },
  })
}
