import FollowItem from './FollowItem'

export default async function FollowsContainer() {
  await new Promise(r => setTimeout(r, 1300))

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 2xl:grid-cols-3 gap-4">
      {[...new Array(20)].map((_, index) => (
        <FollowItem key={index} user={{ image: 'https://i.ibb.co/ccWcyJpy/3561466ac6f721d58fd40ff2dcbaa3d6.jpg' }} />
      ))}
    </div>
  )
}
