import ChatArea from '@/ui/components/modules/Chats/ChatArea'
import ChatWebSocketManager from '@/ui/components/modules/Chats/ChatWebSocketManager' // Вкажи правильний шлях до створеного менеджера

type Props = {
  params: Promise<{ chatId: string }>
}

export default async function ChatPage({ params }: Props) {
  const { chatId } = await params

  return (
    <>
      <ChatWebSocketManager chatId={chatId} />
      <ChatArea chatId={chatId} />
    </>
  )
}
