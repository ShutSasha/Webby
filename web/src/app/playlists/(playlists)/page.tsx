import PlaylistSearchContainer from '@/ui/components/modules/Playlists/PlaylistSearchContainer'

type Props = {
  searchParams: Promise<{ query?: string }>
}

export default async function Page({ searchParams }: Props) {
  const { query } = await searchParams
  const safeQuery = query || ''

  return <PlaylistSearchContainer query={safeQuery} />
}
