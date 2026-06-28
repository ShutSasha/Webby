'use client'

import { useMemo } from 'react'

import { useSearchAdminUsersQuery } from '@/lib/hooks/api/admin/useSearchAdminUsers'
import { useInfiniteScroll } from '@/lib/hooks/useInfiniteScroll'

import TableSkeleton from './AdminUsersTableSkeleton'
import UserTableRow from './UserTableRow'

type Props = {
  searchQuery: string
}

export default function AdminUsersTable({ searchQuery }: Props) {
  const { data, isLoading, isFetchingNextPage, hasNextPage, fetchNextPage } = useSearchAdminUsersQuery(searchQuery, 20)

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
              return <UserTableRow key={user.userId} user={user} lastElementRef={isLast ? lastElementRef : null} />
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
