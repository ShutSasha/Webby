import VideosIcon from '@/assets/icons/Nav/tv-minimal-play.svg'
import UsersIcon from '@/assets/icons/Notifications/users.svg'
import SparklesIcon from '@/assets/icons/Premium/sparkles.svg'

export const PREMIUM_FEATURES = [
  {
    id: 'capacity',
    title: 'Expanded Room Capacity',
    icon: UsersIcon,
    description:
      'Gather large groups for watch parties or mass streams. Forget about limits and invite all your friends at once!',
    freeLimit: 'Limit of 8 participants per room.',
    premiumLimit: 'Increased limit up to 50 participants.',
  },
  {
    id: 'customization',
    title: 'Customization & Higher Limits',
    icon: SparklesIcon,
    description:
      'Stand out from the crowd! Make your rooms and playlists unique with animated thumbnails and high-quality images.',
    freeLimit: 'Standard static images (up to 2 MB).',
    premiumLimit: 'Ability to upload animated .gif and files up to 5 MB.',
  },
  {
    id: 'uploads',
    title: 'Increased Video Storage',
    icon: VideosIcon,
    description:
      'Create massive personal content libraries. More space means more of your favorite videos always ready for co-watching.',
    freeLimit: 'Upload up to 10 videos per account.',
    premiumLimit: 'Upload up to 100 videos.',
  },
]
