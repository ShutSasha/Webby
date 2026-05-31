const compactNumberFormatter = new Intl.NumberFormat('en-US', {
  notation: 'compact',
  maximumFractionDigits: 1,
})

/**
 * Formats a raw view count into a compact, human-readable string (e.g., 1.5k, 2m).
 * Appends the correct singular or plural label based on the count.
 *
 * @param views - The total number of views (must be a positive number).
 * @returns A formatted string such as '1 view', '1.5k views', or '0 views' as a fallback.
 */
export function formatViews(views: number): string {
  if (isNaN(views) || views < 0) {
    return '0 views'
  }

  const formattedNumber = compactNumberFormatter.format(views).toLowerCase()
  const label = views === 1 ? 'view' : 'views'

  return `${formattedNumber} ${label}`
}
