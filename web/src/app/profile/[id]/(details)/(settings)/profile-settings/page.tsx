import Link from 'next/link'

import LockIcon from '@/assets/icons/shared/lock.svg'
import { getUser } from '@/lib/actions/user.actions'
import AboutContainer from '@/ui/components/modules/Profile/Settings/ProfileSettings/AboutContainer'
import UploadAvatarContainer from '@/ui/components/modules/Profile/Settings/ProfileSettings/UploadAvatarContainer'
import UserHeader from '@/ui/components/modules/Profile/Settings/UserHeader'
import AuthPlaceholder from '@/ui/components/shared/AuthPlaceholder'
import EmptyState from '@/ui/components/shared/EmptyState'
import { auth } from '@/workspace/auth'

type Props = {
  params: Promise<{ id: string }>
}

export default async function ProfileSettings({ params }: Props) {
  const { id } = await params
  const session = await auth()
  const userData = await getUser(id)

  if (!session) {
    return (
      <AuthPlaceholder
        title="Sign in to access your settings"
        description="Please log in to manage your account details, update your public profile description,
        customize your display settings, and configure security preferences to keep your account safe."
        icon={<LockIcon className="size-10 text-neutral-500 stroke-1" />}
      />
    )
  }

  if (!userData) {
    return (
      <div className="flex-1 flex items-center justify-center bg-neutral-900/20 rounded-[20px] min-h-[50vh]">
        <EmptyState
          title="Couldn't load profile data"
          description="We ran into an issue while fetching your settings. Please try refreshing the page."
        />
      </div>
    )
  }

  return (
    <div className="flex flex-col w-full max-w-3xl mx-auto gap-8 pb-10">
      <UserHeader image={session.user.image} username={session.user.username} />

      <div className="flex flex-col bg-neutral-900/40 p-6 sm:p-8 rounded-3xl border border-neutral-800 shadow-sm">
        <h2 className="text-sm font-bold text-neutral-500 uppercase tracking-wider mb-6">Account Details</h2>

        <div className="flex flex-col sm:flex-row sm:justify-between sm:items-center gap-6">
          <div className="flex flex-col gap-1">
            <span className="text-sm text-neutral-400">Email address</span>
            <span className="text-neutral-200 font-medium text-[16px]">{userData.user.email}</span>
          </div>

          <UploadAvatarContainer />
        </div>

        <div className="mt-8 pt-6 border-t border-neutral-800 flex flex-col sm:flex-row justify-between items-center
          gap-4">
          <div className="flex flex-col">
            <span className="text-neutral-200 font-medium">Subscription</span>
            <span className="text-sm text-neutral-400">Upgrade to unlock exclusive features.</span>
          </div>
          <Link
            href="/premium"
            className="px-5 py-2.5 rounded-xl bg-linear-to-r from-emerald-500/10 to-emerald-500/5 border
              border-emerald-500/20 text-emerald-400 font-semibold text-sm hover:from-emerald-500/20
              hover:to-emerald-500/10 hover:border-emerald-500/40 transition-all duration-300"
          >
            Get Premium
          </Link>
        </div>
      </div>

      <div className="flex flex-col bg-neutral-900/40 p-6 sm:p-8 rounded-3xl border border-neutral-800 shadow-sm">
        <h2 className="text-sm font-bold text-neutral-500 uppercase tracking-wider mb-6">About Me</h2>
        <AboutContainer userId={id} about={userData.user.about} />
      </div>
    </div>
  )
}
