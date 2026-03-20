type Props = {
  username: string
  message: string
}

export default function RoomChatMessage({ username, message }: Props) {
  return (
    <div>
      <span className="text-cyan-500">{username}:</span>
      <span className="text-neutral-300">{message}</span>
    </div>
  )
}
