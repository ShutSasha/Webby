import PlaylistPublicSearchContainer from '@/ui/components/modules/Playlists/PlaylistPublicSearchContainer'

type Props = {
  searchParams: Promise<{ query?: string }>
}

export default async function Page({ searchParams }: Props) {
  const { query } = await searchParams
  const safeQuery = query || ''

  return <PlaylistPublicSearchContainer query={safeQuery} />
}
