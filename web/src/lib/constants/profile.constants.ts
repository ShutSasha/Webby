export type ProfileTabValue = 'video' | 'playlist'

export const PROFILE_TABS_CONFIG: { label: string; value: ProfileTabValue }[] = [
  { label: 'Video', value: 'video' },
  { label: 'Playlist', value: 'playlist' },
]

export const ALLOWED_PROFILE_TABS = PROFILE_TABS_CONFIG.map(tab => tab.value)
