import FollowConnections from './FollowConnections'

type UserBioSectionProps = {
  username: string
  bio?: string
  userId: string
}

export default function UserBioSection({ username, bio, userId }: UserBioSectionProps) {
  return (
    <div className="flex flex-col justify-between">
      <div className="flex flex-col">
        <h2 className="font-medium text-[24px] text-white leading-tight">{username}</h2>
        {bio && <p className="text-sm max-w-[400px] text-neutral-400">{bio}</p>}
      </div>

      <FollowConnections userId={userId} />
    </div>
  )
}
