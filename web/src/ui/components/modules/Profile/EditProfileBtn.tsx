import Link from 'next/link'

import SettingsIcon from '@/assets/icons/Profile/ic_settings.svg'

export default function EditProfileBtn({ userId }: { userId: string }) {
  return (
    <div className="flex self-end">
      <Link
        href={`/profile/${userId}/profile-settings`}
        className="flex items-center gap-2 py-2 px-4 border border-neutral-700 text-sm rounded-full cursor-pointer
          bg-transparent text-foreground-subtle hover:text-foreground-strong hover:border-neutral-500
          hover:bg-background/50 transition-all duration-300 ease-out"
      >
        <SettingsIcon className="w-4 h-4" />
        <span className="font-medium">Edit Profile</span>
      </Link>
    </div>
  )
}
