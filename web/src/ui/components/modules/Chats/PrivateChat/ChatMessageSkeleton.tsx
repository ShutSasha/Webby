import { cn } from '@/lib/utils/general.utils'

type Props = {
  isMe?: boolean
  className?: string
}

export default function ChatMessageSkeleton({ isMe, className }: Props) {
  return (
    <div
      className={cn(
        'animate-pulse rounded-2xl',
        isMe ? 'self-end bg-emerald-500/20 rounded-br-sm' : 'self-start bg-neutral-800 rounded-bl-sm',
        className,
      )}
    />
  )
}
