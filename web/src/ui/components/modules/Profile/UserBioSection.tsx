import FollowConnections, { FollowConnectionsSkeleton } from './FollowConnections'

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

export function UserBioSectionSkeleton() {
  return (
    <div className="flex flex-col justify-between gap-2">
      <div className="flex flex-col gap-2">
        {/* Username placeholder */}
        <div className="h-7 w-32 bg-neutral-800 animate-pulse rounded-md" />

        {/* Bio text placeholder */}
        <div className="h-4 w-48 md:w-64 bg-neutral-800 animate-pulse rounded-md" />
      </div>

      {/* Підключаємо скелетон підписок */}
      <FollowConnectionsSkeleton />
    </div>
  )
}
