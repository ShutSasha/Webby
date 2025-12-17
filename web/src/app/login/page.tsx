import Image from 'next/image'

import ImageBackground from '@/assets/auth/bg-login.png'
import LoginForm from '@/ui/components/modules/Login/LoginForm'

export default function Page() {
  return (
    <div className="flex flex-1 items-center justify-center">
      <div className="flex w-full bg-neutral-900 max-w-[700px] rounded-[20px] p-5 gap-4 box-border">
        <LoginForm />

        <div className="shrink min-w-0 items-center hidden sm:block">
          <Image src={ImageBackground} alt="Sign up background" loading="eager" className="h-full w-75 rounded-lg" />
        </div>
      </div>
    </div>
  )
}
