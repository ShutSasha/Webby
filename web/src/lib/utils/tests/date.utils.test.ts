import { describe, it, expect } from 'vitest'

import { formatDate } from '../date.utils'

describe('formatDate utility', () => {
  it('should format a valid ISO string into DD/MM/YYYY', () => {
    const validIsoString = '2026-03-16T17:53:56.468988Z'
    const result = formatDate(validIsoString)

    expect(result).toBe('16/03/2026')
  })

  it('should pad single-digit days and months with leading zeros', () => {
    const singleDigitDate = '2026-05-05T10:00:00.000Z'
    const result = formatDate(singleDigitDate)

    expect(result).toBe('05/05/2026')
  })

  it('should return "Invalid Date" when undefined is provided', () => {
    const result = formatDate(undefined)

    expect(result).toBe('Invalid Date')
  })

  it('should return "Invalid Date" when an empty string is provided', () => {
    const result = formatDate('')

    expect(result).toBe('Invalid Date')
  })

  it('should return "Invalid Date" when an invalid date string is provided', () => {
    const invalidDateString = 'not-a-real-date-string'
    const result = formatDate(invalidDateString)

    expect(result).toBe('Invalid Date')
  })
})
