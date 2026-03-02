import { notFound, redirect } from 'next/navigation'

import SettingsNavigation from '@/ui/components/modules/Profile/Settings/SettingsNavigation'
import { auth } from '@/workspace/auth'

type Props = {
  children: React.ReactNode
  params: Promise<{ id: string }>
}

export default async function Layout({ children, params }: Props) {
  const { id } = await params
  const session = await auth()
  await new Promise(r => setTimeout(r, 700))

  if (!session?.user) {
    redirect('/login')
  }

  if (session.user.id !== id) {
    notFound()
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
