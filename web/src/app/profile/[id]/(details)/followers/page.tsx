import { getUser } from '@/lib/actions/user.actions'
import FollowersContainer from '@/ui/components/modules/Profile/CommonFollow/FollowersContainer'
import { FollowersToggle } from '@/ui/components/modules/Profile/CommonFollow/FollowToggle'
import HeaderNavigation from '@/ui/components/modules/Profile/CommonFollow/HeaderNavigation'
import EmptyState from '@/ui/components/shared/EmptyState'

type Props = {
  params: Promise<{ id: string }>
}

export default async function Followers({ params }: Props) {
  const { id } = await params
  const user = await getUser(id)

  if (!user) {
    return (
      <div className="flex flex-col h-full w-full">
        <EmptyState
          title="User not found"
          description="The user you are looking for does not exist or has been deleted."
        />
      </div>
    )
  }

  return (
    <div className="flex flex-col gap-2">
      <HeaderNavigation id={id} image={user.user.avatarUrl} username={user.user.username}>
        <FollowersToggle id={id} />
      </HeaderNavigation>

      <hr className="border-emerald-400/50 mb-1" />

      <FollowersContainer id={id} />
    </div>
  )
}
