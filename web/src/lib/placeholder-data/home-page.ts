import { FC, SVGProps } from 'react'

import FileIcon from '@/assets/icons/HomePage/ic_file.svg'
import FilmIcon from '@/assets/icons/HomePage/ic_film.svg'
import LinkIcon from '@/assets/icons/HomePage/ic_link.svg'
import SlideshareIcon from '@/assets/icons/HomePage/ic_slideshare.svg'
import YoutubeIcon from '@/assets/icons/HomePage/ic_youtube.svg'

type Benefit = {
  icon: FC<SVGProps<SVGSVGElement>>
  title: string
  description: string
}

export const benefitsList: Benefit[] = [
  {
    icon: FilmIcon,
    title: 'movies',
    description: 'Watch movies and series with friends in the distance.',
  },
  {
    icon: YoutubeIcon,
    title: 'services',
    description: 'Watch YouTube, Twitch and other services with your friends.',
  },
  {
    icon: SlideshareIcon,
    title: 'screen share',
    description: 'Share your browser or application tab.',
  },
  {
    icon: FileIcon,
    title: 'files',
    description: 'View video files from your device.',
  },
  {
    icon: LinkIcon,
    title: 'url',
    description: 'Watch live video from the web server.',
  },
]
