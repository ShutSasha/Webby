import CreateVideoContainer from '@/ui/components/modules/Videos/Create/CreateVideoContainer'

export default function CreateVideoPage() {
  return (
    <div className="flex flex-col items-center w-full h-full p-6">
      <h1 className="text-center text-2xl font-bold text-neutral-100 mb-6 tracking-tight">Create Video</h1>
      <CreateVideoContainer />
    </div>
  )
}
