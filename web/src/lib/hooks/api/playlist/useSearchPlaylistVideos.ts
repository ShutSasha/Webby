import { useInfiniteQuery } from '@tanstack/react-query'

import { getPlaylistVideos } from '@/lib/actions/playlist.actions'

const PAGE_SIZE = 20

export const useSearchPlaylistVideosQuery = (playlistId: string, searchQuery: string) => {
  return useInfiniteQuery({
    queryKey: ['search-playlist-videos', playlistId, searchQuery],

    queryFn: async ({ pageParam = 1 }) => {
      return await getPlaylistVideos(playlistId, searchQuery, pageParam, PAGE_SIZE)
    },
    initialPageParam: 1,
    getNextPageParam: (lastPage, allPages) => {
      const items = lastPage?.data?.items || []

      if (items.length === PAGE_SIZE) {
        return allPages.length + 1
      }

      return undefined
    },

    staleTime: 5 * 60 * 1000,
  })
}
