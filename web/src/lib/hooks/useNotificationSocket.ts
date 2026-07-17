import { useEffect, useRef } from 'react'

import { HubConnection, HubConnectionBuilder, LogLevel, HttpTransportType } from '@microsoft/signalr'
import { useQueryClient } from '@tanstack/react-query'
import { signOut, useSession } from 'next-auth/react'

import { useNotificationPopupStore } from '@/stores/notification-popup.store'
import { Notification } from '@/types/notification.types'

import { getOneTimeTicket } from '../actions/notification.actions'

export const useNotificationSocket = () => {
  const { status } = useSession()
  const addPopup = useNotificationPopupStore(state => state.addPopup)
  const connectionRef = useRef<HubConnection | null>(null)
  const queryClient = useQueryClient()
  const { data: session } = useSession()

  useEffect(() => {
    if (status !== 'authenticated') return

    let isMounted = true

    const connectToHub = async (isRetry = false) => {
      try {
        const ticketResponse = await getOneTimeTicket()

        if (!ticketResponse.success || !ticketResponse.data) {
          console.error('SignalR Ticket Error:', ticketResponse.message)
          return
        }

        const ticket = ticketResponse.data

        if (!isMounted) return

        const baseUrl = process.env.NEXT_PUBLIC_SIGNALR_URL || 'http://localhost:5000/hubs/notifications'

        const hubUrl = `${baseUrl}?ticket=${ticket}`

        const connection = new HubConnectionBuilder()
          .withUrl(hubUrl, {
            skipNegotiation: true,
            transport: HttpTransportType.WebSockets,
          })
          .withAutomaticReconnect()
          .configureLogging(LogLevel.Information)
          .build()

        connectionRef.current = connection

        await connection.start()

        connection.on('ReceiveNotification', (notification: Notification) => {
          addPopup(notification)

          queryClient.invalidateQueries({ queryKey: ['user-notifications'] })
          queryClient.invalidateQueries({ queryKey: ['unread-notifications'] })
        })

        connection.on('AccountSuspended', async (userId: string) => {
          if (userId === session?.user?.id) {
            await signOut({ redirectTo: '/login' })
          }
        })

        connection.on('UpdateUnreadNotificationsCount', (count: number) => {
          queryClient.setQueryData(['unread-notifications-count'], count)
        })

        connection.on('AuthError', async (errorMsg: string) => {
          console.error('SignalR AuthError:', errorMsg)

          await connection.stop()

          if (!isRetry && isMounted) {
            connectToHub(true)
          }
        })
      } catch (err) {
        console.error('SignalR Connection Error: ', err)
      }
    }

    connectToHub()

    return () => {
      isMounted = false
      if (connectionRef.current) {
        connectionRef.current.stop()
      }
    }
  }, [status, addPopup, queryClient, session?.user?.id])
}
