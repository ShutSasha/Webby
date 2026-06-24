'use client'

import { useState, useEffect, useMemo } from 'react'

import { useIsClient } from '@/lib/hooks/useIsClient'

type UseSafeVideoTransitionProps = {
  videoUrl: string
  delayMs?: number
}

type UseSafeVideoTransitionResult = {
  activeVideoUrl: string
  isTransitioning: boolean
  isActuallyReady: boolean
}

export const useSafeVideoTransition = ({
  videoUrl,
  delayMs = 800,
}: UseSafeVideoTransitionProps): UseSafeVideoTransitionResult => {
  const isMounted = useIsClient()

  const [activeVideoUrl, setActiveVideoUrl] = useState<string>(videoUrl)
  const [isTransitioning, setIsTransitioning] = useState<boolean>(false)

  useEffect(() => {
    if (videoUrl !== activeVideoUrl) {
      setIsTransitioning(true)

      const timer = setTimeout(() => {
        setActiveVideoUrl(videoUrl)
        setIsTransitioning(false)
      }, delayMs)

      return () => clearTimeout(timer)
    } else {
      setIsTransitioning(false)
    }
  }, [videoUrl, activeVideoUrl, delayMs])

  return useMemo(() => {
    const isActuallyReady = isMounted && !isTransitioning && videoUrl === activeVideoUrl

    return {
      activeVideoUrl,
      isTransitioning,
      isActuallyReady,
    }
  }, [activeVideoUrl, isTransitioning, isMounted, videoUrl])
}
