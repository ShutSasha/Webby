import ChatInput from './ChatInput'
import RoomChat from './RoomChat'
import RoomInteractionFooter from '../Footer/RoomInteractionFooter'

type Props = {
  chatId: string
}

export default function RoomChatContainer({ chatId }: Props) {
  return (
    <div className="flex flex-col flex-1 min-h-0 gap-2">
      <RoomChat chatId={chatId} />
      <ChatInput chatId={chatId} />
      <RoomInteractionFooter />
    </div>
  )
}
