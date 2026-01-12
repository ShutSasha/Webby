import CategoryToggle from '@/ui/components/common/GenreToggle'
import MediaHeader from '@/ui/components/common/MediaHeader'
import MainLayout from '@/ui/components/MainLayout'
import CreateRoomButton from '@/ui/components/modules/Rooms/CreateRoomButton'

export default function RoomPage() {
  return (
    <MainLayout>
      <div className="flex flex-col w-full bg-neutral-900 rounded-[20px] p-5 gap-4 box-border">
        <MediaHeader searchPlaceholder="Search rooms..." actionSlot={<CreateRoomButton />} />

        <div className="flex justify-between">
          <div className='hidden lg:block'>
            <div className="flex gap-3 items-center">
              <span className="w-10 h-0.5 bg-emerald-500 rounded-full" />
              <p className="text-emerald-500 font-black text-[10px] leading-3.5 uppercase">live now</p>
            </div>
            <h2 className="font-black text-[48px] leading-14 text-emerald-500 uppercase">Rooms</h2>
          </div>
          <CategoryToggle />
        </div>
      </div>
    </MainLayout>
  )
}
