import { RefCallback } from 'react'

import Image from 'next/image'
import Link from 'next/link'

import { AdminVideoRecord } from '@/lib/actions/admin.actions'
import { useBanVideoMutation } from '@/lib/hooks/api/admin/useBanVideo'
import { useUnbanVideoMutation } from '@/lib/hooks/api/admin/useUnbanVideo'
import { formatDate } from '@/lib/utils/date.utils'

type VideoTableRowProps = {
  video: AdminVideoRecord
  lastElementRef: RefCallback<HTMLTableRowElement> | null
}

export default function VideoTableRow({ video, lastElementRef }: VideoTableRowProps) {
  const { mutate: banVideo, isPending: isBanning } = useBanVideoMutation()
  const { mutate: unbanVideo, isPending: isUnbanning } = useUnbanVideoMutation()

  const isPending = isBanning || isUnbanning
  const status = video.isBanned ? 'Banned' : 'Available'

  return (
    <tr
      ref={lastElementRef}
      className={`hover:bg-background/20 transition-colors ${isPending ? 'opacity-60 pointer-events-none' : ''}`}
    >
      <td className="py-4 px-6">
        <div className="flex items-center gap-4">
          <div
            className="w-28 h-16 rounded-md bg-emerald-500/20 text-emerald-500 flex items-center justify-center
              font-bold overflow-hidden shrink-0 relative"
          >
            {video.previewUrl ? (
              <Image src={video.previewUrl} alt={video.name} fill className="object-cover" sizes="112px" />
            ) : (
              <span className="text-xs">No image</span>
            )}
          </div>

          <div className="flex flex-col max-w-[300px]">
            <Link
              href={`/videos/wb_${video.videoId}`}
              className="font-medium text-foreground-tertiary truncate hover:text-foreground-subtle"
              title={video.name}
            >
              {video.name}
            </Link>
            <Link
              href={`/profile/${video.user.userId}`}
              className="text-xs text-foreground-faint truncate w-fit hover:text-foreground-disabled"
              title={video.user.username}
            >
              {video.user.username}
            </Link>
          </div>
        </div>
      </td>

      <td className="py-4 px-6">
        <div className="flex items-center gap-2">
          <div className={`size-2 rounded-full ${status === 'Available' ? 'bg-emerald-500' : 'bg-red-500'}`} />
          <span className={status === 'Available' ? 'text-foreground-subtle' : 'text-red-400'}>{status}</span>
        </div>
      </td>

      <td className="py-4 px-6 text-foreground-faint text-sm">{formatDate(video.createdAt)}</td>

      <td className="py-4 px-6 flex justify-end items-center gap-2">
        {status === 'Available' ? (
          <button
            onClick={() => banVideo(video.videoId)}
            disabled={isPending}
            className="w-[68px] flex items-center justify-center px-3 py-1.5 text-xs font-semibold bg-red-500/10
              hover:bg-red-500/20 text-red-500 border border-red-500/20 rounded-lg transition-colors
              disabled:opacity-50"
          >
            {isBanning ? (
              <div className="size-3.5 border-2 border-red-500/30 border-t-red-500 rounded-full animate-spin" />
            ) : (
              'Ban'
            )}
          </button>
        ) : (
          <button
            onClick={() => unbanVideo(video.videoId)}
            disabled={isPending}
            className="w-[68px] flex items-center justify-center px-3 py-1.5 text-xs font-semibold bg-emerald-500/10
              hover:bg-emerald-500/20 text-emerald-500 border border-emerald-500/20 rounded-lg transition-colors
              disabled:opacity-50"
          >
            {isUnbanning ? (
              <div className="size-3.5 border-2 border-emerald-500/30 border-t-emerald-500 rounded-full animate-spin" />
            ) : (
              'Unban'
            )}
          </button>
        )}
      </td>
    </tr>
  )
}
