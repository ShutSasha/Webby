import Image from 'next/image'
import Link from 'next/link'

export default function FollowItem() {
  return (
    <Link
      href={'/'}
      className="bg-neutral-800 hover:bg-neutral-300/10 transition-all duration-200 ease-in rounded-xl p-2 flex
        items-center gap-2 border border-transparent hover:border-emerald-400/25"
    >
      <Image
        src="https://i.ibb.co/ccWcyJpy/3561466ac6f721d58fd40ff2dcbaa3d6.jpg"
        alt=""
        width={200}
        height={200}
        className="w-13 h-13 rounded-full"
      />
      <div className="flex flex-col">
        <p className="text-[18px] font-medium">Username</p>
        <p className="text-neutral-500 text-sm">152 followers</p>
      </div>
    </Link>
  )
}
