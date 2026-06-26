'use client'

import { useMemo } from 'react'

import Image from 'next/image'

import { useSearchAdminUsersQuery } from '@/lib/hooks/api/admin/useSearchAdminUsers'
import { useInfiniteScroll } from '@/lib/hooks/useInfiniteScroll'
import { formatDate } from '@/lib/utils/date.utils'

import TableSkeleton from './AdminUsersTableSkeleton'

type Props = {
  searchQuery: string
}

export default function AdminUsersTable({ searchQuery }: Props) {
  const { data, isLoading, isFetchingNextPage, hasNextPage, fetchNextPage } = useSearchAdminUsersQuery(searchQuery, 10)

  const users = useMemo(() => {
    return data?.pages.flatMap(page => page?.items || []) || []
  }, [data])

  const lastElementRef = useInfiniteScroll({
    isLoading,
    isFetchingNextPage,
    hasNextPage,
    fetchNextPage,
  })

  if (isLoading && users.length === 0) {
    return <TableSkeleton />
  }

  if (!isLoading && users.length === 0) {
    return (
      <div
        className="flex flex-col items-center justify-center py-20 bg-surface border border-border rounded-2xl
          text-foreground-faint"
      >
        <p>No users found matching &quot;{searchQuery}&quot;</p>
      </div>
    )
  }

  return (
    <div className="bg-surface border border-border rounded-2xl overflow-hidden flex flex-col">
      <div className="overflow-x-auto custom-scrollbar">
        <table className="w-full text-left border-collapse min-w-[800px]">
          <thead>
            <tr className="border-b border-border bg-surface/50 text-foreground-faint text-sm">
              <th className="py-4 px-6 font-medium">User</th>
              <th className="py-4 px-6 font-medium">Role</th>
              <th className="py-4 px-6 font-medium">Status</th>
              <th className="py-4 px-6 font-medium">Registered</th>
              <th className="py-4 px-6 font-medium text-right">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-background/60">
            {users.map((user, index) => {
              const isLast = index === users.length - 1
              const role = user.role.toUpperCase()
              const status = user.isBanned ? 'Banned' : 'Active'

              return (
                <tr
                  key={user.userId}
                  ref={isLast ? lastElementRef : null}
                  className="hover:bg-background/20 transition-colors"
                >
                  <td className="py-4 px-6">
                    <div className="flex items-center gap-3">
                      <div
                        className="size-10 rounded-full bg-emerald-500/20 text-emerald-500 flex items-center
                          justify-center font-bold overflow-hidden shrink-0"
                      >
                        {user.avatarUrl && user.avatarUrl.startsWith('http') ? (
                          <Image
                            src={user.avatarUrl}
                            alt={user.username}
                            width={100}
                            height={100}
                            className="size-full object-cover"
                          />
                        ) : (
                          user.username?.charAt(0).toUpperCase() || 'U'
                        )}
                      </div>
                      <div className="flex flex-col">
                        <span className="font-medium text-foreground-tertiary">{user.username}</span>
                        <span className="text-xs text-foreground-faint">{user.email}</span>
                      </div>
                    </div>
                  </td>
                  <td className="py-4 px-6">
                    <span
                      className={`px-2.5 py-1 rounded-md text-xs font-bold ${
                        role === 'ADMIN'
                          ? 'bg-purple-500/10 text-purple-400 border border-purple-500/20'
                          : role === 'MODERATOR'
                            ? 'bg-blue-500/10 text-blue-400 border border-blue-500/20'
                            : 'bg-background text-foreground-muted border border-border/50'
                      }`}
                    >
                      {role}
                    </span>
                  </td>
                  <td className="py-4 px-6">
                    <div className="flex items-center gap-2">
                      <div className={`size-2 rounded-full ${status === 'Active' ? 'bg-emerald-500' : 'bg-red-500'}`} />
                      <span className={status === 'Active' ? 'text-foreground-subtle' : 'text-red-400'}>{status}</span>
                    </div>
                  </td>
                  <td className="py-4 px-6 text-foreground-faint text-sm">{formatDate(user.createdAt)}</td>
                  <td className="py-4 px-6 flex justify-end gap-2">
                    <button
                      className="px-3 py-1.5 text-xs font-semibold bg-background hover:bg-surface-tertiary
                        text-foreground-subtle rounded-lg transition-colors"
                    >
                      Edit Role
                    </button>
                    {status === 'Active' ? (
                      <button
                        className="px-3 py-1.5 text-xs font-semibold bg-red-500/10 hover:bg-red-500/20 text-red-500
                          border border-red-500/20 rounded-lg transition-colors"
                      >
                        Ban
                      </button>
                    ) : (
                      <button
                        className="px-3 py-1.5 text-xs font-semibold bg-emerald-500/10 hover:bg-emerald-500/20
                          text-emerald-500 border border-emerald-500/20 rounded-lg transition-colors"
                      >
                        Unban
                      </button>
                    )}
                  </td>
                </tr>
              )
            })}
          </tbody>
        </table>
      </div>

      {isFetchingNextPage && (
        <div className="flex justify-center py-4 border-t border-background/60">
          <div className="size-6 border-4 border-emerald-500/20 border-t-emerald-500 rounded-full animate-spin" />
        </div>
      )}
    </div>
  )
}
