'use client'

import { useGetRoomByIdQuery } from '@/lib/hooks/api/room/useGetRoomByIdQuery'

type Props = {
  roomId: string
  initialName: string
}

export default function RoomTitle({ roomId, initialName }: Props) {
  const { data: room } = useGetRoomByIdQuery(roomId)

  const displayName = room?.name || initialName

  return <p className="text-foreground-secondary text-[20px] font-bold truncate">{displayName}</p>
}
