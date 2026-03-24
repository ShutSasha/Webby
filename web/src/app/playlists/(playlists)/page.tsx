import { Suspense } from 'react'

import PlaylistSearchContainer from '@/ui/components/modules/Playlists/PlaylistSearchContainer'

type Props = {
  searchParams: Promise<{ page?: string; query?: string }>
}

export default async function Page({ searchParams }: Props) {
  const { page, query } = await searchParams
  const safeQuery = query || ''

  let currentPage = parseInt(page ?? '1', 10)

  if (Number.isNaN(currentPage) || currentPage < 1) {
    currentPage = 1
  }

  return (
    <Suspense key={`${currentPage}+${safeQuery}`}>
      <PlaylistSearchContainer page={currentPage} query={safeQuery} />
    </Suspense>
  )
}
