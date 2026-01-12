import React from 'react'

import { s } from 'framer-motion/client'

export default function RoomCard() {
  return (
    <div className="rounded-xl h-[180px] py-3 px-2.5">
      <div className="flex items-center justify-between">
        {/* Live */}
        <div className="bg-red-700 flex items-center gap-1 px-2 py-1 rounded-lg">
          <span className="h-1.5 w-1.5 rounded-full text-neutral-300 block" />
          <p className="uppercase font-black text-[9px] leading-3 bg-neutral-300">live</p>
        </div>
      </div>
    </div>
  )
}
