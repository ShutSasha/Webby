import type { NextAuthConfig } from 'next-auth'

export const authConfig = {
  pages: {
    signIn: '/', // або '/login'
  },
  callbacks: {
    async jwt({ token, user }) {
      if (user) {
        token._id = user._id
        token.username = user.username
        token.sub = user._id
      }
      return token
    },

    async session({ session, token }) {
      if (token && session.user) {
        session.user._id = token._id
        session.user.username = token.username
        session.user.name = token.username
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
