import { redirect } from 'next/navigation'

import { auth } from '@/workspace/auth'

export default async function SecurePage() {
  const session = await auth()

  if (!session?.user) {
    redirect('/login')
  }

  return (
    <div className="flex flex-col">
      <p>user id: {session.user.id}</p>
    </div>
  )
}
