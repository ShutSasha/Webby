import LockIcon from '@/assets/icons/shared/lock.svg'
import { getUserAchievements } from '@/lib/actions/achievement.actions'
import AchievementItem from '@/ui/components/modules/Profile/Settings/Achievements/AchievementItem'
import UserHeader from '@/ui/components/modules/Profile/Settings/UserHeader'
import AuthPlaceholder from '@/ui/components/shared/AuthPlaceholder'
import EmptyState from '@/ui/components/shared/EmptyState'
import { auth } from '@/workspace/auth'

type Props = {
  params: Promise<{ id: string }>
}

export default async function AchievementsPage({ params }: Props) {
  const { id: userId } = await params
  const session = await auth()
  const userAchievements = await getUserAchievements(userId)

  if (!session) {
    return (
      <AuthPlaceholder
        title="Sign in to view your achievements"
        description="Please log in to unlock your trophy room. Track your completed milestones,
        check your ongoing quest progress, and view all the exclusive badges you have earned on the platform."
        icon={<LockIcon className="size-10 text-foreground-faint stroke-1" />}
      />
    )
  }

  if (!userAchievements.success || !userAchievements.data) {
    return (
      <div className="flex flex-col gap-6 pb-10 h-full">
        <UserHeader image={session.user.image} username={session.user.username} />
        <div className="flex-1 mt-10">
          <EmptyState
            title="Couldn't load achievements"
            description="We ran into a problem while retrieving the achievements. Please try refreshing the page."
          />
        </div>
      </div>
    )
  }

  const pinned = userAchievements.data.pinnedAchievements
  const all = userAchievements.data.achievements

  return (
    <div className="flex flex-col pb-10">
      <UserHeader image={session.user.image} username={session.user.username} />

      <div className="w-full flex flex-col items-center gap-6 mt-6">
        <div className="w-full max-w-[800px] border border-border rounded-2xl p-6 sm:p-8 flex flex-col">
          <div className="flex flex-col mb-6">
            <h2 className="text-[11px] font-bold text-foreground-muted uppercase tracking-wider mb-2">
              Pinned Achievements
            </h2>
            <p className="text-sm text-foreground-subtle leading-relaxed">
              Showcase your proudest milestones. You can pin any unlocked achievement to easily access and display it
              here.
            </p>
          </div>

          {pinned.length === 0 ? (
            <div
              className="w-full flex flex-col items-center justify-center py-10 px-4 border-2 border-dashed
                border-border rounded-[20px] bg-surface/20"
            >
              <p className="text-foreground-muted font-medium">No pinned achievements yet</p>
              <p className="text-foreground-disabled text-sm mt-1 text-center max-w-sm">
                Click on the three dots of any unlocked achievement below to pin it here.
              </p>
            </div>
          ) : (
            <div className="flex flex-row flex-wrap items-center justify-center gap-3">
              {pinned.map(achievement => (
                <AchievementItem
                  achievementId={achievement.achievementId}
                  title={achievement.title}
                  description={achievement.description}
                  key={`pinned-${achievement.achievementId}`}
                  image={achievement.iconUrl}
                  isPinned={true}
                  isUnlocked={true}
                  achievementProgressValue={achievement.achievementProgressValue}
                  targetValue={achievement.targetValue}
                  unlockedAt={achievement.unlockedAt}
                />
              ))}
            </div>
          )}
        </div>

        <div className="w-full max-w-[800px] border border-border rounded-2xl p-6 sm:p-8 flex flex-col">
          <div className="flex flex-col mb-6">
            <h2 className="text-[11px] font-bold text-foreground-muted uppercase tracking-wider mb-2">
              All Achievements
            </h2>
            <p className="text-sm text-foreground-subtle leading-relaxed">
              Track your progress, complete quests, and collect exclusive badges across the platform.
            </p>
          </div>

          {all.length === 0 ? (
            <div className="flex items-center justify-center py-12">
              <p className="text-foreground-faint">There are no achievements available.</p>
            </div>
          ) : (
            <div className="flex flex-row flex-wrap items-center justify-center gap-3">
              {all.map(achievement => (
                <AchievementItem
                  achievementId={achievement.achievementId}
                  title={achievement.title}
                  description={achievement.description}
                  key={`all-${achievement.achievementId}`}
                  image={achievement.iconUrl}
                  isPinned={false}
                  isUnlocked={achievement.isUnlocked}
                  achievementProgressValue={achievement.achievementProgressValue}
                  targetValue={achievement.targetValue}
                  unlockedAt={achievement.unlockedAt}
                />
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
