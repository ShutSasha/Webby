import { useInfiniteQuery } from '@tanstack/react-query'

import { searchPlaylists } from '@/lib/actions/playlist.actions'
import { unwrapServerAction } from '@/lib/utils/general.utils'

const PAGE_SIZE = 20

export const useSearchPlaylistsQuery = (searchQuery: string) => {
  return useInfiniteQuery({
    queryKey: ['search-playlists', searchQuery],

    queryFn: async ({ pageParam = 1 }) => {
      return unwrapServerAction(await searchPlaylists(searchQuery, pageParam, PAGE_SIZE))
    },
    initialPageParam: 1,
    getNextPageParam: (lastPage, allPages) => {
      const items = lastPage?.items || []

      if (items.length === PAGE_SIZE) {
        return allPages.length + 1
      }

      return undefined
    },

    staleTime: 5 * 60 * 1000,
  })
}
