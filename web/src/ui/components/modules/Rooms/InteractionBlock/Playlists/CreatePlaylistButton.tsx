'use client'

import MediaButton from '@/ui/components/common/MediaButton'

export default function CreatePlaylistButton() {
  const handleClick = () => {
    alert('Create room logic here')
  }

  return (
    <MediaButton actionLabel="Create a playlist">
      <h3 className="text-neutral-300 text-center mb-2">Create a playlist</h3>
      <button className="bg-emerald-500 text-neutral-900" onClick={handleClick}>
        click me
      </button>
    </MediaButton>
  )
}
