import MainContainer from '@/ui/components/layouts/MainContainer'
import MainLayout from '@/ui/components/layouts/MainLayout'

type Props = Readonly<{
  children: React.ReactNode
}>

export default function PaymentLayout({ children }: Props) {
  return (
    <MainLayout>
      <MainContainer>{children}</MainContainer>
    </MainLayout>
  )
}
