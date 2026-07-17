import { useMemo } from 'react'

export function isPlatformSource(url: string): boolean {
  if (!url) return false

  const isNativeFile = /\.(mp4|webm)($|\?)/i.test(url)

  return !isNativeFile
}

export function usePlayerControls(url: string): boolean {
  const isThirdPartyPlatform = useMemo(() => isPlatformSource(url), [url])

  return isThirdPartyPlatform
}
