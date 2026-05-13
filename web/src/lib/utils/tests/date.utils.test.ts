import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'

import { formatDate, formatRelativeTime, formatTimeAgo, formatVideoTime, isValidDate } from '../date.utils'

describe('isValidDate utility', () => {
  it('should return true for valid modern dates', () => {
    expect(isValidDate('2026-03-16T17:53:56Z')).toBe(true)
    expect(isValidDate(new Date())).toBe(true)
    expect(isValidDate(1742136000000)).toBe(true)
  })

  it('should return true for old dates', () => {
    expect(isValidDate(5)).toBe(true)
    expect(isValidDate('1990-01-01')).toBe(true)
  })

  it('should return false for invalid data types or garbage strings', () => {
    expect(isValidDate(undefined)).toBe(false)
    expect(isValidDate(null)).toBe(false)
    expect(isValidDate('not-a-date')).toBe(false)
    expect(isValidDate({})).toBe(false)
  })
})

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

describe('formatRelativeTime utility', () => {
  it('should format date with month and time', () => {
    const date = '2026-05-11T12:00:00.000Z'
    const result = formatRelativeTime(date)
    expect(result).toContain('May 11, 03:00 PM')
  })

  it('should return "Unknown time" when null or undefined is provided', () => {
    expect(formatRelativeTime(undefined)).toBe('Unknown time')
    expect(formatRelativeTime(null)).toBe('Unknown time')
  })

  it('should return "Unknown time" when garbage string is provided', () => {
    expect(formatRelativeTime('garbage-string')).toBe('Unknown time')
  })
})

describe('formatVideoTime utility', () => {
  it('should format seconds into MM:SS', () => {
    expect(formatVideoTime(45)).toBe('0:45')
    expect(formatVideoTime(125)).toBe('2:05')
  })

  it('should format seconds into HH:MM:SS', () => {
    expect(formatVideoTime(3665)).toBe('1:01:05')
    expect(formatVideoTime(36650)).toBe('10:10:50')
  })

  it('should format seconds into DD:HH:MM:SS', () => {
    expect(formatVideoTime(366500)).toBe('4:05:48:20')
    expect(formatVideoTime(3665000)).toBe('11:10:03:20')
  })

  it('should return 00:00 for invalid or negative numbers', () => {
    expect(formatVideoTime(-10)).toBe('00:00')
    expect(formatVideoTime(NaN)).toBe('00:00')
  })
})
