import type { NextAuthConfig } from 'next-auth'

// TODO configure new user object

export const authConfig = {
  pages: {
    signIn: '/login',
  },
  callbacks: {
    async jwt({ token, user }) {
      if (user) {
        const u = user as any

        token.id = u.userId
        token.email = u.email
        token.image = u.avatarUrl
        token.username = u.username
      }
      return token
    },

    async session({ session, token }) {
      if (token && session.user) {
        session.user.id = token.id as string
        session.user.name = token.username as string
        session.user.email = token.email as string
        session.user.image = token.image as string
      }
      return session
    },
    authorized({ auth, request: { nextUrl } }) {
      const isLoggedIn = !!auth?.user
      const isOnChats = nextUrl.pathname.startsWith('/chats')

      if (isOnChats) {
        if (isLoggedIn) return true
        return false
      }
      return true
    },
  },
  providers: [],
} satisfies NextAuthConfig
