/**
 * Converts an ISO date string (e.g., 2026-03-16T17:53:56.468988Z) into dd/mm/yyyy format.
 * * @param dateString - The ISO date string retrieved from the server.
 * @returns A formatted date string or 'Unknown time' if input is undefined.
 */
export const formatDate = (dateString: string | Date | number | unknown): string => {
  if (
    !dateString ||
    (typeof dateString !== 'string' && typeof dateString !== 'number' && !(dateString instanceof Date))
  ) {
    return 'Unknown time'
  }

  if (!dateString) return 'Unknown time'

  const date = new Date(dateString)

  if (isNaN(date.getTime())) {
    return 'Unknown time'
  }

  /**
   * We use 'en-GB' locale because it natively uses the DD/MM/YYYY format with slashes.
   * This eliminates the need for manual string manipulation like .split().join().
   */
  return new Intl.DateTimeFormat('en-GB', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
  }).format(date)
}

const TIME_INTERVALS = [
  { label: 'year', seconds: 31536000 },
  { label: 'month', seconds: 2592000 },
  { label: 'week', seconds: 604800 },
  { label: 'day', seconds: 86400 },
  { label: 'hour', seconds: 3600 },
  { label: 'minute', seconds: 60 },
] as const

let rtfCache: Intl.RelativeTimeFormat | null = null

export function formatTimeAgo(dateInput: string | Date | number | unknown): string {
  if (!dateInput || (typeof dateInput !== 'string' && typeof dateInput !== 'number' && !(dateInput instanceof Date))) {
    return 'Unknown time'
  }

  const date = new Date(dateInput)

  if (isNaN(date.getTime())) {
    return 'Unknown time'
  }

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

export function formatVideoTime(seconds: number) {
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

export function formatRelativeTime(dateString: string) {
  const date = new Date(dateString)
  return new Intl.DateTimeFormat('en-US', {
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  }).format(date)
}
