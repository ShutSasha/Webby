import { useEffect, useRef } from 'react'

import { HubConnection, HubConnectionBuilder, LogLevel, HttpTransportType } from '@microsoft/signalr'
import { useQueryClient } from '@tanstack/react-query'
import { useSession } from 'next-auth/react'

import { useNotificationPopupStore } from '@/stores/notification-popup.store'
import { Notification } from '@/types/notification.types'

export const useNotificationSocket = () => {
  const { data: session } = useSession()
  const addPopup = useNotificationPopupStore(state => state.addPopup)
  const setUnreadCount = useNotificationPopupStore(state => state.setUnreadCount)
  const connectionRef = useRef<HubConnection | null>(null)
  const queryClient = useQueryClient()
  const token = session?.user.accessToken
  
  useEffect(() => {
    if (!token) return

    const baseUrl = process.env.NEXT_PUBLIC_SIGNALR_URL || 'http://localhost:5000/hubs/notifications'

    const hubUrl = `${baseUrl}?accessToken=${token}`

    const connection = new HubConnectionBuilder()
      .withUrl(hubUrl, {
        skipNegotiation: true,
        transport: HttpTransportType.WebSockets,
      })
      .withAutomaticReconnect()
      .configureLogging(LogLevel.Information)
      .build()

    connectionRef.current = connection

    connection
      .start()
      .then(() => {
        console.log('SignalR Connected Successfully.')

        connection.on('ReceiveNotification', (notification: Notification) => {
          addPopup(notification)

          queryClient.invalidateQueries({ queryKey: ['user-notifications'] })
          queryClient.invalidateQueries({ queryKey: ['unread-notifications'] })
        })

        connection.on('UpdateUnreadNotificationsCount', (count: number) => {
          setUnreadCount(count)
          queryClient.setQueryData(['unread-notifications-count'], count)
        })

        connection.on('AuthError', (errorMsg: string) => {
          console.error('SignalR AuthError:', errorMsg)
          connection.stop()
        })
      })
      .catch(err => console.error('SignalR Connection Error: ', err))

    return () => {
      connection.stop()
    }
  }, [token, addPopup, setUnreadCount, queryClient])
}
