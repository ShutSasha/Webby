'use client'

import { startTransition, useActionState, useEffect, useState } from 'react'

import Link from 'next/link'
import { useForm } from 'react-hook-form'

import PasswordIcon from '@/assets/auth/ic_password.svg'
import { changePassworddAction } from '@/lib/actions/auth.actions'
import { useToastStore } from '@/stores/toast-store'
import AuthInput from '@/ui/components/modules/Auth/AuthInput'
import Button from '@/ui/components/shared/Button'

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
    <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col justify-between items-center my-4">
      <div className="flex flex-col gap-3 max-w-md w-full">
        <div className="flex flex-col gap-1">
          <label htmlFor="currentPassword" className="text-sm">
            Current password
          </label>
          <AuthInput
            {...register('currentPassword')}
            Icon={PasswordIcon}
            showPassword={showPassword}
            togglePassword={togglePassword}
            type="password"
            placeholder="Current password"
          />
        </div>
        <div className="flex flex-col gap-1">
          <label htmlFor="newPassword" className="text-sm">
            New password
          </label>
          <AuthInput
            {...register('newPassword')}
            Icon={PasswordIcon}
            showPassword={showPassword}
            togglePassword={togglePassword}
            type="password"
            placeholder="New password"
          />
        </div>
        <div className="flex flex-col gap-1">
          <label htmlFor="confirmNewPassword" className="text-sm">
            Confirm new password
          </label>
          <AuthInput
            {...register('confirmPassword')}
            Icon={PasswordIcon}
            showPassword={showPassword}
            togglePassword={togglePassword}
            type="password"
            placeholder="Confirm new password"
          />
        </div>
        <Link href={'/forgot-password'} className="text-right text-emerald-500 text-sm hover:underline cursor-pointer">
          Forgot your password?
        </Link>
        <Button
          disabled={isPending}
          type="submit"
          viewType="confirm"
          className={`text-[16px] leading-[22px] font-semibold w-fit mx-auto ${
            isPending ? 'bg-surface-tertiary hover:bg-surface-tertiary cursor-not-allowed' : 'cursor-pointer'
          }`}
          paddingClasses="px-5 py-2"
        >
          {isPending ? 'Saving...' : 'Save'}
        </Button>
        <div className={`${state?.errors ? 'block' : 'hidden'}`}>
          {state?.errors &&
            Object.entries(state.errors as Record<string, string>).map(([field, message]) => (
              <p key={field} className="text-red-500 text-sm">
                {message}
              </p>
            ))}
        </div>
      </div>
    </form>
  )
}
