'use client'

const mockUsers = [
  { id: '1', username: 'cdidk', role: 'ADMIN', status: 'Active', registered: '2025-10-12', avatar: 'J' },
  { id: '2', username: 'anime_lover', role: 'USER', status: 'Banned', registered: '2026-01-05', avatar: 'A' },
  { id: '3', username: 'streamer_pro', role: 'MODERATOR', status: 'Active', registered: '2026-03-20', avatar: 'S' },
  { id: '4', username: 'toxic_guy', role: 'USER', status: 'Active', registered: '2026-05-11', avatar: 'T' },
]

export default function AdminUsersPage() {
  return (
    <div className="flex flex-col gap-6 animate-in fade-in duration-500">
      {/* Search Bar */}
      <div className="flex items-center gap-4">
        <div className="flex-1 bg-[#0A0A0A] border border-neutral-800/60 rounded-xl px-4 py-3 flex items-center gap-3">
          <svg className="size-5 text-neutral-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
            />
          </svg>
          <input
            type="text"
            placeholder="Search users by username..."
            className="bg-transparent border-none outline-none w-full text-neutral-200 placeholder:text-neutral-600"
          />
        </div>
      </div>

      {/* Users Table */}
      <div className="bg-[#0A0A0A] border border-neutral-800/60 rounded-2xl overflow-hidden">
        <table className="w-full text-left border-collapse">
          <thead>
            <tr className="border-b border-neutral-800/60 bg-neutral-900/50 text-neutral-500 text-sm">
              <th className="py-4 px-6 font-medium">User</th>
              <th className="py-4 px-6 font-medium">Role</th>
              <th className="py-4 px-6 font-medium">Status</th>
              <th className="py-4 px-6 font-medium">Registered</th>
              <th className="py-4 px-6 font-medium text-right">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-neutral-800/60">
            {mockUsers.map(user => (
              <tr key={user.id} className="hover:bg-neutral-800/20 transition-colors">
                <td className="py-4 px-6">
                  <div className="flex items-center gap-3">
                    <div
                      className="size-10 rounded-full bg-emerald-500/20 text-emerald-500 flex items-center
                        justify-center font-bold"
                    >
                      {user.avatar}
                    </div>
                    <span className="font-medium text-neutral-200">{user.username}</span>
                  </div>
                </td>
                <td className="py-4 px-6">
                  <span
                    className={`px-2.5 py-1 rounded-md text-xs font-bold ${
                      user.role === 'ADMIN'
                        ? 'bg-purple-500/10 text-purple-400 border border-purple-500/20'
                        : user.role === 'MODERATOR'
                          ? 'bg-blue-500/10 text-blue-400 border border-blue-500/20'
                          : 'bg-neutral-800 text-neutral-400'
                    }`}
                  >
                    {user.role}
                  </span>
                </td>
                <td className="py-4 px-6">
                  <div className="flex items-center gap-2">
                    <div
                      className={`size-2 rounded-full ${user.status === 'Active' ? 'bg-emerald-500' : 'bg-red-500'}`}
                    />
                    <span className={user.status === 'Active' ? 'text-neutral-300' : 'text-red-400'}>
                      {user.status}
                    </span>
                  </div>
                </td>
                <td className="py-4 px-6 text-neutral-500 text-sm">{user.registered}</td>
                <td className="py-4 px-6 flex justify-end gap-2">
                  <button
                    className="px-3 py-1.5 text-xs font-semibold bg-neutral-800 hover:bg-neutral-700 text-neutral-300
                      rounded-lg transition-colors"
                  >
                    Edit Role
                  </button>
                  {user.status === 'Active' ? (
                    <button
                      className="px-3 py-1.5 text-xs font-semibold bg-red-500/10 hover:bg-red-500/20 text-red-500 border
                        border-red-500/20 rounded-lg transition-colors"
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
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}
