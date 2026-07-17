'use client'
import { useState } from 'react'

import PlusIcon from '@/assets/icons/ic_plus_create.svg'
import { usePlaylistStore } from '@/stores/playlist.store'
import { MediaType } from '@/types/general.types'
import ActionButton from '@/ui/components/shared/ActionButton'

import SaveToPlaylistModal from '../SaveToPlaylistModal'

type ButtonProps = {
  videoId: string
  userId: string
  mediaType?: MediaType | null
}

export default function SaveToPlaylistButton({ videoId, userId, mediaType }: ButtonProps) {
  const [isOpen, setIsOpen] = useState(false)
  const activeMediaType = usePlaylistStore(state => state.activeMediaType)

  const currentMediaType = mediaType || activeMediaType

  return (
    <>
      <ActionButton onClick={() => setIsOpen(true)} label="Add to playlist" btnClassName="self-end">
        <PlusIcon className="size-4" />
      </ActionButton>

      <SaveToPlaylistModal
        isOpen={isOpen}
        onClose={() => setIsOpen(false)}
        videoId={videoId}
        userId={userId}
        mediaType={currentMediaType}
      />
    </>
  )
}
