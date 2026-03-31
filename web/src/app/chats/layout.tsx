import { Metadata } from 'next'

import MainContainer from '@/ui/components/layouts/MainContainer'
import MainLayout from '@/ui/components/layouts/MainLayout'

export const metadata: Metadata = {
  title: 'Chats',
}

type Props = Readonly<{
  children: React.ReactNode
}>

export default function Layout({ children }: Props) {
  return (
    <MainLayout>
      <MainContainer>{children}</MainContainer>
    </MainLayout>
  )
}
