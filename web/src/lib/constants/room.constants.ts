export type RoomTabValue = 'all' | 'cinema' | 'music' | 'gaming' | 'education'

export const ROOM_TABS_CONFIG: { label: string; value: RoomTabValue }[] = [
  { label: 'All', value: 'all' },
  { label: 'Cinema', value: 'cinema' },
  { label: 'Music', value: 'music' },
  { label: 'Gaming', value: 'gaming' },
  { label: 'Education', value: 'education' },
]
