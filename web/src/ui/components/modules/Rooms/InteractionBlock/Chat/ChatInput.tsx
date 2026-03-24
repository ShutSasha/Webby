'use client'

import SendIcon from '@/assets/icons/shared/send_message.svg'

export default function ChatInput() {
  return (
    <div className="relative">
      <input
        placeholder="Send a message"
        className="py-2 pl-4 pr-12 bg-neutral-900 placeholder:text-neutral-700 w-full ring-0 outline-0 rounded-lg
          text-neutral-300"
      />
      <SendIcon
        className="size-6 text-neutral-700 absolute -translate-y-1/2 top-1/2 right-3 cursor-pointer
          hover:text-emerald-500 transition-colors duration-300 ease-in-out"
      />
    </div>
  )
}
