import Link from 'next/link'

type FollowConnectionsProps = {
  userId: string
  followersCount: number
  followsCount: number
}

export default function FollowConnections({ userId, followersCount, followsCount }: FollowConnectionsProps) {
  const linkStyles = 'underline hover:text-emerald-400 transition-colors duration-300 text-foreground-subtle'

  return (
    <div className="flex gap-2 items-center text-sm">
      <Link href={`/profile/${userId}/followers`} className={linkStyles}>
        {followersCount} Followers
      </Link>

      <span className="w-1 h-1 rounded-full bg-foreground-subtle" />

      <Link href={`/profile/${userId}/follows`} className={linkStyles}>
        {followsCount} Follows
      </Link>
    </div>
  )
}
