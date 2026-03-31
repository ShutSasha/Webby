import CreateRoomButton from '@/ui/components/modules/Rooms/CreateRoomButton'
import PageToggle from '@/ui/components/modules/Rooms/PageToggle'
import MediaHeader from '@/ui/components/shared/MediaHeader'

type Props = Readonly<{
  children: React.ReactNode
}>

export default function Layout({ children }: Props) {
  return (
    <>
      <MediaHeader
        searchPlaceholder="Search a playlist"
        actionSlot={<CreateRoomButton />}
        pageToggle={<PageToggle />}
      />
      {children}
    </>
  )
}
