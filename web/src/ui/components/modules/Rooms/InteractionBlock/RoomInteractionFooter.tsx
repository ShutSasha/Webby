import CubeIcon from '@/assets/icons/Room/cube-points.svg'
import PollIcon from '@/assets/icons/Room/poll-icon.svg'
import ReactionSmileIcon from '@/assets/icons/Room/reaction-smile.svg'

export default function RoomInteractionFooter() {
  return (
    <div className="flex flex-row justify-between items-center">
      <div className="flex items-center gap-2">
        <UserPoints />
        <Reactions />
      </div>
      <Poll />
    </div>
  )
}

function UserPoints() {
  return (
    <div
      className="flex flex-row items-center bg-neutral-900 py-1.5 px-3 gap-2.5 rounded-md group hover:bg-emerald-500
        transition-colors duration-300 ease-in-out cursor-pointer"
    >
      <CubeIcon
        className="size-5 text-neutral-300 group-hover:text-neutral-900 transition-colors duration-300 ease-in-out"
      />
      <p className="group-hover:text-neutral-900 text-neutral-300 transition-colors duration-300 ease-in-out leading-5">
        12
      </p>
    </div>
  )
}

function Reactions() {
  return (
    <div
      className="flex flex-row items-center bg-neutral-900 py-1.5 px-3 gap-2.5 rounded-md group hover:bg-emerald-500
        transition-colors duration-300 ease-in-out cursor-pointer"
    >
      <ReactionSmileIcon
        className="size-5 text-neutral-300 group-hover:text-neutral-900 transition-colors duration-300 ease-in-out"
      />
    </div>
  )
}

function Poll() {
  return (
    <div
      className="flex flex-row items-center bg-neutral-900 py-1.5 px-3 gap-2.5 rounded-md group hover:bg-emerald-500
        transition-colors duration-300 ease-in-out cursor-pointer"
    >
      <PollIcon
        className="size-5 text-neutral-300 group-hover:text-neutral-900 transition-colors duration-300 ease-in-out"
      />
      <p className="group-hover:text-neutral-900 text-neutral-300 transition-colors duration-300 ease-in-out leading-5">
        Poll
      </p>
    </div>
  )
}
