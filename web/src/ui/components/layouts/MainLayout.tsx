import DesktopNav from '../modules/Nav/DesktopNav'

type Props = Readonly<{
  children: React.ReactNode
}>

export default function MainLayout({ children }: Props) {
  return (
    <main className="flex flex-row min-h-screen bg-neutral-800 text-neutral-300 relative z-0">
      <DesktopNav />
      <div className="flex flex-1 p-4">{children}</div>
    </main>
  )
}
