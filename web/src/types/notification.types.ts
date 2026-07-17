import { PaginatedData } from './general.types'

export type NotificationStatus = 'Read' | 'Unread'
export type NotificationTargetType = 'System' | 'Video' | 'Playlist' | 'Room' | string

export type Notification = {
  notificationId: string
  userId: string
  title: string
  message: string
  targetType: NotificationTargetType
  targetIdentifier: string
  createdAt: string
  notificationStatus: NotificationStatus
}

export type GetNotificationsResponse = PaginatedData<Notification>