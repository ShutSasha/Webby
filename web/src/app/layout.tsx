import type { Metadata } from 'next'

import Providers from '@/ui/components/Providers/Providers'
import { inter } from '@/ui/fonts'
import './globals.css'

export const metadata: Metadata = {
  title: 'Webby',
}

type Props = Readonly<{
  children: React.ReactNode
}>

export default async function RootLayout({ children }: Props) {
  return (
    <html lang="uk">
      <Providers>
        <body className={`${inter.className} antialiased`}>{children}</body>
      </Providers>
    </html>
  )
}
