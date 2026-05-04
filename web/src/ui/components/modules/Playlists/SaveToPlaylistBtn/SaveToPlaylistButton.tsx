'use client'
import { useState } from 'react'

import PlusIcon from '@/assets/icons/ic_plus_create.svg'
import { VideoSource } from '@/types/video.types'
import ActionButton from '@/ui/components/shared/ActionButton'

import SaveToPlaylistModal from '../SaveToPlaylistModal'

type ButtonProps = {
  videoId: string
  userId: string
  videoSource: VideoSource | undefined
}

export default function SaveToPlaylistButton({ videoId, userId, videoSource }: ButtonProps) {
  const [isOpen, setIsOpen] = useState(false)

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
        videoSource={videoSource}
      />
    </>
  )
}
