export type RoomTabValue = 'all' | 'cinema' | 'music' | 'games' | 'education'

export const ROOM_TABS_CONFIG: { label: string; value: RoomTabValue }[] = [
  { label: 'All', value: 'all' },
  { label: 'Cinema', value: 'cinema' },
  { label: 'Music', value: 'music' },
  { label: 'Games', value: 'games' },
  { label: 'Education', value: 'education' },
]

export const ALLOWED_ROOM_TABS = ROOM_TABS_CONFIG.map(tab => tab.value)
