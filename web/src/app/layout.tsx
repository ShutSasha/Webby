import type { Metadata } from 'next'

import NotificationPopupContainer from '@/ui/components/modules/Toast/NotificationPopupContainer'
import ToastContainer from '@/ui/components/modules/Toast/ToastContainer'
import Providers from '@/ui/components/providers/Providers'
import { inter } from '@/ui/fonts'

import '@/ui/global.css'

export const metadata: Metadata = {
  title: 'Webby',
}

type Props = Readonly<{
  children: React.ReactNode
}>

export default async function RootLayout({ children }: Props) {
  return (
    <html lang="uk">
      <body className={`${inter.className} antialiased`}>
        <Providers>
          {children}
          <NotificationPopupContainer />
        </Providers>
        <ToastContainer />
      </body>
    </html>
  )
}
