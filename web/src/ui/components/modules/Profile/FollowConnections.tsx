import Link from 'next/link'

type FollowConnectionsProps = {
  userId: string
  followersCount?: number
  followsCount?: number
}

export default async function FollowConnections({
  userId,
  followersCount = 0,
  followsCount = 0,
}: FollowConnectionsProps) {
  const linkStyles = 'underline hover:text-emerald-400 transition-colors duration-300'

  return (
    <div className="flex gap-2 items-center text-sm">
      <Link href={`/profile/${userId}/followers`} className={linkStyles}>
        {followersCount} Followers
      </Link>

      <span className="bg-neutral-300 w-1 h-1 rounded-full" />

      <Link href={`/profile/${userId}/follows`} className={linkStyles}>
        {followsCount} Follows
      </Link>
    </div>
  )
}

export function FollowConnectionsSkeleton() {
  return (
    <div className="flex gap-2 items-center">
      {/* Followers placeholder */}
      <div className="h-4 w-20 bg-neutral-800 animate-pulse rounded-md" />

      {/* Dot separator */}
      <span className="bg-neutral-700 w-1 h-1 rounded-full" />

      {/* Follows placeholder */}
      <div className="h-4 w-16 bg-neutral-800 animate-pulse rounded-md" />
    </div>
  )
}
