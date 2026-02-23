import FollowItem from '@/ui/components/modules/Profile/CommonFollow/FollowItem'
import { FollowsToggle } from '@/ui/components/modules/Profile/CommonFollow/FollowToggle'
import HeaderNavigation from '@/ui/components/modules/Profile/CommonFollow/HeaderNavigation'

type Props = {
  params: Promise<{ id: string }>
}

export default async function Follows({ params }: Props) {
  const { id } = await params

  return (
    <div className="flex flex-col gap-2">
      {/* Header Navigation */}
      <HeaderNavigation id={id} image="https://i.pinimg.com/originals/44/64/20/4464203a781eed3650f1fdd624c4d02a.jpg">
        <FollowsToggle id={id} />
      </HeaderNavigation>

      <hr className="border-emerald-400/50 mb-1" />

      {/* Grid List */}
      <div className="grid grid-cols-1 md:grid-cols-2 2xl:grid-cols-3 gap-4">
        {[...new Array(20)].map((_, index) => (
          <FollowItem key={index} user={{ image: 'https://i.ibb.co/ccWcyJpy/3561466ac6f721d58fd40ff2dcbaa3d6.jpg' }} />
        ))}
      </div>
    </div>
  )
}
