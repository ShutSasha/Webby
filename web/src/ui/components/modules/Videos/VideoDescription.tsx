'use client'

import { useState, useRef, useLayoutEffect } from 'react'

import { formatDate } from '@/lib/utils/date.utils'
import { cn } from '@/lib/utils/general.utils'

type Props = {
  text: string
  views: number
  date: string
  videoTags?: string[] | null
}

export default function VideoDescription(props: Props) {
  const [isExpanded, setIsExpanded] = useState(false)
  const [isTruncated, setIsTruncated] = useState(false)
  const textRef = useRef<HTMLParagraphElement>(null)

  useLayoutEffect(() => {
    const element = textRef.current
    if (element) {
      const hasOverflow = element.scrollHeight > element.clientHeight
      setIsTruncated(hasOverflow)
    }
  }, [props.text, props.videoTags])

  return (
    <div className="flex flex-col gap-2 rounded-xl bg-black/40 p-3 mt-4 overflow-hidden">
      <p className="text-sm font-bold text-neutral-300">
        {props.views} {props.views === 1 ? 'view' : 'views'} | {formatDate(props.date)}
      </p>

      <div className="relative">
        <p
          ref={textRef}
          className={cn(
            'text-sm text-neutral-300 wrap-break-word whitespace-pre-wrap leading-relaxed transition-all',
            !isExpanded && 'line-clamp-3',
          )}
        >
          {props.text}

          {props.videoTags && props.videoTags.length > 0 && (
            <>
              {'\n\n'}
              {props.videoTags.map(tag => (
                <span key={tag} className="text-neutral-500 font-medium mr-2 transition-colors">
                  #{tag}
                </span>
              ))}
            </>
          )}
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
