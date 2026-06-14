import { Metadata } from 'next'

import MainContainer from '@/ui/components/layouts/MainContainer'
import MainLayout from '@/ui/components/layouts/MainLayout'
import ChatSidebar from '@/ui/components/modules/Chats/ChatSidebar'

export const metadata: Metadata = {
  title: 'Chats',
}

type Props = Readonly<{
  children: React.ReactNode
}>

export default function Layout({ children }: Props) {
  return (
    <MainLayout>
      <MainContainer className="flex flex-row h-[calc(100vh-40px)]">
        <ChatSidebar />

        <div className="flex-1 rounded-2xl overflow-hidden flex flex-col border border-neutral-800/60">{children}</div>
      </MainContainer>
    </MainLayout>
  )
}
