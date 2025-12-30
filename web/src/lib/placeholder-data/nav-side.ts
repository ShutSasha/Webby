import { ComponentProps, FC, SVGProps } from 'react'

import Link from 'next/link'

import LogoIcon from '@/assets/icons/ic_logo.svg'
import MailIcon from '@/assets/icons/ic_mail.svg'
import NotificationIcon from '@/assets/icons/ic_notifications.svg'
import PlaylistIcon from '@/assets/icons/ic_playlist.svg'
import SearchIcon from '@/assets/icons/ic_search.svg'

export type NavItem = {
  key: string
  href: ComponentProps<typeof Link>['href']
  icon: FC<SVGProps<SVGSVGElement>>
  text: string
  iconSize: string
}

export const navItems: NavItem[] = [
  { key: 'search', href: '/', icon: SearchIcon, text: 'Search', iconSize: 'size-6' },
  { key: 'playlists', href: '/', icon: PlaylistIcon, text: 'Playlists', iconSize: 'size-6' },
  { key: 'messages', href: '/', icon: MailIcon, text: 'Messages', iconSize: 'size-6' },
  { key: 'notifications', href: '/', icon: NotificationIcon, text: 'Notifications', iconSize: 'size-6' },
]

export const navMobileItems: NavItem[] = [
  { key: 'search', href: '/', icon: SearchIcon, text: 'Search', iconSize: 'max-[420px]:size-5 size-6' },
  { key: 'playlists', href: '/', icon: PlaylistIcon, text: 'Playlists', iconSize: 'max-[420px]:size-5 size-6' },
  { key: 'logo', href: '/', icon: LogoIcon, text: '', iconSize: 'max-[420px]:size-10 size-12' },
  { key: 'messages', href: '/', icon: MailIcon, text: 'Messages', iconSize: 'max-[420px]:size-5 size-6' },
  {
    key: 'notifications',
    href: '/',
    icon: NotificationIcon,
    text: 'Notifications',
    iconSize: 'max-[420px]:size-5 size-6',
  },
]
