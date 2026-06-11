import ChatArea from '@/ui/components/modules/Chats/ChatArea'

type Props = {
  params: Promise<{ chatId: string }>
}

export default async function ChatPage({ params }: Props) {
  const { chatId } = await params

  return <ChatArea />
}
