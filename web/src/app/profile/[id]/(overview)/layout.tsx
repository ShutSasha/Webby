import { Suspense } from 'react'

import UserProfileInfo, { UserProfileInfoSkeleton } from '@/ui/components/modules/Profile/UserProfileInfo'

type Props = {
  children: React.ReactNode
  params: Promise<{ id: string }>
}

export default async function ProfileLayout({ children, params }: Props) {
  const { id } = await params

  return (
    <>
      <Suspense fallback={<UserProfileInfoSkeleton />}>
        <UserProfileInfo id={id} />
      </Suspense>

      {children}
    </>
  )
}
