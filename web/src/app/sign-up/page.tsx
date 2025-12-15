import Image from 'next/image'

import ImageBackground from '@/assets/auth/home-page-man-sign-up.png'
import SignUpForm from '@/ui/components/modules/SignUp/SignUpForm'

export default function Page() {
  return (
    <div className="flex flex-1 items-center justify-center">
      <div className="flex w-full bg-neutral-900 max-w-[700px] rounded-[20px] p-5 gap-4 box-border">
        <SignUpForm />

        <div className="shrink min-w-0 items-center hidden sm:block">
          <Image src={ImageBackground} alt="Sign up background" loading="eager" className="h-full w-110 rounded-lg" />
        </div>
      </div>
    </div>
  )
}
