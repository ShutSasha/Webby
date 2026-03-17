import RoomInteractionFooter from '../RoomInteractionFooter'
import ChatInput from './ChatInput'
import RoomChat from './RoomChat'

export default function RoomChatContainer() {
  return (
    <div className="flex flex-col flex-1 min-h-0 gap-2">
      <RoomChat />
      <ChatInput />
      <RoomInteractionFooter />
    </div>
  )
}
