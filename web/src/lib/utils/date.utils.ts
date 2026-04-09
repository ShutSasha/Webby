/**
 * Converts an ISO date string (e.g., 2026-03-16T17:53:56.468988Z) into dd/mm/yyyy format.
 * * @param dateString - The ISO date string retrieved from the server.
 * @returns A formatted date string or an empty string if input is undefined.
 */
export const formatDate = (dateString: string | undefined): string => {
  if (!dateString) return ''

  const date = new Date(dateString)

  if (isNaN(date.getTime())) {
    return 'Invalid Date'
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

export function formatTimeAgo(dateInput: string | Date | number): string {
  const date = new Date(dateInput)
  const now = new Date()

  const diffInSeconds = Math.floor((now.getTime() - date.getTime()) / 1000)

  if (diffInSeconds < 60) {
    return 'just now'
  }

  const rtf = new Intl.RelativeTimeFormat('en-US', { numeric: 'always' })

  const intervals = [
    { label: 'year', seconds: 31536000 },
    { label: 'month', seconds: 2592000 },
    { label: 'week', seconds: 604800 },
    { label: 'day', seconds: 86400 },
    { label: 'hour', seconds: 3600 },
    { label: 'minute', seconds: 60 },
  ] as const

  for (const interval of intervals) {
    const count = Math.floor(diffInSeconds / interval.seconds)

    if (count >= 1) {
      return rtf.format(-count, interval.label)
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
