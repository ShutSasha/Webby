import { describe, it, expect } from 'vitest'

import { formatViews } from '../video.utils'

describe('formatViews utility', () => {
  it('should return "0 views" for zero', () => {
    expect(formatViews(0)).toBe('0 views')
  })

  it('should return "0 views" for negative numbers', () => {
    expect(formatViews(-5)).toBe('0 views')
    expect(formatViews(-1000)).toBe('0 views')
  })

  it('should return "0 views" for NaN', () => {
    expect(formatViews(NaN)).toBe('0 views')
  })

  it('should use singular "view" when count is exactly 1', () => {
    expect(formatViews(1)).toBe('1 view')
  })

  it('should use plural "views" for numbers greater than 1', () => {
    expect(formatViews(2)).toBe('2 views')
    expect(formatViews(999)).toBe('999 views')
  })

  it('should format thousands using "k"', () => {
    expect(formatViews(1000)).toBe('1k views')
    expect(formatViews(1500)).toBe('1.5k views')
    expect(formatViews(9999)).toBe('10k views')
    expect(formatViews(10500)).toBe('10.5k views')
    expect(formatViews(999000)).toBe('999k views')
  })

  it('should format millions using "m"', () => {
    expect(formatViews(1000000)).toBe('1m views')
    expect(formatViews(1500000)).toBe('1.5m views')
    expect(formatViews(25000000)).toBe('25m views')
  })

  it('should format billions using "b"', () => {
    expect(formatViews(1000000000)).toBe('1b views')
    expect(formatViews(3500000000)).toBe('3.5b views')
  })
})
