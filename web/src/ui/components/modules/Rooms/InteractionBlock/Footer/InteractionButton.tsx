import { cn } from '@/lib/utils/general.utils'

type InteractionButtonProps = {
  icon: React.ReactNode
  text?: React.ReactNode
  onClick?: () => void
  className?: string
  children?: React.ReactNode
}

export default function InteractionButton({ icon, text, onClick, className, children }: InteractionButtonProps) {
  return (
    <div
      onClick={onClick}
      className={cn(
        'relative flex flex-row items-center bg-surface py-1.5 px-3 gap-2.5 rounded-md group select-none',
        'hover:bg-neutral-100 dark:hover:bg-emerald-500 transition-colors duration-300 ease-in-out cursor-pointer',
        className,
      )}
    >
      {icon}

      {text && (
        <p
          className="dark:group-hover:text-foreground-inverse-subtle text-foreground-subtle transition-colors
            duration-300 ease-in-out leading-5 font-medium"
        >
          {text}
        </p>
      )}
      {children}
    </div>
  )
}
