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
