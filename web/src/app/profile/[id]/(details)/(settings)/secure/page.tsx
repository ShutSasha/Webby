import { redirect } from 'next/navigation'

import ResetPasswordForm from '@/ui/components/modules/Profile/Settings/Secure/ResetPasswordForm'
import UserHeader from '@/ui/components/modules/Profile/Settings/UserHeader'
import { auth } from '@/workspace/auth'

export default async function SecurePage() {
  const session = await auth()

  if (!session?.user) {
    redirect('/login')
  }

  return (
    <div className="flex flex-col">
      <UserHeader image={session.user.image} username={session.user.username} />
      <ResetPasswordForm userId={session.user.id} />
    </div>
  )
}
