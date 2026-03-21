import MediaHeader from '@/ui/components/common/MediaHeader'
import MainLayout from '@/ui/components/MainLayout'
import PageToggle from '@/ui/components/modules/Playlists/PageToggle'
import CreateRoomButton from '@/ui/components/modules/Rooms/CreateRoomButton'

type Props = Readonly<{
  children: React.ReactNode
}>

export default function Layout({ children }: Props) {
  return (
    <MainLayout>
      <div className="flex flex-col w-full bg-neutral-900 rounded-[20px] p-5 gap-4 box-border">
        {/* Create playlist btn */}
        <MediaHeader
          searchPlaceholder="Search a playlist"
          actionSlot={<CreateRoomButton />}
          pageToggle={<PageToggle />}
        />
        {children}
      </div>
    </MainLayout>
  )
}
