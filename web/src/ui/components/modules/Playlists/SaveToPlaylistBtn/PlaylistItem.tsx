import Image from 'next/image'

import { BLUR_DATA_URLS } from '@/ui/images'

type Props = {
  playlistId: string
  image: string
  name: string
  count: number
  status: 'added' | 'removed' | 'none'
  disabled?: boolean
  onToggle: (playlistId: string) => void
}

export default function PlaylistItem({ playlistId, image, name, count, status, disabled, onToggle }: Props) {
  return (
    <div
      className={`flex items-center justify-between py-2 px-2.5 transition-all duration-300 ease-in-out rounded-xl
        ${disabled ? 'opacity-50 cursor-not-allowed' : 'group cursor-pointer hover:bg-neutral-800/50'} `}
      onClick={() => {
        if (!disabled) onToggle(playlistId)
      }}
    >
      <div className="flex gap-4">
        <Image
          src={image}
          width={100}
          height={100}
          alt=""
          className="aspect-square size-10 object-cover rounded-lg"
          loading="lazy"
          placeholder="blur"
          blurDataURL={BLUR_DATA_URLS['neutral800']}
        />
        <div className="flex flex-col">
          <p className="text-sm text-foreground-subtle line-clamp-1" title={name}>
            {name}
          </p>
          <p className="text-sm text-foreground0">{count === 1 ? '1 video' : `${count} videos`}</p>
        </div>
      </div>
      <div className="flex items-center justify-center size-6 shrink-0 ml-3">
        {status === 'added' && (
          <div
            className="size-5 bg-emerald-500 rounded-full flex items-center justify-center
              shadow-[0_0_8px_rgba(16,185,129,0.3)] animate-in zoom-in duration-200"
          >
            <svg
              className="w-3.5 h-3.5 text-foreground-inverse-subtle"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              strokeWidth={3}
            >
              <path strokeLinecap="round" strokeLinejoin="round" d="M5 13l4 4L19 7" />
            </svg>
          </div>
        )}
        {status === 'removed' && (
          <div
            className="size-5 bg-red-500 rounded-full flex items-center justify-center
              shadow-[0_0_8px_rgba(239,68,68,0.3)] animate-in zoom-in duration-200"
          >
            <svg
              className="w-3.5 h-3.5 text-foreground-inverse-subtle"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              strokeWidth={3.5}
            >
              <path strokeLinecap="round" strokeLinejoin="round" d="M5 12h14" />
            </svg>
          </div>
        )}
        {status === 'none' && (
          <div
            className="size-5 border-2 border-neutral-600 rounded-full group-hover:border-neutral-400 transition-colors"
          />
        )}
      </div>
    </div>
  )
}
