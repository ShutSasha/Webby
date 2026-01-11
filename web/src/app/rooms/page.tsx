import MainLayout from '@/ui/components/MainLayout'
import Search from '@/ui/components/Search'

export default function RoomPage() {
  return (
    <MainLayout>
      <div className="w-full bg-neutral-900 rounded-[20px] p-5 gap-4 box-border">
        <Search placeholder="Search a room..." />
      </div>
    </MainLayout>
  )
}
