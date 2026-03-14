import { redirect } from 'next/navigation'

import { getAllAchievements } from '@/app/api/achievements'
import { getUser } from '@/app/api/user'
import AchievementItem from '@/ui/components/modules/Profile/Settings/Achievements/AchievementItem'
import Splitter from '@/ui/components/modules/Profile/Settings/Achievements/Splitter'
import UserHeader from '@/ui/components/modules/Profile/Settings/UserHeader'
import { auth } from '@/workspace/auth'

type Props = {
  params: Promise<{ id: string }>
}

export default async function AchievementsPage({ params }: Props) {
  const { id } = await params
  const session = await auth()

  // TODO: change to one edpoint
  const userData = await getUser(id)
  const allAchievements = await getAllAchievements()
  const unpinnedAchievements = allAchievements?.data?.filter(
    achievement => !userData?.pinnedUserAchievements.some(item => item.achievementId === achievement.achievementId),
  )

  if (!session?.user) {
    redirect('/login')
  }

  return (
    <div>
      <UserHeader image={session.user.image} username={session.user.username} />
      <Splitter text="Your pinned achievements" />
      <div className="flex flex-row items-center justify-center gap-3">
        {userData?.pinnedUserAchievements.map(achivement => (
          <AchievementItem
            title={achivement.title}
            description="Upload 5 videos on the Webby platform"
            key={achivement.achievementId}
            image={achivement.iconUrl}
            isPinned={true}
          />
        ))}
      </div>
      <Splitter text="All achievements" />
      <div className="grid grid-cols-5 gap-4">
        {unpinnedAchievements &&
          unpinnedAchievements.length > 0 &&
          unpinnedAchievements.map(achievement => (
            <AchievementItem
              title={achievement.title}
              description={achievement.description}
              key={achievement.achievementId}
              image={achievement.iconUrl}
              className="w-full"
              isPinned={false}
            />
          ))}
      </div>
    </div>
  )
}
