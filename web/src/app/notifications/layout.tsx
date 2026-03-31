import type { Metadata } from 'next'

import MainContainer from '@/ui/components/layouts/MainContainer'
import MainLayout from '@/ui/components/layouts/MainLayout'

export const metadata: Metadata = {
  title: 'Notifications',
}

type Props = Readonly<{
  children: React.ReactNode
}>

export default async function Layout({ children }: Props) {
  return (
    <MainLayout>
      <MainContainer>{children}</MainContainer>
    </MainLayout>
  )
}
