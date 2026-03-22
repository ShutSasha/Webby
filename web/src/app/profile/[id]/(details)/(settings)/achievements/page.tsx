import { notFound, redirect } from 'next/navigation'

import { getUserAchievements } from '@/app/api/achievements'
import { clog } from '@/lib/utils/utils'
import AchievementItem from '@/ui/components/modules/Profile/Settings/Achievements/AchievementItem'
import Splitter from '@/ui/components/modules/Profile/Settings/Achievements/Splitter'
import UserHeader from '@/ui/components/modules/Profile/Settings/UserHeader'
import { auth } from '@/workspace/auth'

type Props = {
  params: Promise<{ id: string }>
}

export default async function AchievementsPage({ params }: Props) {
  const { id: userId } = await params
  const session = await auth()
  const userAchievements = await getUserAchievements(userId)

  clog('userAchievements', userAchievements)

  if (!session?.user) {
    redirect('/login')
  }

  if (!userAchievements.success || !userAchievements.data) {
    notFound()
  }

  return (
    <div>
      <UserHeader image={session.user.image} username={session.user.username} />
      <Splitter text="Your pinned achievements" />
      <div className="flex flex-row items-center justify-center gap-3">
        {userAchievements.data.pinnedAchievements.map(achivement => (
          <AchievementItem
            achievementId={achivement.achievementId}
            title={achivement.title}
            description={achivement.description}
            key={achivement.achievementId}
            image={achivement.iconUrl}
            isPinned={true}
            isUnlocked
          />
        ))}
      </div>
      <Splitter text="All achievements" />
      <div className="grid grid-cols-5 gap-4">
        {userAchievements.data.achievements.map(achievement => (
          <AchievementItem
            achievementId={achievement.achievementId}
            title={achievement.title}
            description={achievement.description}
            key={achievement.achievementId}
            image={achievement.iconUrl}
            className="w-full"
            isPinned={false}
            isUnlocked={achievement.isUnlocked}
          />
        ))}
      </div>
    </div>
  )
}
