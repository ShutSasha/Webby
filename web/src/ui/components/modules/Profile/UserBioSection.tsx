import FollowConnections from './FollowConnections'

type UserBioSectionProps = {
  username: string
  bio?: string
  userId: string
  followersCount: number
  followsCount: number
}

export default function UserBioSection({ username, bio, userId, followersCount, followsCount }: UserBioSectionProps) {
  return (
    <div className="flex flex-col justify-between">
      <div className="flex flex-col">
        <h2 className="font-medium text-[24px] text-white leading-tight">{username}</h2>
        {bio && <p className="text-sm max-w-[400px] text-neutral-400">{bio}</p>}
      </div>

      <FollowConnections followersCount={followersCount} followsCount={followsCount} userId={userId} />
    </div>
  )
}
