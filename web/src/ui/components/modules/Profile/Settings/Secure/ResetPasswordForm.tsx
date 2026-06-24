'use client'

import { startTransition, useActionState, useEffect, useState } from 'react'

import Link from 'next/link'
import { useForm } from 'react-hook-form'

import PasswordIcon from '@/assets/auth/ic_password.svg'
import { changePassworddAction } from '@/lib/actions/auth.actions'
import { cn } from '@/lib/utils/general.utils'
import { useToastStore } from '@/stores/toast-store'
import AuthInput from '@/ui/components/modules/Auth/AuthInput'

const initialState = {
  success: false,
  errors: null,
  timestamp: Date.now(),
}

export default function ResetPasswordForm() {
  const [showPassword, setShowPassword] = useState<boolean>(false)
  const togglePassword = () => setShowPassword(prev => !prev)
  const [state, formAction, isPending] = useActionState(changePassworddAction, initialState)
  const addToast = useToastStore(state => state.addToast)

  const { register, handleSubmit, reset } = useForm({
    defaultValues: {
      currentPassword: '',
      newPassword: '',
      confirmPassword: '',
    },
  })

  useEffect(() => {
    if (state?.success) {
      addToast('Password changed successfully!', 'success')
      reset()
    }
  }, [state, addToast, reset])

  const onSubmit = (data: any) => {
    const formData = new FormData()
    formData.append('currentPassword', data.currentPassword)
    formData.append('newPassword', data.newPassword)
    formData.append('confirmPassword', data.confirmPassword)

    startTransition(() => {
      formAction(formData)
    })
  }

  return (
    <div className="w-full flex flex-col items-center gap-6 mt-6">
      <form
        onSubmit={handleSubmit(onSubmit)}
        className="w-full max-w-[800px] border border-border rounded-2xl p-6 sm:p-8 flex flex-col"
      >
        <div className="flex flex-col mb-8">
          <h2 className="text-[11px] font-bold text-foreground-muted uppercase tracking-wider mb-2">Change Password</h2>
          <p className="text-sm text-foreground-subtle leading-relaxed max-w-2xl">
            Update your password to keep your account secure. Make sure to choose a strong, unique password.
          </p>
        </div>

        <div className="flex flex-col gap-6 w-full max-w-md mx-auto">
          <div className="flex flex-col gap-1.5">
            <label htmlFor="currentPassword" className="text-sm font-medium text-foreground-secondary">
              Current password
            </label>
            <AuthInput
              {...register('currentPassword')}
              Icon={PasswordIcon}
              showPassword={showPassword}
              togglePassword={togglePassword}
              type="password"
              placeholder="Enter your current password"
            />
          </div>

          <div className="flex flex-col gap-1.5">
            <label htmlFor="newPassword" className="text-sm font-medium text-foreground-secondary">
              New password
            </label>
            <AuthInput
              {...register('newPassword')}
              Icon={PasswordIcon}
              showPassword={showPassword}
              togglePassword={togglePassword}
              type="password"
              placeholder="Create a new password"
            />
          </div>

          <div className="flex flex-col gap-1.5">
            <label htmlFor="confirmNewPassword" className="text-sm font-medium text-foreground-secondary">
              Confirm new password
            </label>
            <AuthInput
              {...register('confirmPassword')}
              Icon={PasswordIcon}
              showPassword={showPassword}
              togglePassword={togglePassword}
              type="password"
              placeholder="Confirm your new password"
            />
          </div>

          {state?.errors && (
            <div className="flex flex-col gap-1 bg-red-500/10 border border-red-500/20 rounded-lg p-3 mt-2">
              {Object.entries(state.errors as Record<string, string>).map(([field, message]) => (
                <p key={field} className="text-red-500 text-sm font-medium">
                  {message}
                </p>
              ))}
            </div>
          )}
        </div>

        <div className="mt-10 flex justify-center w-full">
          <button
            disabled={isPending}
            type="submit"
            className={cn(
              'px-8 py-2.5 rounded-xl font-medium transition-all duration-300 min-w-40 border',
              isPending &&
                'bg-surface-secondary text-foreground-muted cursor-not-allowed opacity-70 border-border/50 shadow-none',
              !isPending &&
                `bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 hover:bg-emerald-500/20 border-emerald-500/30
                hover:border-emerald-500/50 shadow-sm`,
            )}
          >
            {isPending ? 'Saving...' : 'Save changes'}
          </button>
        </div>
      </form>

      <div
        className="w-full max-w-[800px] border border-border rounded-2xl p-6 sm:p-8 flex flex-col sm:flex-row
          justify-between items-start sm:items-center gap-6"
      >
        <div className="flex flex-col max-w-lg">
          <h2 className="text-[11px] font-bold text-foreground-muted uppercase tracking-wider mb-2">
            Recovery & Social Logins
          </h2>
          <p className="text-sm text-foreground-subtle leading-relaxed">
            {`If you use a social login (like Google) and haven't set a password yet, or if you simply forgot your current
            one, you can reset it.`}
          </p>
        </div>

        <Link
          href="/forgot-password"
          className="shrink-0 px-6 py-2.5 rounded-xl text-sm font-medium border border-border bg-transparent
            text-foreground-secondary hover:text-foreground-strong hover:bg-surface-secondary transition-colors
            duration-200"
        >
          Reset Password
        </Link>
      </div>
    </div>
  )
}
