import { getVideoInfo } from '@/lib/actions/video.actions'
import EditVideoContainer from '@/ui/components/modules/Videos/Edit/EditVideoContainer'
import EmptyState from '@/ui/components/shared/EmptyState'
import { auth } from '@/workspace/auth'

type Props = {
  params: Promise<{ id: string }>
}

export default async function EditVideoPage({ params }: Props) {
  const { id } = await params
  const [session, response] = await Promise.all([auth(), getVideoInfo(id)])

  if (!response.success || !response.data) {
    return (
      <div className="flex-1 flex items-center justify-center bg-neutral-900/20 rounded-[20px]">
        <EmptyState
          title="Video not found"
          description="We couldn't find the video you're looking for. It might have been deleted or you may not have permission to edit it."
        />
      </div>
    )
  }

  if (!session || session.user.id !== response.data.user.userId) {
    return (
      <div className="flex-1 flex items-center justify-center bg-neutral-900/20 rounded-[20px]">
        <EmptyState
          title="Access Denied"
          description="You don't have permission to edit this video. Make sure you are logged into the correct account."
        />
      </div>
    )
  }

  return (
    <div className="flex flex-col items-center w-full h-full p-6">
      <h1 className="text-center text-2xl font-bold text-neutral-100 mb-6 tracking-tight">Edit Video</h1>
      <EditVideoContainer initialVideo={response.data} />
    </div>
  )
}
