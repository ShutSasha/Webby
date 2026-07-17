import FilledPlay from '@/assets/icons/Player/filled-play.svg'

type Props = {
  playing: boolean
  buffering: boolean
}

export default function PlayerCenterButton({ playing, buffering }: Props) {
  if (buffering) return null

  return (
    <div
      className={`${playing === false ? 'opacity-100' : 'opacity-0'} absolute w-12 h-12 md:w-18 md:h-18 bg-black/40
        rounded-full top-1/2 left-1/2 -translate-y-1/2 -translate-x-1/2 flex items-center justify-center pl-1
        transition-opacity duration-400 ease-in-out pointer-events-none`}
    >
      <FilledPlay className="w-5 h-5 md:w-8 md:h-8 text-white/90" />
    </div>
  )
}
