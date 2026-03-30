'use client'

import MediaButton from '../../shared/MediaButton'

export default function CreateVideoButton() {
  const handleClick = () => {
    alert('Create video logic here')
  }

  return (
    <MediaButton>
      <h3 className="text-neutral-300 text-center mb-2">Create a video</h3>
      <button className="bg-emerald-500 text-neutral-900" onClick={handleClick}>
        click me
      </button>
    </MediaButton>
  )
}
