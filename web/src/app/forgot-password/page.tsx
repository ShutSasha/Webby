'use client'

import { useEffect } from 'react'

import { useCommonStore } from '@/stores/common.store'
import CodeStep from '@/ui/components/common/ResetPassword/CodeStep'
import EmailStep from '@/ui/components/common/ResetPassword/EmailStep'
import NewPasswordStep from '@/ui/components/common/ResetPassword/NewPasswordStep'
import MainLayout from '@/ui/components/MainLayout'

export default function ForgotPasswordPage() {
  const step = useCommonStore(state => state.forgotPasswordStep)
  const reset = useCommonStore(state => state.resetForgotPassword)

  useEffect(() => {
    reset()
  }, [reset])

  return (
    <MainLayout>
      <div className="flex flex-1 items-center justify-center py-10">
        <div
          className="flex w-full bg-neutral-900 max-w-[700px] rounded-[20px] p-10 box-border flex-col items-center
            gap-6"
        >
          {step === 1 && <EmailStep />}
          {step === 2 && <CodeStep />}
          {step === 3 && <NewPasswordStep />}
        </div>
      </div>
    </MainLayout>
  )
}
