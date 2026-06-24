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
    <div className="w-full flex justify-center mt-6">
      <form
        onSubmit={handleSubmit(onSubmit)}
        className="w-full max-w-[800px] border border-border rounded-2xl p-6 sm:p-8 flex flex-col items-center"
      >
        <div className="flex flex-col mb-8">
          <h2 className="text-[11px] font-bold text-foreground-muted uppercase tracking-wider mb-2">
            Password Management
          </h2>
          <p className="text-sm text-foreground-subtle max-w-xl">
            {` Update your password to keep your account secure. If you use a social login (like Google) and haven't set a
            password yet, or if you simply forgot it, you can reset it.`}
          </p>
        </div>

        <div className="flex flex-col gap-6 w-full max-w-md">
          <div className="flex flex-col gap-1.5">
            <div className="flex justify-between items-center">
              <label htmlFor="currentPassword" className="text-sm font-medium text-foreground-secondary">
                Current password
              </label>
              <Link
                href="/forgot-password"
                className="text-xs text-foreground-muted hover:text-foreground transition-colors"
              >
                Forgot your password?
              </Link>
            </div>
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
    </div>
  )
}
