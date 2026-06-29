import { RefCallback } from 'react'

import Image from 'next/image'

import { AdminUserRecord } from '@/lib/actions/admin.actions'
import { useBanUserMutation } from '@/lib/hooks/api/admin/useBanUser'
import { useChangeUserRoleMutation } from '@/lib/hooks/api/admin/useChangeUserRole'
import { useUnbanUserMutation } from '@/lib/hooks/api/admin/useUnbanUser'
import { formatDate } from '@/lib/utils/date.utils'

import RoleSelectMenu from './RoleSelectMenu'

type UserTableRowProps = {
  user: AdminUserRecord
  lastElementRef: RefCallback<HTMLTableRowElement> | null
}

export default function UserTableRow({ user, lastElementRef }: UserTableRowProps) {
  const { mutate: banUser, isPending: isBanning } = useBanUserMutation()
  const { mutate: unbanUser, isPending: isUnbanning } = useUnbanUserMutation()
  const { mutate: changeRole, isPending: isChangingRole } = useChangeUserRoleMutation()

  const isPending = isBanning || isUnbanning || isChangingRole
  const role = user.role.toUpperCase()
  const status = user.isBanned ? 'Banned' : 'Active'

  return (
    <tr
      ref={lastElementRef}
      className={`hover:bg-background/20 transition-colors ${isPending ? 'opacity-60 pointer-events-none' : ''}`}
    >
      <td className="py-4 px-6">
        <div className="flex items-center gap-3">
          <div
            className="size-10 rounded-full bg-emerald-500/20 text-emerald-500 flex items-center justify-center
              font-bold overflow-hidden shrink-0"
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
        <RoleSelectMenu
          currentRole={user.role}
          onRoleChange={newRole => changeRole({ userId: user.userId, newRole })}
          disabled={isPending}
        />

        {status === 'Active' ? (
          <button
            onClick={() => banUser(user.userId)}
            disabled={isPending}
            className="w-[68px] flex items-center justify-center px-3 py-1.5 text-xs font-semibold bg-red-500/10
              hover:bg-red-500/20 text-red-500 border border-red-500/20 rounded-lg transition-colors
              disabled:opacity-50"
          >
            {isBanning ? (
              <div className="size-3.5 border-2 border-red-500/30 border-t-red-500 rounded-full animate-spin" />
            ) : (
              'Ban'
            )}
          </button>
        ) : (
          <button
            onClick={() => unbanUser(user.userId)}
            disabled={isPending}
            className="w-[68px] flex items-center justify-center px-3 py-1.5 text-xs font-semibold bg-emerald-500/10
              hover:bg-emerald-500/20 text-emerald-500 border border-emerald-500/20 rounded-lg transition-colors
              disabled:opacity-50"
          >
            {isUnbanning ? (
              <div className="size-3.5 border-2 border-emerald-500/30 border-t-emerald-500 rounded-full animate-spin" />
            ) : (
              'Unban'
            )}
          </button>
        )}
      </td>
    </tr>
  )
}
