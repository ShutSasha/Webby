import LockIcon from '@/assets/icons/shared/lock.svg'
import ResetPasswordForm from '@/ui/components/modules/Profile/Settings/Secure/ResetPasswordForm'
import UserHeader from '@/ui/components/modules/Profile/Settings/UserHeader'
import AuthPlaceholder from '@/ui/components/shared/AuthPlaceholder'
import { auth } from '@/workspace/auth'

export default async function SecurePage() {
  const session = await auth()

  if (!session) {
    return (
      <AuthPlaceholder
        title="Sign in to access your settings"
        description="Please log in to manage your account details, update your public profile description,
        customize your display settings, and configure security preferences to keep your account safe."
        icon={<LockIcon className="size-10 text-foreground0 stroke-1" />}
      />
    )
  }

  return (
    <div className="flex flex-col">
      <UserHeader image={session.user.image} username={session.user.username} />
      <ResetPasswordForm />
    </div>
  )
}
