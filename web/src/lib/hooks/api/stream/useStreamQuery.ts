import { useQuery, UseQueryOptions } from '@tanstack/react-query'

import { getStreamInfo } from '@/lib/actions/stream.actions'
import { Stream } from '@/types/stream.types'

export const useStreamQuery = (
  streamerId: string | undefined,
  options?: Omit<UseQueryOptions<Stream>, 'queryKey' | 'queryFn'>,
) => {
  return useQuery({
    queryKey: ['stream', streamerId],
    queryFn: async () => {
      if (!streamerId) throw new Error('Streamer ID is required')
      const res = await getStreamInfo(streamerId)
      if (!res.success || !res.data) throw new Error(res.message)
      return res.data
    },
    enabled: !!streamerId && options?.enabled,
    staleTime: 60 * 1000,
  })
}
