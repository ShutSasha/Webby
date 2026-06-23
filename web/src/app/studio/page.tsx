import LockIcon from '@/assets/icons/shared/lock.svg'
import StudioVideosContainer from '@/ui/components/modules/Studio/StudioVideosContainer'
import CreateVideoButton from '@/ui/components/modules/Videos/CreateVideoButton'
import AuthPlaceholder from '@/ui/components/shared/AuthPlaceholder'
import { auth } from '@/workspace/auth'

export default async function StudioPage() {
  const session = await auth()

  if (!session) {
    return (
      <AuthPlaceholder
        title="Sign in to view your content"
        description="Please log in to manage your videos, playlists, and track your channel activity."
        icon={<LockIcon className="size-10 text-foreground0 stroke-1" />}
      />
    )
  }

  return (
    <>
      <div className="flex items-center justify-between mb-2">
        <h1 className="text-2xl font-bold text-foreground-secondary tracking-tight">Webby studio</h1>
        <CreateVideoButton />
      </div>

      <div className="flex flex-col flex-1 overflow-y-auto -mx-2 px-2">
        <div
          className="flex items-center text-xs font-medium text-foreground-muted border-b border-neutral-800 pb-3 px-4
            sticky top-0 bg-neutral-900 z-10"
        >
          <div className="flex-1 min-w-[300px]">Video</div>
          <div className="w-28 text-center shrink-0">Status</div>
          <div className="w-28 text-center shrink-0">Visibility</div>
          <div className="w-32 text-center shrink-0">Date</div>
          <div className="w-24 text-center shrink-0">Views</div>
          <div className="w-12 shrink-0"></div>
        </div>

        <StudioVideosContainer currentUserId={session.user.id} />
      </div>
    </>
  )
}
