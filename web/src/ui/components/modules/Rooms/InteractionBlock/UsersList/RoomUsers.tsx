'use client'

import { useParams } from 'next/navigation'

import UserList from './UserList'

export default function RoomUsers() {
  const params = useParams()
  const roomId = params?.id as string

  if (!roomId) return null

  return <UserList roomId={roomId} />
}
