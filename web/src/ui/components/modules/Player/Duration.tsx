import { formatVideoTime } from '@/lib/utils/date.utils'

export default function Duration({ className, seconds }: { className?: string; seconds: number }) {
  return (
    <time dateTime={`P${Math.round(seconds)}S`} className={className} style={{ fontFeatureSettings: '"tnum"' }}>
      {formatVideoTime(seconds)}
    </time>
  )
}
