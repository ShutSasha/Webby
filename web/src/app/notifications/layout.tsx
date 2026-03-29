import type { Metadata } from 'next'

import MainLayout from '@/ui/components/layouts/MainLayout'

export const metadata: Metadata = {
  title: 'Notifications',
}

type Props = Readonly<{
  children: React.ReactNode
}>

export default async function Layout({ children }: Props) {
  return <MainLayout>{children}</MainLayout>
}
