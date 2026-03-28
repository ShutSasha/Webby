import MainLayout from '@/ui/components/MainLayout'

type Props = Readonly<{
  children: React.ReactNode
}>

export default function Layout({ children }: Props) {
  return (
    <MainLayout>
      <div className="flex flex-col w-full bg-neutral-900 rounded-[20px] p-5 gap-4 box-border">{children}</div>
    </MainLayout>
  )
}
