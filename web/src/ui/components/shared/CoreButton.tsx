'use client'

import { ButtonHTMLAttributes, ReactNode } from 'react'

import { cn } from '@/lib/utils/general.utils'

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  children: ReactNode
  variant?: 'primary' | 'secondary' | 'ghost' | 'danger'
  isLoading?: boolean
}

export default function CoreButton({
  children,
  variant = 'primary',
  isLoading = false,
  disabled,
  className,
  ...props
}: ButtonProps) {
  const baseStyles =
    'relative inline-flex items-center justify-center gap-2 font-semibold transition-all duration-300 ease-out rounded-xl text-sm overflow-hidden disabled:opacity-50 disabled:cursor-not-allowed'

  const variants = {
    primary: 'bg-emerald-500 hover:bg-emerald-600 text-foreground-inverse cursor-pointer',
    secondary: 'bg-background hover:bg-surface-tertiary/80 text-foreground-tertiary',
    danger: 'bg-red-500/10 hover:bg-red-500/20 text-red-500',
    ghost: 'bg-transparent hover:bg-background text-foreground-muted hover:text-foreground-secondary',
  }

  const sizeStyles = 'px-6 py-2.5'

  return (
    <button
      disabled={disabled || isLoading}
      className={cn(baseStyles, variants[variant], sizeStyles, className)}
      {...props}
    >
      {isLoading && (
        <div className="absolute inset-0 flex items-center justify-center bg-inherit">
          <div
            className={cn(
              'size-5 border-2 rounded-full animate-spin',
              variant === 'primary'
                ? 'border-neutral-900/20 border-t-neutral-900'
                : 'border-emerald-500/20 border-t-emerald-500',
            )}
          />
        </div>
      )}

      <span className={cn('flex items-center gap-2', isLoading && 'opacity-0')}>{children}</span>
    </button>
  )
}
