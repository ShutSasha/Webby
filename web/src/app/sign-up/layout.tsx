import type { Metadata } from 'next'

import MainLayout from '@/ui/components/MainLayout'

export const metadata: Metadata = {
  title: 'Webby Sign Up',
}

type Props = Readonly<{
  children: React.ReactNode
}>

export default async function SignUpLayout({ children }: Props) {
  return <MainLayout>{children}</MainLayout>
}
