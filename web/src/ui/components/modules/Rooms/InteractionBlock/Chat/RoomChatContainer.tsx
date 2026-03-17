import RoomChat from './RoomChat'

export default function RoomChatContainer() {
  return (
    <div className="flex flex-col flex-1 min-h-0 gap-2">
      <RoomChat />
      <div className="bg-red-400">input</div>
      <div className="bg-purple-600">RoomInteractionFooter</div>
    </div>
  )
}
