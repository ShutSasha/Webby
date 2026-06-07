'use client'
import { useState } from 'react'

import PlusIcon from '@/assets/icons/ic_plus_create.svg'
import { usePlaylistStore } from '@/stores/playlist.store'
import ActionButton from '@/ui/components/shared/ActionButton'

import SaveToPlaylistModal from '../SaveToPlaylistModal'

type ButtonProps = {
  videoId: string
  userId: string
}

export default function SaveToPlaylistButton({ videoId, userId }: ButtonProps) {
  const [isOpen, setIsOpen] = useState(false)
  const activeMediaType = usePlaylistStore(state => state.activeMediaType)

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
        mediaType={activeMediaType}
      />
    </>
  )
}
