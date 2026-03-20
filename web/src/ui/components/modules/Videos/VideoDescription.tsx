'use client'

import { useState, useRef, useLayoutEffect } from 'react'

import { cn } from '@/lib/utils/utils'

export default function VideoDescription({ text }: { text: string }) {
  const [isExpanded, setIsExpanded] = useState(false)
  const [isTruncated, setIsTruncated] = useState(false)
  const textRef = useRef<HTMLParagraphElement>(null)

  useLayoutEffect(() => {
    const element = textRef.current
    if (element) {
      const hasOverflow = element.scrollHeight > element.clientHeight
      setIsTruncated(hasOverflow)
    }
  }, [text])

  return (
    <div className="flex flex-col gap-2 rounded-xl bg-black/40 p-3 mt-4 overflow-hidden">
      <p className="text-sm font-bold text-neutral-300">100 views | 17/02/2026</p>

      <div className="relative">
        <p
          ref={textRef}
          className={cn(
            'text-sm text-neutral-300 break-all leading-relaxed transition-all',
            !isExpanded && 'line-clamp-3',
          )}
        >
          {text}
        </p>

        {(isTruncated || isExpanded) && (
          <button
            onClick={() => setIsExpanded(!isExpanded)}
            className="text-sm font-bold text-neutral-300 mt-1 hover:underline cursor-pointer transition-all block"
          >
            {isExpanded ? 'Show less' : '...more'}
          </button>
        )}
      </div>
    </div>
  )
}
