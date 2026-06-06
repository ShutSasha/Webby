import { useMutation, useQueryClient } from '@tanstack/react-query'

import { addRoomMembersAction } from '@/lib/actions/room.actions'

export const useAddRoomMembersMutation = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async ({ roomId, userIds }: { roomId: string; userIds: string[] }) => {
      return await addRoomMembersAction(roomId, userIds)
    },
    onSuccess: (res, variables) => {
      if (res.success) {
        queryClient.invalidateQueries({ queryKey: ['room-members', variables.roomId] })
      }
    },
  })
}
