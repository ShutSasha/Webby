const CHAT_COLORS = [
  'text-red-400',
  'text-orange-400',
  'text-amber-400',
  'text-emerald-400',
  'text-cyan-400',
  'text-blue-400',
  'text-indigo-400',
  'text-violet-400',
  'text-purple-400',
  'text-fuchsia-400',
  'text-rose-400',
]

const colorCache = new Map<string, string>()

export const getUserColor = (userId: string): string => {
  if (colorCache.has(userId)) {
    return colorCache.get(userId)!
  }

  let hash = 0
  for (let i = 0; i < userId.length; i++) {
    hash = userId.charCodeAt(i) + ((hash << 5) - hash)
  }

  const index = Math.abs(hash) % CHAT_COLORS.length
  const color = CHAT_COLORS[index]

  colorCache.set(userId, color)

  return color
}
