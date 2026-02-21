'use server'

import { Route } from 'next'
import { redirect } from 'next/dist/client/components/navigation'
import { cookies } from 'next/headers'
import { AuthError, type User } from 'next-auth'

import $api from '@/app/api'
import { serverLog } from '@/lib/utils/utils'

import { signIn } from '../../../auth'

type ReturnError = {
  success: boolean
  errors: Record<string, string>
}

export async function login(email: string, password: string): Promise<User | ReturnError> {
  try {
    const { data: response } = await $api.post('/auth/sign-in', { email, password })

    return response.data as User
  } catch (error: any) {
    serverLog('LOGIN_ERROR', error, true)
    const errors = error.response?.data?.errors
    const message = error.response?.data?.message
    return {
      success: false,
      errors: errors ? errors : message ? { global: message } : { global: 'unknown login error' },
    }
  }
}

export async function authenticate(prevState: any, formData: FormData) {
  try {
    await signIn('credentials', {
      ...Object.fromEntries(formData),
      redirect: false,
    })

    return { success: true, errors: null }
  } catch (error) {
    if (error instanceof AuthError) {
      const cause = error.cause?.err?.message

      try {
        if (cause) {
          const serverErrors = JSON.parse(cause)
          return {
            success: false,
            errors: serverErrors,
          }
        }
      } catch {}

      return {
        success: false,
        errors: { email: 'Unknown login error' },
      }
    }

    throw error
  }
}

export async function registerUser(prevState: any, formData: FormData) {
  const email = formData.get('email') as string
  const username = formData.get('username') as string
  const password = formData.get('password') as string
  const repeatPassword = formData.get('repeatPassword') as string

  if (password !== repeatPassword) {
    return {
      success: false,
      errors: { repeatPassword: 'Password does not match' },
    }
  }

  try {
    await $api.post('/auth/sign-up', {
      email,
      username,
      password,
    })
  } catch (error: any) {
    serverLog('REGISTER_ERROR', error, true)

    const serverErrors = error.response?.data?.errors || {}

    return {
      success: false,
      errors: serverErrors,
    }
  }

  redirect(`/sign-up/email-verify?email=${encodeURIComponent(email)}` satisfies Route)
}

export async function verifyUser({ email, code }: { email: string; code: string }) {
  try {
    await $api.post('/auth/verify-user', {
      email,
      verificationCode: code,
    })
  } catch (error: any) {
    serverLog('VERIFY_USER_ERROR', error, true)

    const serverErrors = error.response?.data?.errors || {}

    return {
      success: false,
      errors: serverErrors,
    }
  }
}

export async function googleAuthenticate(prevState: string | undefined, formData: FormData) {
  try {
    await signIn('google')
  } catch (error) {
    if (error instanceof AuthError) {
      return 'google log in failed'
    }
    throw error
  }
}

export async function getEncryptedToken() {
  const cookieStore = await cookies()
  const token =
    cookieStore.get('authjs.session-token')?.value || cookieStore.get('__Secure-authjs.session-token')?.value

  return token
}
