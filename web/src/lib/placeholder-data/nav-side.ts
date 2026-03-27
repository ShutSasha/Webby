import { ComponentProps, FC, SVGProps } from 'react'

import Link from 'next/link'

import LogoIcon from '@/assets/icons/ic_logo.svg'
import MenuIcon from '@/assets/icons/ic_menu.svg'
import SearchIcon from '@/assets/icons/ic_search.svg'
import BellIcon from '@/assets/icons/Nav/bell.svg'
import ChatsIcon from '@/assets/icons/Nav/chats.svg'
import PlaylistIcon from '@/assets/icons/Nav/list-video.svg'
import VideosIcon from '@/assets/icons/Nav/tv-minimal-play.svg'
import RoomsIcon from '@/assets/icons/Nav/tv.svg'

export type NavItem = {
  key: string
  href: ComponentProps<typeof Link>['href']
  icon: FC<SVGProps<SVGSVGElement>>
  text: string
  iconSize: string
}

export const navDesktopElements: NavItem[] = [
  { key: 'search', href: '/', icon: SearchIcon, text: 'Search', iconSize: 'size-6' },
  { key: 'videos', href: '/videos', icon: VideosIcon, text: 'Videos', iconSize: 'size-6' },
  { key: 'rooms', href: '/rooms', icon: RoomsIcon, text: 'Rooms', iconSize: 'size-6' },
  { key: 'playlists', href: '/playlists', icon: PlaylistIcon, text: 'Playlists', iconSize: 'size-6' },
  { key: 'chats', href: '/chats', icon: ChatsIcon, text: 'Chats', iconSize: 'size-6' },
  { key: 'notifications', href: '/notifications', icon: BellIcon, text: 'Notifications', iconSize: 'size-6' },
]

export const navMobileElements: NavItem[] = [
  { key: 'search', href: '/rooms', icon: SearchIcon, text: 'Search', iconSize: 'max-[420px]:size-5 size-6' },
  { key: 'playlists', href: '/', icon: PlaylistIcon, text: 'Playlists', iconSize: 'max-[420px]:size-5 size-6' },
  { key: 'logo', href: '/', icon: LogoIcon, text: '', iconSize: 'max-[420px]:size-10 size-12' },
  { key: 'chats', href: '/', icon: ChatsIcon, text: 'Chats', iconSize: 'max-[420px]:size-5 size-6' },
  {
    key: 'menu',
    href: '/',
    icon: MenuIcon,
    text: 'Menu',
    iconSize: 'max-[420px]:size-5 size-6',
  },
]
