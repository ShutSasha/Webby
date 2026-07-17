import HiddenEye from '@/assets/auth/ic_eye_off.svg'

type Props = {
  playlistName: string
  hiddenVideosCount: number
}

export default function PlaylistContainerHeader({ playlistName, hiddenVideosCount }: Props) {
  return (
    <div className="flex flex-col mb-2 px-1 shrink-0">
      <h2 className="text-foreground-secondary font-bold text-xl line-clamp-2 leading-tight" title={playlistName}>
        {playlistName}
      </h2>

      {hiddenVideosCount > 0 && (
        <div className="flex items-center gap-1.5 text-foreground-faint text-xs font-medium bg-surface/50 w-fit
          rounded-md">
          <HiddenEye className="size-3.5 shrink-0" />

          <span>
            {hiddenVideosCount} {hiddenVideosCount > 1 ? 'videos' : 'video'} hidden
          </span>
        </div>
      )}
    </div>
  )
}
