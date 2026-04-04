import { useInfiniteQuery } from '@tanstack/react-query'

import { searchUserPlaylists } from '@/lib/actions/playlist.actions'

const PAGE_SIZE = 20

export const useSearchUserPlaylistsQuery = (
  userId: string | undefined,
  searchQuery: string,
  videoId?: string,
  enabled: boolean = true,
) => {
  return useInfiniteQuery({
    queryKey: ['search-user-playlists', userId, searchQuery, videoId],

    queryFn: async ({ pageParam = 1 }) => {
      if (!userId) throw new Error('User ID is required')
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

    enabled: enabled && !!userId,
    staleTime: 5 * 60 * 1000,
  })
}
