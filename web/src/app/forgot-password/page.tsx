'use client'

import { useEffect } from 'react'

import { useCommonStore } from '@/stores/common.store'
import CodeStep from '@/ui/components/modules/Auth/ResetPassword/CodeStep'
import EmailStep from '@/ui/components/modules/Auth/ResetPassword/EmailStep'
import NewPasswordStep from '@/ui/components/modules/Auth/ResetPassword/NewPasswordStep'

export default function ForgotPasswordPage() {
  const step = useCommonStore(state => state.forgotPasswordStep)
  const reset = useCommonStore(state => state.resetForgotPassword)

  useEffect(() => {
    reset()
  }, [reset])

  return (
    <div className="flex flex-1 items-center justify-center py-10">
      <div className="flex w-full bg-surface max-w-[700px] rounded-[20px] p-10 box-border flex-col items-center gap-6">
        {step === 1 && <EmailStep />}
        {step === 2 && <CodeStep />}
        {step === 3 && <NewPasswordStep />}
      </div>
    </div>
  )
}
