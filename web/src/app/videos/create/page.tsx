import LockIcon from '@/assets/icons/shared/lock.svg'
import CreateVideoContainer from '@/ui/components/modules/Videos/Create/CreateVideoContainer'
import AuthPlaceholder from '@/ui/components/shared/AuthPlaceholder'
import { auth } from '@/workspace/auth'

export default async function CreateVideoPage() {
  const session = await auth()

  if (!session) {
    return (
      <AuthPlaceholder
        title="Sign in to create your videos"
        description="Please log in to upload your custom video files and create room sessions to share and watch content live with your friends."
        icon={<LockIcon className="size-10 text-foreground0 stroke-1" />}
      />
    )
  }

  return (
    <div className="flex flex-col items-center w-full h-full p-6">
      <h1 className="text-center text-2xl font-bold text-foreground-secondary mb-6 tracking-tight">Create Video</h1>
      <CreateVideoContainer />
    </div>
  )
}
