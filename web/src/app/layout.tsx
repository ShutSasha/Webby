import type { Metadata } from 'next'
import { Geist, Geist_Mono } from 'next/font/google'
import './globals.css'
import { auth } from '../../auth'
import Providers from './components/Providers'

const geistSans = Geist({
  variable: '--font-geist-sans',
  subsets: ['latin'],
})

const geistMono = Geist_Mono({
  variable: '--font-geist-mono',
  subsets: ['latin'],
})

export const metadata: Metadata = {
  title: 'Secure Shop',
  description: 'Diffie-Hellman & DES implementation',
}

export default async function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode
}>) {
  const session = await auth()

  return (
    <html lang="uk">
      <Providers>
        <body className={`${geistSans.variable} ${geistMono.variable} bg-gray-900 text-gray-200 antialiased`}>
          <div className="flex min-h-[calc(100vh-64px)] flex-col">{children}</div>
        </body>
      </Providers>
    </html>
  )
}
