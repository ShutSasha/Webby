import { useInfiniteQuery } from '@tanstack/react-query'

import { searchUserPlaylists } from '@/lib/actions/playlist.actions'

const PAGE_SIZE = 20

export const useSearchUserPlaylistsQuery = (userId: string, searchQuery: string, videoId?: string) => {
  return useInfiniteQuery({
    queryKey: ['search-user-playlists', userId, searchQuery, videoId],

    queryFn: async ({ pageParam = 1 }) => {
      return await searchUserPlaylists(userId, searchQuery, pageParam, PAGE_SIZE, videoId)
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
