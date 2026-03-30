import { useInfiniteQuery } from '@tanstack/react-query'

import { searchPlaylists } from '@/lib/actions/playlist.actions'

const PAGE_SIZE = 20

export const useSearchPlaylistsQuery = (searchQuery: string) => {
  return useInfiniteQuery({
    queryKey: ['search-playlists', searchQuery],

    queryFn: async ({ pageParam = 1 }) => {
      return await searchPlaylists(searchQuery, pageParam, PAGE_SIZE)
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
