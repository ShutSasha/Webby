import { redirect } from 'next/navigation'

import BadgeItem from '@/ui/components/modules/Profile/Settings/Badges/BadgeItem'
import Splitter from '@/ui/components/modules/Profile/Settings/Badges/Splitter'
import UserHeader from '@/ui/components/modules/Profile/Settings/UserHeader'
import { auth } from '@/workspace/auth'

export default async function BadgesPage() {
  const session = await auth()

  if (!session?.user) {
    redirect('/login')
  }

  return (
    <div>
      <UserHeader image={session.user.image} username={session.user.username} />
      <Splitter text="Your pinned badges" />
      <div className="flex flex-row items-center justify-center gap-3">
        {[...new Array(3)].map((_, index) => (
          <BadgeItem
            title="Title"
            description="Upload 5 videos on the Webby platform"
            key={index}
            image="https://i.ibb.co/ch9ZCDTr/badge1.png"
          />
        ))}
      </div>
      <Splitter text="All badges" />
      <div className="grid grid-cols-5 gap-4">
        {[...new Array(20)].map((_, index) => (
          <BadgeItem
            title="Title"
            description="Upload 5 videos on the Webby platform"
            key={index}
            image="https://i.ibb.co/ch9ZCDTr/badge1.png"
            className="w-full"
          />
        ))}
      </div>
    </div>
  )
}
