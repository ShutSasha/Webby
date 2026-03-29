'use client'

import { useState } from 'react'

import { useForm } from 'react-hook-form'

import { resendVerifyCode } from '@/app/api/auth'
import MailIcon from '@/assets/auth/ic_mail.svg'
import { parseAxiosError } from '@/lib/utils/utils'
import { useCommonStore } from '@/stores/common.store'
import AuthInput from '@/ui/components/modules/Auth/AuthInput'
import Button from '@/ui/components/shared/Button'

export default function EmailStep() {
  const { setForgotPasswordStep, setForgotPasswordEmail } = useCommonStore()
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const { register, handleSubmit } = useForm({ defaultValues: { email: '' } })

  const onSubmit = async (data: { email: string }) => {
    setLoading(true)
    setError(null)

    try {
      await resendVerifyCode({ email: data.email })

      setForgotPasswordEmail(data.email)
      setForgotPasswordStep(2)
    } catch (err) {
      const serverErrors = parseAxiosError(err)
      setError(serverErrors.global || 'Failed to send code')
    } finally {
      setLoading(false)
    }
  }

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="w-full space-y-6">
      <div className="text-center space-y-2">
        <h2 className="text-white text-2xl font-bold">Forgot Password</h2>
        <p className="text-neutral-400 text-sm">Enter the email associated with your account</p>
      </div>

      <AuthInput
        {...register('email', { required: true })}
        Icon={MailIcon}
        placeholder="Enter your email"
        type="email"
      />

      {error && <p className="text-red-500 text-sm text-center">{error}</p>}

      <Button type="submit" viewType="confirm" className="w-full rounded-xl py-3 font-semibold" disabled={loading}>
        {loading ? 'Sending code...' : 'Next Step'}
      </Button>
    </form>
  )
}
