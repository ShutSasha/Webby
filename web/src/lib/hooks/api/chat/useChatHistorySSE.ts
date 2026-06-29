import { useEffect } from 'react'

import { fetchEventSource } from '@microsoft/fetch-event-source'
import { useQueryClient } from '@tanstack/react-query'

import { serverLog } from '@/lib/utils/general.utils'

export const useChatHistorySSE = (jwtToken?: string | null) => {
  const queryClient = useQueryClient()

  useEffect(() => {
    if (!jwtToken) return

    const controller = new AbortController()

    const connectToSSE = async () => {
      try {
        await fetchEventSource(`${process.env.NEXT_PUBLIC_SSE_URL}/sse`, {
          method: 'GET',
          headers: {
            Authorization: `Bearer ${jwtToken}`,
            Accept: 'text/event-stream',
          },
          signal: controller.signal,

          onmessage(event) {
            if (event.event === 'UPDATE_CHAT_HISTORY') {
              queryClient.invalidateQueries({ queryKey: ['chat-history'] })
            }
          },
          async onopen(response) {
            if (response.ok) return

            throw new Error(`SSE connection failed with status: ${response.status}`)
          },
          onerror(error) {
            serverLog('SSE_CONNECTION_ERROR', error, false)

            if (error instanceof Error && error.message.match(/status: (401|403|404)/)) {
              throw error
            }

            return 5000
          },
        })
      } catch (error) {
        serverLog('SSE_FATAL_ERROR', error, false)
      }
    }

    connectToSSE()

    return () => {
      controller.abort()
    }
  }, [queryClient, jwtToken])
}
