'use client'

import { startTransition, useActionState, useState } from 'react'

import { useForm } from 'react-hook-form'

import { authenticate } from '@/app/api/auth'
import PasswordIcon from '@/assets/auth/ic_password.svg'
import AuthInput from '@/ui/components/AuthInput'
import Input from '@/ui/components/Input'

const initialState = {
  success: false,
  errors: null,
}

export default function ResetPasswordForm() {
  const [showPassword, setShowPassword] = useState<boolean>(false)
  const togglePassword = () => setShowPassword(prev => !prev)
  const [state, formAction, isPending] = useActionState(authenticate, initialState)

  const { register, handleSubmit } = useForm({
    defaultValues: {
      currentPassword: '',
      newPassword: '',
      confirmNewPassword: '',
    },
  })

  const onSubmit = (data: any) => {
    const formData = new FormData()
    Object.entries(data).forEach(([key, value]) => formData.append(key, value as string))

    startTransition(() => {
      formAction(formData)
    })
  }

  return (
    <div className="flex flex-col justify-between items-center my-4">
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
            {...register('confirmNewPassword')}
            Icon={PasswordIcon}
            showPassword={showPassword}
            togglePassword={togglePassword}
            type="password"
            placeholder="Confirm new password"
          />
        </div>
      </div>
    </div>
  )
}
