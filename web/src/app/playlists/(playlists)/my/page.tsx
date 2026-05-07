import LockIcon from '@/assets/icons/shared/lock.svg'
import PlaylistUserSearchContainer from '@/ui/components/modules/Playlists/PlaylistUserSearchContainer'
import AuthPlaceholder from '@/ui/components/shared/AuthPlaceholder'
import { auth } from '@/workspace/auth'

type Props = {
  searchParams: Promise<{ query?: string }>
}

export default async function Page({ searchParams }: Props) {
  const { query } = await searchParams
  const safeQuery = query || ''
  const session = await auth()

  if (!session?.user) {
    return (
      <AuthPlaceholder
        title="Sign in to view your playlists"
        description=" Keep track of your favorite videos, create custom collections, and manage your saved content by logging into
        your account."
        icon={<LockIcon className="size-10 text-neutral-500 stroke-1" />}
      />
    )
  }

  return <PlaylistUserSearchContainer query={safeQuery} userId={session.user.id} />
}
