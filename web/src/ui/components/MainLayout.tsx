import DesktopNav from './modules/Nav/DesktopNav'
import MobileNav from './modules/Nav/MobileNav'

type Props = Readonly<{
  children: React.ReactNode
}>

export default function MainLayout({ children }: Props) {
  return (
    <main className="flex flex-col-reverse md:flex-row min-h-screen bg-neutral-800 text-neutral-300 relative">
      <DesktopNav />
      <MobileNav />
      <div className="flex flex-1 p-4">{children}</div>
    </main>
  )
}
