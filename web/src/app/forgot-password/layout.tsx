import { Metadata } from 'next'

import MainLayout from '@/ui/components/layouts/MainLayout'

export const metadata: Metadata = {
  title: 'Forgot password',
}

type Props = Readonly<{
  children: React.ReactNode
}>

export default function Layout({ children }: Props) {
  return <MainLayout>{children}</MainLayout>
}
