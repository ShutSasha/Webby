'use client'

import { useEffect, useRef, useState } from 'react'

import { Role } from '@/types/user.types'

type Props = {
  currentRole: Role
  onRoleChange: (newRole: Role) => void
  disabled?: boolean
}

export default function RoleSelectMenu({ currentRole, onRoleChange, disabled }: Props) {
  const [isOpen, setIsOpen] = useState(false)
  const menuRef = useRef<HTMLDivElement>(null)

  const roles: Role[] = ['User', 'Moderator']

  useEffect(() => {
    if (!isOpen) return

    const handleClickOutside = (e: MouseEvent) => {
      if (menuRef.current && !menuRef.current.contains(e.target as Node)) {
        setIsOpen(false)
      }
    }

    document.addEventListener('mousedown', handleClickOutside)
    return () => document.removeEventListener('mousedown', handleClickOutside)
  }, [isOpen])

  return (
    <div className="relative inline-block" ref={menuRef}>
      <button
        onClick={() => setIsOpen(prev => !prev)}
        disabled={disabled}
        className={`px-3 py-1.5 min-w-[90px] text-xs font-semibold rounded-lg transition-colors border
          focus:outline-none flex items-center justify-center ${
            isOpen
              ? 'bg-surface border-emerald-500 text-foreground-secondary shadow-[0_0_0_1px_var(--color-emerald-500)]'
              : 'bg-background hover:bg-surface-tertiary text-foreground-subtle border-transparent'
          } ${disabled ? 'opacity-50 cursor-not-allowed' : 'cursor-pointer'}`}
      >
        {currentRole}
      </button>

      {isOpen && (
        <div
          className="absolute top-full left-1/2 -translate-x-1/2 mt-1.5 w-full min-w-[110px] z-50 bg-surface border
            border-border shadow-lg shadow-black/10 py-1 rounded-lg animate-in fade-in zoom-in-95 duration-150
            overflow-hidden"
        >
          {roles.map(role => (
            <button
              key={role}
              onClick={() => {
                if (role !== currentRole) {
                  onRoleChange(role)
                }
                setIsOpen(false)
              }}
              className={`w-full text-center px-4 py-2 text-xs transition-colors cursor-pointer ${
                role === currentRole
                  ? 'bg-emerald-500/10 text-emerald-500 font-bold'
                  : 'text-foreground-secondary hover:bg-surface-tertiary/50'
              }`}
            >
              {role}
            </button>
          ))}
        </div>
      )}
    </div>
  )
}
