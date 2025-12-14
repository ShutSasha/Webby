import { ComponentProps, FC, SVGProps } from 'react'

import Link from 'next/link'

import MailIcon from '@/assets/icons/ic_mail.svg'
import NotificationIcon from '@/assets/icons/ic_notifications.svg'
import PlaylistIcon from '@/assets/icons/ic_playlist.svg'
import SearchIcon from '@/assets/icons/ic_search.svg'

export type NavItem = {
  href: ComponentProps<typeof Link>['href']
  icon: FC<SVGProps<SVGSVGElement>>
  text: string
}

export const navItems: NavItem[] = [
  { href: '/', icon: SearchIcon, text: 'Search' },
  { href: '/', icon: PlaylistIcon, text: 'Playlists' },
  { href: '/', icon: MailIcon, text: 'Messages' },
  { href: '/', icon: NotificationIcon, text: 'Notifications' },
]
