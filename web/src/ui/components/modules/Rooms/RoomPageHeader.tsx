import CreateRoomButton from './CreateRoomButton'
import Search from '../../Search'

export default function RoomPageHeader() {
  return (
    <div
      className="flex flex-col md:flex-row items-start md:items-center justify-between gap-4 pb-4 border-b
        border-neutral-800/60"
    >
      <div className="flex flex-col shrink-0 order-1">
        <h1 className="text-2xl font-bold text-neutral-100">Live Rooms</h1>
        <p className="text-sm text-neutral-400 mt-0.5">Watch and discuss videos together</p>
      </div>

      <div className="w-full md:max-w-md order-3 md:order-2">
        <Search placeholder="Search a room..." containerClassName="w-full" />
      </div>

      <div className="shrink-0 order-2 md:order-3 self-end md:self-auto">
        <CreateRoomButton />
      </div>
    </div>
  )
}
