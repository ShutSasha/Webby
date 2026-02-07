import type { NextAuthConfig } from 'next-auth'

export const authConfig = {
  pages: {
    signIn: '/login',
  },
  callbacks: {
    async signIn({ user, account }) {
      if (account?.provider === 'google') {
        try {
          console.log('authConfig.signIn - Google user:', user)

          user.userId = user.id as string
          user.username = user.name as string
          user.avatarUrl = user.image as string
          return true
        
          // TODO Sync with backend

          // if (backendUser) {
          //   user.userId = backendUser.id
          //   user.username = backendUser.username
          //   user.avatarUrl = backendUser.image
          //   return true
          // }
        } catch (error) {
          console.error('Backend sync error:', error)
          return false
        }
      }
      return true
    },

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
