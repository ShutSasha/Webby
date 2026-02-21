import { notFound, redirect } from 'next/navigation'

import MainLayout from '@/ui/components/MainLayout'
import { auth } from '@/workspace/auth'

export default async function Settings({ params }: { params: Promise<{ id: string }> }) {
  const session = await auth()
  const { id } = await params

  if (!session?.user) {
    redirect('/login')
  }

  if (session.user.id !== id) {
    notFound()
  }

  return (
    <MainLayout>
      <div className="flex flex-col gap-4 w-full max-w-5xl 2xl:max-w-7xl mx-auto">
        <div className="bg-neutral-900 rounded-[20px] p-5 flex justify-between gap-4">
          <p>settings page</p>
        </div>
      </div>
    </MainLayout>
  )
}
