import PlaylistItem from '@/ui/components/modules/Playlists/PlaylistItem'
import GridCardsContainer from '@/ui/components/shared/GridCardsContainer'

export default function Page() {
  return (
    <GridCardsContainer>
      {Array.from({ length: 10 }).map((_, index) => (
        <PlaylistItem
          key={index}
          id={'98b84aad-372d-4c06-861a-fd361d2cbcb4'}
          src="https://i.ibb.co/V7LDjpM/anime-style-clouds.jpg"
          name="name"
          creator="creator"
          videoCount={1}
        />
      ))}
    </GridCardsContainer>
  )
}
