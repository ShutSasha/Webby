import FollowItem from './FollowItem'

export default async function FollowersContainer() {
  await new Promise(r => setTimeout(r, 1300))

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 2xl:grid-cols-3 gap-4">
      {[...new Array(20)].map((_, index) => (
        <FollowItem key={index} user={{ image: 'https://i.ibb.co/9fP2m4R/d5beed264cd6af09fc7a53548326c4bf.jpg' }} />
      ))}
    </div>
  )
}
