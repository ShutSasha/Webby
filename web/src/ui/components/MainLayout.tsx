import SideNav from './modules/Nav/SideNav'

type Props = Readonly<{
  children: React.ReactNode
}>

export default function MainLayout({ children }: Props) {
  return (
    <main className="flex min-h-screen bg-neutral-800 text-neutral-300">
      <SideNav />
      <div className="flex flex-1 p-4">{children}</div>
    </main>
  )
}
