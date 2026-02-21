import Link from 'next/link'

import SettingsIcon from '@/assets/icons/Profile/ic_settings.svg'

export default function EditProfileBtn({ userId }: { userId: string }) {
  return (
    <div className="flex group/settings self-end">
      <Link
        href={`/profile/${userId}/settings`}
        className="flex items-center gap-2 py-2 px-3 border border-border text-sm rounded-full cursor-pointer
          group-hover/settings:bg-emerald-400 transition-all duration-400 ease-in"
      >
        <SettingsIcon className="w-4 h-4 group-hover/settings:text-neutral-900 transition-all duration-300 font-medium" />
        <p className="group-hover/settings:text-neutral-900 transition-all duration-300 font-medium">Edit Profile</p>
      </Link>
    </div>
  )
}
