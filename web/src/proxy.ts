import NextAuth from 'next-auth'

import { authConfig } from '../auth.config'

const { auth } = NextAuth(authConfig)

export default auth(req => {
  const { nextUrl } = req
  const isLoggedIn = !!req.auth
  const userRole = req.auth?.user?.role

  const isAuthPage = nextUrl.pathname.startsWith('/login') || nextUrl.pathname.startsWith('/sign-up')
  const isOnSettings = nextUrl.pathname.endsWith('/settings')
  const isAdminPage = nextUrl.pathname.startsWith('/admin')

  if (isAuthPage) {
    if (isLoggedIn) {
      return Response.redirect(new URL('/', nextUrl))
    }
    return
  }

  if (isAdminPage) {
    if (!isLoggedIn) {
      return Response.redirect(new URL('/login', nextUrl))
    }
    if (!['Admin', 'Moderator'].includes(userRole)) {
      return Response.redirect(new URL('/', nextUrl))
    }
  }

  if (isOnSettings && !isLoggedIn) {
    return Response.redirect(new URL('/login', nextUrl))
  }
})

export const config = {
  matcher: ['/((?!api|_next/static|_next/image|favicon.ico).*)'],
}
