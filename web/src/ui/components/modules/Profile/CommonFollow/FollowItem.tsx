import Link from 'next/link'

import MailIcon from '@/assets/icons/ic_mail_with_background.svg'
import { useCreateChatMutation } from '@/lib/hooks/api/chat/useCreateChat'
import SafeImage from '@/ui/components/shared/SafeImage'
import { BLUR_DATA_URLS } from '@/ui/images'

import UnfollowButton from './UnfollowButton'

type Props = {
  userId: string
  avatarUrl: string
  followersCount: number
  username: string
  isFollow: boolean
  userProfileId: string
  currentUserId?: string
}

export default function FollowItem({
  userId,
  avatarUrl,
  followersCount,
  username,
  isFollow,
  userProfileId,
  currentUserId,
}: Props) {
  const isOwner = userProfileId === currentUserId
  const { mutate: createChat, isPending } = useCreateChatMutation()

  const handleCreateChat = () => {
    if (isPending) return
    createChat(userId)
  }

  return (
    <div
      className="bg-neutral-800 hover:bg-neutral-300/10 transition-all duration-200 ease-in rounded-xl p-2 flex
        items-center justify-between gap-2 border border-transparent hover:border-emerald-400/25"
    >
      <Link href={`/profile/${userId}`} className="flex items-center gap-2">
        {/* User info */}
        <SafeImage
          src={avatarUrl}
          alt=""
          width={60}
          height={60}
          className="w-13 h-13 rounded-full"
          loading="lazy"
          placeholder="blur"
          blurDataURL={BLUR_DATA_URLS['neutral900']}
        />
        <div className="flex flex-col">
          <p className="text-[16px] font-medium">{username}</p>
          <p className="text-foreground0 text-[12px]">{followersCount} followers</p>
        </div>
      </Link>

      {/* Actions */}
      <div className="flex items-center gap-5 pr-2">
        {currentUserId !== userId && (
          <div
            className={`relative flex items-center justify-center
            ${isPending ? 'cursor-not-allowed opacity-70' : 'cursor-pointer group/mail'}`}
            onClick={handleCreateChat}
          >
            {isPending ? (
              <div
                className="w-[18px] h-[18px] rounded-full border-2 border-neutral-600 border-t-emerald-500 animate-spin"
              />
            ) : (
              <>
                <MailIcon
                  className="w-4.5 h-4.5 text-foreground0/90 group-hover/mail:text-emerald-500/90 transition-all
                    duration-300"
                />
                <div
                  className="absolute w-8 h-8 left-1/2 -translate-x-1/2 top-1/2 -translate-y-1/2
                    group-hover/mail:bg-emerald-500/20 transition-all duration-300 rounded-full"
                />
              </>
            )}
          </div>
        )}
        {isFollow && isOwner && <UnfollowButton targetId={userId} currentUserId={currentUserId} />}
      </div>
    </div>
  )
}
