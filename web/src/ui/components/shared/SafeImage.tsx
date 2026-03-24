'use client'

import { useState, useEffect } from 'react'

import Image, { ImageProps } from 'next/image'

const DEFAULT_IMAGE = 'https://i.ibb.co/4RLdNrBC/785ca39a2a95c19e66b01b3e0615d32c.jpg'

interface Props extends Omit<ImageProps, 'src'> {
  src: string | null | undefined
}

// TODO: think about useEffect logic, it triggers so many re-renders for many reusable components
export default function SafeImage({ src, alt, ...props }: Props) {
  const validateSrc = (url: string | null | undefined): string => {
    if (!url || typeof url !== 'string') return DEFAULT_IMAGE

    const trimmed = url.trim()

    const isValidFormat = trimmed.startsWith('/') || trimmed.startsWith('http')

    const isGarbage = trimmed === 'null' || trimmed === 'undefined' || trimmed === ''

    return !isGarbage && isValidFormat ? trimmed : DEFAULT_IMAGE
  }

  const [imgSrc, setImgSrc] = useState(() => validateSrc(src))

  useEffect(() => {
    setImgSrc(validateSrc(src))
  }, [src])

  return <Image {...props} src={imgSrc} alt={alt || 'image'} onError={() => setImgSrc(DEFAULT_IMAGE)} />
}
