export function formatViews(views: number): string {
  const formatter = new Intl.NumberFormat('en-US', {
    notation: 'compact',
    maximumFractionDigits: 1,
  })

  const formattedNumber = formatter.format(views).toLowerCase()

  const label = views === 1 ? 'view' : 'views'

  return `${formattedNumber} ${label}`
}
