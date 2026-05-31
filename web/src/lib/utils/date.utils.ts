const gbDateFormatter = new Intl.DateTimeFormat('en-GB', {
  day: '2-digit',
  month: '2-digit',
  year: 'numeric',
})

const usRelativeDateFormatter = new Intl.DateTimeFormat('en-US', {
  month: 'short',
  day: 'numeric',
  hour: '2-digit',
  minute: '2-digit',
})

const TIME_INTERVALS = [
  { label: 'year', seconds: 31536000 },
  { label: 'month', seconds: 2592000 },
  { label: 'week', seconds: 604800 },
  { label: 'day', seconds: 86400 },
  { label: 'hour', seconds: 3600 },
  { label: 'minute', seconds: 60 },
] as const

let rtfCache: Intl.RelativeTimeFormat | null = null

export function isValidDate(input: unknown): input is string | number | Date {
  if (!input || (typeof input !== 'string' && typeof input !== 'number' && !(input instanceof Date))) {
    return false
  }

  const date = new Date(input)
  if (isNaN(date.getTime())) {
    return false
  }

  return true
}

/**
 * Converts a date input into a standard British format: DD/MM/YYYY.
 * Example: '16/03/2026'.
 * * @param dateInput - The date string, number, or Date object.
 * @returns Formatted date string or 'Unknown time'.
 */
export const formatDate = (dateInput: unknown): string => {
  if (!isValidDate(dateInput)) return 'Unknown time'

  const date = new Date(dateInput)

  return gbDateFormatter.format(date)
}

/**
 * Calculates the relative time difference between now and the provided date.
 * Example: '5 minutes ago', 'just now', '2 days ago'.
 * * @param dateInput - The date string, number, or Date object.
 * @returns Relative time string or 'Unknown time'.
 */
export function formatTimeAgo(dateInput: unknown): string {
  if (!isValidDate(dateInput)) return 'Unknown time'

  const date = new Date(dateInput)
  const diffInSeconds = Math.floor((Date.now() - date.getTime()) / 1000)

  if (diffInSeconds < 60) {
    return 'just now'
  }

  if (!rtfCache) {
    rtfCache = new Intl.RelativeTimeFormat('en-US', { numeric: 'always' })
  }

  for (const interval of TIME_INTERVALS) {
    const count = Math.floor(diffInSeconds / interval.seconds)

    if (count >= 1) {
      return rtfCache.format(-count, interval.label)
    }
  }

  return 'just now'
}

/**
 * Formats video duration in seconds into a player-friendly format.
 * Example: 45 -> '00:45', 3665 -> '1:01:05'.
 * * @param seconds - Total video duration in seconds.
 * @returns Formatted duration string: [D:]HH:MM:SS or MM:SS.
 */
export function formatVideoTime(seconds: number) {
  if (isNaN(seconds) || seconds < 0) return '00:00'

  const date = new Date(seconds * 1000)
  const dd = date.getUTCDate() - 1
  const hh = date.getUTCHours()
  const mm = date.getUTCMinutes()
  const ss = pad(date.getUTCSeconds())
  if (dd) {
    return `${dd}:${pad(hh)}:${pad(mm)}:${ss}`
  }
  if (hh) {
    return `${hh}:${pad(mm)}:${ss}`
  }
  return `${mm}:${ss}`
}

function pad(string: string | number) {
  return `0${string}`.slice(-2)
}

/**
 * Formats a date into a short US format with time.
 * Example: 'May 11, 12:00 PM'.
 * * @param dateInput - The date string, number, or Date object.
 * @returns Formatted string or 'Unknown time'.
 */
export function formatRelativeTime(dateInput: unknown) {
  if (!isValidDate(dateInput)) return 'Unknown time'
  const date = new Date(dateInput)
  return usRelativeDateFormatter.format(date)
}
