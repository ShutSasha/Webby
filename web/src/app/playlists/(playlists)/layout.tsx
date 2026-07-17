import PageToggle from '@/ui/components/modules/Playlists/PageToggle'
import CreatePlaylistButton from '@/ui/components/modules/Rooms/InteractionBlock/Playlists/CreatePlaylistButton'
import MediaHeader from '@/ui/components/shared/MediaHeader'

type Props = Readonly<{
  children: React.ReactNode
}>

export default function Layout({ children }: Props) {
  return (
    <>
      <MediaHeader
        searchPlaceholder="Search a playlist"
        actionSlot={<CreatePlaylistButton />}
        pageToggle={<PageToggle />}
      />
      {children}
    </>
  )
}
