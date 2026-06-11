'use client'

import Image from 'next/image'

import SendIcon from '@/assets/icons/shared/send_message.svg'

const MOCK_MESSAGES = [
  { id: 1, text: 'Worem ipsum dolor sit amet, consectetur adipiscing elit.', isMe: false, time: '17:50' },
  { id: 2, text: 'Lorem ipsum dolor sit', isMe: false, time: '17:50' },
  { id: 3, text: 'Lorem ipsum dolor sit', isMe: false, time: '17:50' },
  { id: 4, text: 'Lorem ipsum dolor sit', isMe: true, time: '17:50' },
  { id: 5, text: 'Lorem ipsum dolor sit', isMe: true, time: '17:50' },
  { id: 6, text: 'Worem ipsum dolor sit amet, consectetur adipiscing elit.', isMe: false, time: '17:50' },
  { id: 7, text: 'Lorem ipsum dolor sit', isMe: false, time: '17:50' },
]

export default function ChatArea() {
  return (
    <div className="flex flex-col h-full w-full">
      {/* Header */}
      <div className="h-[72px] shrink-0 border-b border-neutral-800/60 flex items-center justify-between px-6 bg-neutral-900/20">
        <div className="flex items-center gap-3">
          <Image
            src="https://i.pravatar.cc/150?u=header"
            alt="Avatar"
            width={40}
            height={40}
            className="rounded-full object-cover size-10"
          />
          <span className="font-medium text-neutral-200">@guntersteam</span>
        </div>
      </div>

      {/* Messages Body */}
      <div className="flex-1 overflow-y-auto custom-scrollbar p-6 flex flex-col gap-2">
        {MOCK_MESSAGES.map((msg) => (
          <div
            key={msg.id}
            className={`flex flex-col max-w-[70%] ${msg.isMe ? 'self-end items-end' : 'self-start items-start'}`}
          >
            <div
              className={`px-4 py-2.5 flex items-end gap-3 shadow-sm ${
                msg.isMe
                  ? 'bg-[#1c1c1c] text-neutral-300 rounded-2xl rounded-br-sm border border-neutral-800/50' // Відправлені (сірі)
                  : 'bg-emerald-500 text-neutral-950 rounded-2xl rounded-bl-sm font-medium' // Отримані (зелені)
              }`}
            >
              <p className="text-[15px] leading-relaxed">{msg.text}</p>
              <span
                className={`text-[10px] shrink-0 translate-y-0.5 ${msg.isMe ? 'text-neutral-600' : 'text-emerald-900/60'}`}
              >
                {msg.time}
              </span>
            </div>
          </div>
        ))}
      </div>

      {/* Input Area */}
      <div className="p-4 bg-neutral-900/20 border-t border-neutral-800/60">
        <div className="relative flex items-center">
          <input
            type="text"
            placeholder="Write a message..."
            className="w-full bg-[#141414] border border-neutral-800 text-neutral-200 placeholder:text-neutral-600 rounded-xl py-3.5 pl-5 pr-12 outline-none focus:border-neutral-600 transition-colors"
          />
          <button className="absolute right-3 p-1.5 text-neutral-500 hover:text-emerald-500 transition-colors">
            <SendIcon className="size-5" />
          </button>
        </div>
      </div>
    </div>
  )
}