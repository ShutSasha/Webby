import { useParams } from 'next/navigation'

import CubeIcon from '@/assets/icons/Room/cube-points.svg'
import { useGetRoomMemberPointsQuery } from '@/lib/hooks/api/room/useGetRoomMemberPoints'
import { formatPoints } from '@/lib/utils/video.utils'

import InteractionButton from './InteractionButton'

export default function UserPoints() {
  const params = useParams()
  const roomId = params?.id as string | undefined

  const { data, isLoading } = useGetRoomMemberPointsQuery(roomId)
  const points = data?.points ?? 0

  return (
    <InteractionButton
      icon={
        <CubeIcon
          className="size-5 text-foreground-subtle dark:group-hover:text-foreground-inverse-subtle transition-colors
            duration-300 ease-in-out"
        />
      }
      text={isLoading ? <span className="animate-pulse">...</span> : formatPoints(points)}
    />
  )
}
