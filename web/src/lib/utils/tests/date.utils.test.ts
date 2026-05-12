import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'

import { formatDate, formatTimeAgo } from '../date.utils'

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

  it('should return "Unknown time" when undefined or null is provided', () => {
    expect(formatDate(undefined)).toBe('Unknown time')
    expect(formatDate(null)).toBe('Unknown time')
  })

  it('should return "Unknown time" when an empty string is provided', () => {
    expect(formatDate('')).toBe('Unknown time')
  })

  it('should return "Unknown time" when an invalid string is provided', () => {
    const invalidDateString = 'not-a-real-date-string'

    expect(formatDate(invalidDateString)).toBe('Unknown time')
  })
})

describe('formatTimeAgo utility', () => {
  beforeEach(() => {
    vi.useFakeTimers()
    vi.setSystemTime(new Date('2026-05-11T12:00:00.000Z'))
  })

  afterEach(() => {
    vi.useRealTimers()
  })

  it('should return "just now" for times less than 60 seconds ago', () => {
    const date = new Date('2026-05-11T11:59:10.000Z')
    expect(formatTimeAgo(date)).toBe('just now')
  })

  it('should return relative minutes', () => {
    const date = new Date('2026-05-11T11:55:00.000Z')
    expect(formatTimeAgo(date)).toBe('5 minutes ago')
  })

  it('should return relative hours', () => {
    const date = new Date('2026-05-11T09:00:00.000Z')
    expect(formatTimeAgo(date)).toBe('3 hours ago')
  })

  it('should return relative days', () => {
    const date = new Date('2026-05-09T12:00:00.000Z')
    expect(formatTimeAgo(date)).toBe('2 days ago')
  })

  it('should return "Unknown time" when invalid garbage string is provided', () => {
    expect(formatTimeAgo('server-sent-garbage')).toBe('Unknown time')
  })

  it('should return "Unknown time" when null or undefined is provided', () => {
    expect(formatTimeAgo(undefined)).toBe('Unknown time')
    expect(formatTimeAgo(null)).toBe('Unknown time')
  })

  it('should handle future dates safely due to clock desync (return "just now")', () => {
    const futureDate = new Date('2026-05-11T12:05:00.000Z')
    expect(formatTimeAgo(futureDate)).toBe('just now')
  })
})
