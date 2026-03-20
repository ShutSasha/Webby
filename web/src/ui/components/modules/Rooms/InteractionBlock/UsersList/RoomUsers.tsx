'use client'

import { useEffect, useState } from 'react'

import UserList from './UserList'

async function fetchUsers() {
  return Array.from({ length: 20 }, (_, i) => ({
    id: i,
    name: `guuuntersteam`,
    avatar: 'https://i.ibb.co/PGL4ymBS/thumb-1920-415519.jpg',
  }))
}

export default function RoomUsers() {
  const [users, setUsers] = useState<any[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const loadData = async () => {
      try {
        const data = await fetchUsers()
        setUsers(data)
      } catch (error) {
        console.error('Failed to fetch users:', error)
      } finally {
        setLoading(false)
      }
    }

    loadData()
  }, [])

  if (loading) return <p className="text-center text-neutral-500 py-4">Loading users...</p>

  return <UserList initialUsers={users} />
}
