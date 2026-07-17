import type { Metadata } from 'next'
import NextTopLoader from 'nextjs-toploader'

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
    <html lang="uk" suppressHydrationWarning>
      <body className={`${inter.className} antialiased`}>
        <NextTopLoader
          color="#10b981"
          initialPosition={0.4}
          crawlSpeed={200}
          height={3}
          crawl={true}
          showSpinner={false}
          easing="ease-out"
          speed={400}
          shadow="none"
          zIndex={99999}
        />

        <Providers>
          {children}
          <NotificationPopupContainer />
        </Providers>
        <ToastContainer />
      </body>
    </html>
  )
}
