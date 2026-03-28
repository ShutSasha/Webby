import { redirect } from 'next/navigation'

import SettingsNavigation from '@/ui/components/modules/Profile/Settings/SettingsNavigation'
import EmptyState from '@/ui/components/shared/EmptyState'
import { auth } from '@/workspace/auth'

type Props = {
  children: React.ReactNode
  params: Promise<{ id: string }>
}

export default async function Layout({ children, params }: Props) {
  const [{ id }, session] = await Promise.all([params, auth()])

  if (!session?.user) {
    redirect('/login')
  }

  if (session.user.id !== id) {
    return (
      <div className="flex-1 flex items-center justify-center bg-neutral-900/20 rounded-[20px]">
        <EmptyState
          title="Access Denied"
          description="You can only view and manage the settings for your own profile."
        />
      </div>
    )
  }

  return (
    <div className="flex flex-1 gap-2">
      {/* Left menu */}
      <SettingsNavigation id={id} />
      <span className="w-px bg-border rounded-full" />
      <div className="w-full">{children}</div>
    </div>
  )
}
