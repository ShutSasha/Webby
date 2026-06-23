import Link from 'next/link'

const activeStyle = 'text-foreground-strong cursor-default'
const inactiveStyle = 'text-foreground0 hover:text-emerald-400 transition-colors duration-300'

type Props = {
  id: string
}

export async function FollowersToggle({ id }: Props) {
  return (
    <div className="flex items-center justify-center gap-2 text-[20px] font-semibold">
      <span className={activeStyle}>Followers</span>
      <span className="text-foreground-disabled">/</span>
      <Link href={`/profile/${id}/follows`} className={inactiveStyle}>
        Follows
      </Link>
    </div>
  )
}

export async function FollowsToggle({ id }: Props) {
  return (
    <div className="flex items-center justify-center gap-2 text-[20px] font-semibold">
      <Link href={`/profile/${id}/followers`} className={inactiveStyle}>
        Followers
      </Link>
      <span className="text-foreground-disabled">/</span>
      <span className={activeStyle}>Follows</span>
    </div>
  )
}
