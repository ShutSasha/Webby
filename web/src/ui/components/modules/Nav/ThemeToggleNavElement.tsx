'use client'

import { useTheme } from 'next-themes'

import MoonIcon from '@/assets/icons/Nav/moon.svg'
import SunIcon from '@/assets/icons/Nav/sun.svg'
import { useIsClient } from '@/lib/hooks/useIsClient'

import { DesktopNavElement } from './DesktopNavElement'

const LoadingIcon = ({ className }: { className?: string }) => (
  <div className={`rounded-full bg-background animate-pulse ${className}`} />
)

type Props = {
  isExpanded: boolean
}

export function ThemeToggleNavElement({ isExpanded }: Props) {
  const { theme, setTheme } = useTheme()
  const isMounted = useIsClient()

  const toggleTheme = () => {
    setTheme(theme === 'dark' ? 'light' : 'dark')
  }

  if (!isMounted) {
    return (
      <DesktopNavElement
        isExpanded={isExpanded}
        text="Theme"
        Icon={LoadingIcon}
        iconSize="size-5 2xl:size-6"
        onClick={() => {}}
      />
    )
  }

  const isDark = theme === 'dark'

  return (
    <DesktopNavElement
      isExpanded={isExpanded}
      text={isDark ? 'Light mode' : 'Dark mode'}
      Icon={isDark ? SunIcon : MoonIcon}
      iconSize="size-5 2xl:size-6"
      onClick={toggleTheme}
    />
  )
}
