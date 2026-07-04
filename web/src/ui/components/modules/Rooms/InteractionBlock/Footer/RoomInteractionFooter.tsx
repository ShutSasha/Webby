'use client'

import Poll from './Poll'
import Reactions from './Reactions'
import UserPoints from './UserPoints'

export default function RoomInteractionFooter() {
  return (
    <div className="flex flex-row justify-between items-center relative">
      <div className="flex items-center gap-2">
        <UserPoints />
        <Reactions />
      </div>
      <Poll />
    </div>
  )
}
