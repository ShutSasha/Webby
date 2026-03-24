import RoomChatMessage from './RoomChatMessage'

export default function RoomChat() {
  return (
    <div
      className="flex-1 overflow-y-auto min-h-0 flex flex-col gap-1 pr-1 [&::-webkit-scrollbar]:w-1.5
        [&::-webkit-scrollbar-track]:bg-transparent [&::-webkit-scrollbar-thumb]:bg-neutral-900
        [&::-webkit-scrollbar-thumb]:border-0 [&::-webkit-scrollbar-thumb]:rounded-full
        hover:[&::-webkit-scrollbar-thumb]:bg-neutral-800"
    >
      {[...new Array(20)].map((_, key) => (
        <RoomChatMessage
          key={key}
          username="temma0101"
          message="u just scared cuz u know how good i would be with a knife"
        />
      ))}
    </div>
  )
}
