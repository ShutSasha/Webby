import { Suspense } from 'react'

import MainLayout from '@/ui/components/MainLayout'
import UserProfileInfo, { UserProfileInfoSkeleton } from '@/ui/components/modules/Profile/UserProfileInfo'

type Props = {
  children: React.ReactNode
  params: Promise<{ id: string }>
}

export default async function ProfileLayout({ children, params }: Props) {
  const { id } = await params

  return (
    <MainLayout>
      <div className="flex flex-col gap-4 w-full max-w-5xl 2xl:max-w-7xl mx-auto">
        <Suspense fallback={<UserProfileInfoSkeleton />}>
          <UserProfileInfo id={id} />
        </Suspense>

        {children}
      </div>
    </MainLayout>
  )
}
