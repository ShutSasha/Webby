'use client'

import { useState } from 'react'

import Image, { ImageProps } from 'next/image'

const FALLBACK_IMAGES = {
  user: '/fallbacks/default-user-avatar.jpg',
  video: '/fallbacks/video-not-found.jpg',
  room: '/fallbacks/room-not-found.jpg',
  playlist: '/fallbacks/playlist-not-found.jpg',
  default: '/fallbacks/no-image.png',
} as const

export type FallbackType = keyof typeof FALLBACK_IMAGES

interface SafeImageProps extends Omit<ImageProps, 'src'> {
  src?: string | null
  fallbackType?: FallbackType
}

const checkIsValidSrc = (url?: string | null): boolean => {
  if (!url || typeof url !== 'string') return false

  const trimmed = url.trim()
  if (trimmed === 'null' || trimmed === 'undefined' || trimmed === '') return false

  if (trimmed.startsWith('/')) return true

  try {
    new URL(trimmed)
    return true
  } catch {
    return false
  }
}

export default function SafeImage({ src, alt, fallbackType = 'default', ...props }: SafeImageProps) {
  const [failedSrc, setFailedSrc] = useState<string | null>(null)

  const isValid = checkIsValidSrc(src)

  const shouldUseFallback = !isValid || src === failedSrc

  const finalSrc = shouldUseFallback ? FALLBACK_IMAGES[fallbackType] : (src as string)

  return (
    <Image
      {...props}
      src={finalSrc}
      alt={alt || 'Image'}
      onError={() => {
        if (src) setFailedSrc(src)
      }}
    />
  )
}
