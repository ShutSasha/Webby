import { FC, SVGProps } from 'react'

interface ProfileActionButtonProps {
  Icon: FC<SVGProps<SVGSVGElement>>
  label: string
  onClick?: () => void
  iconClassName?: string
  btnClassName?: string
}

export default async function ProfileActionButton({
  Icon,
  label,
  onClick,
  iconClassName,
  btnClassName,
}: ProfileActionButtonProps) {
  return (
    <button
      onClick={onClick}
      className={`flex items-center gap-2 px-4 py-2 bg-neutral-800 rounded-full transition-all duration-300
        hover:bg-neutral-700/40 cursor-pointer active:scale-90 group border border-transparent ${btnClassName}`}
    >
      <Icon className={`w-4 h-4 transition-colors ${iconClassName}`} />

      <span className="text-white text-sm font-medium leading-none">{label}</span>
    </button>
  )
}
