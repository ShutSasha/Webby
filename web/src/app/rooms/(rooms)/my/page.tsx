import LockIcon from '@/assets/icons/shared/lock.svg'
import UserRoomsContainer from '@/ui/components/modules/Rooms/UserRoomsContainer'
import AuthPlaceholder from '@/ui/components/shared/AuthPlaceholder'
import { auth } from '@/workspace/auth'

type Props = {
  searchParams: Promise<{ query?: string }>
}

export default async function Page({ searchParams }: Props) {
  const { query } = await searchParams
  const safeQuery = query || ''
  const session = await auth()

  if (!session) {
    return (
      <AuthPlaceholder
        title="Sign in to view your rooms"
        description="Manage your custom rooms, adjust privacy settings, and host synchronized viewing sessions by logging into your account."
        icon={<LockIcon className="size-10 text-neutral-500 stroke-1" />}
      />
    )
  }

  return <UserRoomsContainer query={safeQuery} />
}
