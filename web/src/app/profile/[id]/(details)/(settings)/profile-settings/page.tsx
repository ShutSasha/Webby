import { redirect } from 'next/navigation'

import AboutContainer from '@/ui/components/modules/Profile/Settings/ProfileSettings/AboutContainer'
import UploadAvatarContainer from '@/ui/components/modules/Profile/Settings/ProfileSettings/UploadAvatarContainer'
import UserHeader from '@/ui/components/modules/Profile/Settings/UserHeader'
import { auth } from '@/workspace/auth'

export default async function ProfileSettings() {
  const session = await auth()

  if (!session?.user) {
    redirect('/login')
  }

  return (
    <div className="flex flex-col">
      <UserHeader image={session.user.image} username={session.user.username} />
      <div className="flex justify-between items-center my-4">
        <p>User email: test@gmail.com</p>
        <UploadAvatarContainer />
      </div>
      <p className="mb-2">About:</p>
      <AboutContainer />
    </div>
  )
}
