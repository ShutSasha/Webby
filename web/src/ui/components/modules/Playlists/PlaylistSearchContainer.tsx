import { searchPlaylists } from '@/app/api/playlists'

import PlaylistItem from './PlaylistItem'
import GridCardsContainer from '../../shared/GridCardsContainer'

type Props = {
  page: number
  query: string
}

const PAGE_SIZE = 10

export default async function PlaylistSearchContainer({ page, query }: Props) {
  const searchResponse = await searchPlaylists(query, page, PAGE_SIZE)

  if (!searchResponse.data?.items?.length) {
    return (
      <div className="flex flex-1 items-center justify-center">
        <p className="text-neutral-500 text-center">
          Playlists by query <span className="text-neutral-300">{`'${query}'`}</span> not found
        </p>
      </div>
    )
  }

  const { items: playlists } = searchResponse.data

  return (
    <GridCardsContainer>
      {playlists.map(playlist => (
        <PlaylistItem
          key={playlist.playlistId}
          id={playlist.playlistId}
          src={playlist.playlistCover}
          name={playlist.name}
          creator="NOT FOUND"
          videoCount={playlist.countOfVideos}
        />
      ))}
    </GridCardsContainer>
  )
}
