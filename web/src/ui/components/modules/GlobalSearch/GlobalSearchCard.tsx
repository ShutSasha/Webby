import { memo, useState } from 'react'

import Link from 'next/link'
import { useSession } from 'next-auth/react'

import { EntityType, FALLBACK_MAP, getRoute } from '@/lib/utils/global-search-modal.utils'
import { BLUR_DATA_URLS } from '@/ui/images'

import { SearchCardMenu } from './SearchCardMenu'
import SafeImage from '../../shared/SafeImage'
import SaveToPlaylistModal from '../Playlists/SaveToPlaylistModal'

type Props = {
  id: string
  title: string
  subtitle: string
  thumbnail: string
  type: EntityType
  showAddButton?: boolean
  onCloseSearchModal: () => void
}

const GlobalSearchCard = ({
  id,
  title,
  subtitle,
  thumbnail,
  type,
  showAddButton = true,
  onCloseSearchModal,
}: Props) => {
  const { data: session } = useSession()
  const currentUserId = session?.user?.id

  const [isPlaylistModalOpen, setIsPlaylistModalOpen] = useState(false)

  const checkHasMenu = (type: EntityType) => type !== 'User' && type !== 'Room'

  const handleNavigateTo = () => {
    onCloseSearchModal()
  }

  const targetUrl = getRoute(type, id)

  return (
    <>
      <Link
        href={targetUrl}
        className="group flex items-center justify-between gap-3 p-2 rounded-xl hover:bg-neutral-800/40
          transition-colors border border-transparent hover:border-neutral-800/60 cursor-pointer"
        onClick={handleNavigateTo}
      >
        <div className="flex items-center gap-3 overflow-hidden">
          <div
            className={`relative shrink-0 overflow-hidden bg-neutral-800 flex items-center justify-center
              ${type === 'User' ? 'w-12 h-12 rounded-full' : 'w-24 h-14 rounded-lg'}`}
          >
            <SafeImage
              src={thumbnail}
              alt={title}
              fallbackType={FALLBACK_MAP[type]}
              width={160}
              height={90}
              className="object-cover aspect-video h-full w-full"
              placeholder="blur"
              blurDataURL={BLUR_DATA_URLS['neutral800']}
            />
          </div>
          <div className="flex flex-col overflow-hidden">
            <h4 className="text-sm font-semibold text-neutral-200 truncate">{title}</h4>
            <p className="text-xs text-neutral-500 truncate mt-0.5">{subtitle}</p>
          </div>
        </div>

        <div className="flex items-center gap-1 shrink-0">
          {showAddButton && checkHasMenu(type) && (
            <SearchCardMenu type={type} id={id} onOpenPlaylistModal={() => setIsPlaylistModalOpen(true)} />
          )}
        </div>
      </Link>

      {isPlaylistModalOpen && (
        <SaveToPlaylistModal
          isOpen={isPlaylistModalOpen}
          onClose={() => setIsPlaylistModalOpen(false)}
          videoId={id}
          userId={currentUserId}
        />
      )}
    </>
  )
}

export default memo(GlobalSearchCard)

export function GlobalSearchCardSkeleton({ isUser = false }: { isUser?: boolean }) {
  return (
    <div className="flex items-center gap-3 p-2 w-full animate-pulse">
      <div className={`shrink-0 bg-neutral-800/80 ${isUser ? 'w-12 h-12 rounded-full' : 'w-24 h-14 rounded-lg'}`} />
      <div className="flex flex-col gap-2 w-full">
        <div className="h-3.5 bg-neutral-800/80 rounded w-[60%]" />
        <div className="h-3 bg-neutral-800/60 rounded w-[40%]" />
      </div>
    </div>
  )
}
