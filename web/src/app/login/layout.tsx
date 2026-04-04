import type { Metadata } from 'next'

import MainLayout from '@/ui/components/layouts/MainLayout'

export const metadata: Metadata = {
  title: 'Webby Login',
}

type Props = Readonly<{
  children: React.ReactNode
}>

export default async function SignUpLayout({ children }: Props) {
  return <MainLayout>{children}</MainLayout>
}
