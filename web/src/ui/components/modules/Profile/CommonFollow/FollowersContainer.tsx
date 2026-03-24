import { getUserFollowers } from '@/app/api/user'

import FollowItem from './FollowItem'

type Props = {
  id: string
}

export default async function FollowersContainer({ id }: Props) {
  const followers = await getUserFollowers(id)

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 2xl:grid-cols-3 gap-4">
      {followers &&
        followers.length > 0 &&
        followers.map(follower => <FollowItem key={follower.userId} {...follower} />)}
    </div>
  )
}
