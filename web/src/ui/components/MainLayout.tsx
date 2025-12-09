import SideNav from './SideNav'

type Props = Readonly<{
  children: React.ReactNode
}>

export default function MainLayout({ children }: Props) {
  return (
    <main className="flex min-h-screen bg-neutral-800 text-neutral-300">
      <SideNav />
      <div className="w-full p-4">{children}</div>
    </main>
  )
}
