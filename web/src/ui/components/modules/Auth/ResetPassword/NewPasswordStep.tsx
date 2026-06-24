'use client'

import { useActionState, startTransition, useState, useEffect } from 'react'

import Link from 'next/link'
import { useForm } from 'react-hook-form'

import PasswordIcon from '@/assets/auth/ic_password.svg'
import { resetPasswordAction } from '@/lib/actions/auth.actions'
import { useCommonStore } from '@/stores/common.store'
import { useToastStore } from '@/stores/toast-store'
import AuthInput from '@/ui/components/modules/Auth/AuthInput'
import Button from '@/ui/components/shared/Button'

export default function NewPasswordStep() {
  const [showPassword, setShowPassword] = useState<boolean>(false)
  const togglePassword = () => setShowPassword(prev => !prev)

  const forgotPasswordEmail = useCommonStore(state => state.forgotPasswordEmail)
  const addToast = useToastStore(state => state.addToast)

  const [state, formAction, isPending] = useActionState(resetPasswordAction, {
    success: false,
    errors: null,
    timestamp: Date.now(),
  })

  const { register, handleSubmit, reset } = useForm({
    defaultValues: {
      newPassword: '',
      confirmPassword: '',
    },
  })

  useEffect(() => {
    if (state?.success) {
      addToast('Password has been successfully reset!', 'success')
      reset()
    }
  }, [state, addToast, reset])

  const onSubmit = (data: any) => {
    const formData = new FormData()

    formData.append('email', forgotPasswordEmail)
    formData.append('newPassword', data.newPassword)
    formData.append('confirmPassword', data.confirmPassword)

    startTransition(() => {
      formAction(formData)
    })
  }

  if (state.success) {
    return (
      <div className="flex flex-col items-center text-center space-y-4 animate-in fade-in zoom-in duration-300">
        <p className="text-emerald-500 font-medium text-lg">Password has been reset!</p>
        <p className="text-foreground-muted text-sm">You can now use your new password to log in.</p>
        <Button viewType="confirm" className="px-8">
          <Link href={'/login'}>Go to Login</Link>
        </Button>
      </div>
    )
  }

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="w-full space-y-6">
      <div className="text-center space-y-2">
        <h2 className="text-foreground-strong text-2xl font-bold">Create new password</h2>
        <p className="text-foreground-muted text-sm">Enter a new strong password for your account</p>
      </div>

      <div className="space-y-4">
        <div className="flex flex-col gap-1">
          <label className="text-sm text-foreground-muted ml-1">New password</label>
          <AuthInput
            {...register('newPassword', { required: true })}
            type="password"
            Icon={PasswordIcon}
            placeholder="Enter new password"
            showPassword={showPassword}
            togglePassword={togglePassword}
          />
        </div>

        <div className="flex flex-col gap-1">
          <label className="text-sm text-foreground-muted ml-1">Confirm password</label>
          <AuthInput
            {...register('confirmPassword', { required: true })}
            type="password"
            Icon={PasswordIcon}
            placeholder="Repeat new password"
            showPassword={showPassword}
            togglePassword={togglePassword}
          />
        </div>
      </div>

      <Button disabled={isPending} type="submit" viewType="confirm" className="w-full py-3 rounded-xl font-semibold">
        {isPending ? 'Saving...' : 'Reset Password'}
      </Button>

      {state.errors && (
        <div className="space-y-1 bg-red-500/10 p-3 rounded-lg border border-red-500/20">
          {Object.entries(state.errors).map(([field, message]) => (
            <p key={field} className="text-red-500 text-xs text-center">
              {message as string}
            </p>
          ))}
        </div>
      )}
    </form>
  )
}
