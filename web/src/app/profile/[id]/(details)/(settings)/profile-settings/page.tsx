import { redirect } from 'next/navigation'

import { getUser } from '@/app/api/user'
import AboutContainer from '@/ui/components/modules/Profile/Settings/ProfileSettings/AboutContainer'
import UploadAvatarContainer from '@/ui/components/modules/Profile/Settings/ProfileSettings/UploadAvatarContainer'
import UserHeader from '@/ui/components/modules/Profile/Settings/UserHeader'
import EmptyState from '@/ui/components/shared/EmptyState'
import { auth } from '@/workspace/auth'

type Props = {
  params: Promise<{ id: string }>
}

export default async function ProfileSettings({ params }: Props) {
  const { id } = await params
  const session = await auth()
  const userData = await getUser(id)

  if (!session?.user) {
    redirect('/login')
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
    <div className="flex flex-col">
      <UserHeader image={session.user.image} username={session.user.username} />
      <div className="flex justify-between items-center my-4">
        <p>User email: {userData.user.email}</p>
        <UploadAvatarContainer />
      </div>
      <p className="mb-2">About:</p>
      <AboutContainer userId={id} about={userData.user.about} />
    </div>
  )
}
