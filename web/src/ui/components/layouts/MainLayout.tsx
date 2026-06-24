import { auth } from '@/workspace/auth'

import DesktopNav from '../modules/Nav/DesktopNav'

type Props = Readonly<{
  children: React.ReactNode
}>

export default async function MainLayout({ children }: Props) {
  const session = await auth()
  const isAdmin = session?.user?.role === 'Admin'

  return (
    <main
      className="flex flex-row min-h-screen bg-background text-foreground relative z-0 transition-colors duration-300"
    >
      <DesktopNav isAdmin={isAdmin} />
      <div className="flex flex-1 p-4">{children}</div>
    </main>
  )
}
