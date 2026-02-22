import { notFound, redirect } from 'next/navigation'

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
    <div>
      <p>settings</p>
      <p>
        Lorem ipsum dolor sit amet, consectetur adipisicing elit. Totam saepe blanditiis cupiditate beatae enim minima
        voluptas mollitia? Necessitatibus, odit distinctio.
      </p>
      <p>
        Lorem ipsum dolor sit amet, consectetur adipisicing elit. Totam saepe blanditiis cupiditate beatae enim minima
        voluptas mollitia? Necessitatibus, odit distinctio.
      </p>
      <p>
        Lorem ipsum dolor sit amet, consectetur adipisicing elit. Totam saepe blanditiis cupiditate beatae enim minima
        voluptas mollitia? Necessitatibus, odit distinctio.
      </p>
    </div>
  )
}
