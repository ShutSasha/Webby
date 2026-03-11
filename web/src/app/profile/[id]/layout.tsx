import MainLayout from '@/ui/components/MainLayout'

type Props = {
  children: React.ReactNode
  params: Promise<{ id: string }>
}

export default async function Layout({ children }: Props) {
  return (
    <MainLayout>
      <div className="flex flex-col gap-4 w-full max-w-5xl 2xl:max-w-7xl mx-auto">{children}</div>
    </MainLayout>
  )
}
