import FollowItem from '@/ui/components/modules/Profile/FollowItem'

export default async function Follows() {
  return (
    <div className="flex flex-col gap-3">
      <p className="text-center text-[20px] font-semibold">Follows</p>
      <hr className="text-emerald-400" />
      <div className="grid grid-cols-4 gap-4">
        {[...new Array(20)].map((it, index) => (
          <FollowItem key={index} />
        ))}
      </div>
    </div>
  )
}
