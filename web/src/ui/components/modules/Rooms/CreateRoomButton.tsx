'use client'

import MediaButton from '../../shared/MediaButton'

export default function CreateRoomButton() {
  const handleClick = () => {
    alert('Create room logic here')
  }

  return (
    <MediaButton actionLabel="Create a room">
      <h3 className="text-neutral-300 text-center mb-2">Create a room</h3>
      <button className="bg-emerald-500 text-neutral-900" onClick={handleClick}>
        click me
      </button>
    </MediaButton>
  )
}
