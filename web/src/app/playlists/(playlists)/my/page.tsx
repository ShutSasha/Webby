import PlaylistAuthPlaceholder from '@/ui/components/modules/Playlists/PlaylistAuthPlaceholder'
import PlaylistUserSearchContainer from '@/ui/components/modules/Playlists/PlaylistUserSearchContainer'
import { auth } from '@/workspace/auth'

type Props = {
  searchParams: Promise<{ query?: string }>
}

export default async function Page({ searchParams }: Props) {
  const { query } = await searchParams
  const safeQuery = query || ''
  const session = await auth()

  if (!session?.user) {
    return <PlaylistAuthPlaceholder />
  }

  return <PlaylistUserSearchContainer query={safeQuery} userId={session.user.id} />
}
