import type { NextAuthConfig } from 'next-auth'

export const authConfig = {
  pages: {
    signIn: '/login',
  },
  callbacks: {
    async jwt({ token, user }) {
      if (user) {
        token.id = user.userId
        token.username = user.username
        token.email = user.email
        token.image = user.avatarUrl
      }
      return token
    },

    async session({ session, token }) {
      if (token && session.user) {
        session.user.id = token.id
        session.user.username = token.username
        session.user.email = token.email
        session.user.image = token.image
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
