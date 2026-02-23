'use server'

import { Route } from 'next'
import { redirect } from 'next/dist/client/components/navigation'
import { cookies } from 'next/headers'
import { AuthError } from 'next-auth'

import $api from '@/app/api'
import { parseAxiosError, serverLog } from '@/lib/utils/utils'
import { LoginRes } from '@/types/auth'
import { FormActionState } from '@/types/general'
import { signIn } from '@/workspace/auth'

type ReturnError = {
  success: boolean
  errors: Record<string, string>
}

export async function login(email: string, password: string): Promise<LoginRes | ReturnError> {
  try {
    const { data: response } = await $api.post('/auth/sign-in', { email, password })

    return response.data as LoginRes
  } catch (error: unknown) {
    serverLog('LOGIN_ERROR', error, true)

    return {
      success: false,
      errors: parseAxiosError(error),
    }
  }
}

export async function authenticate(_prevState: FormActionState, formData: FormData) {
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
          return {
            success: false,
            errors: JSON.parse(cause),
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

export async function registerUser(_prevState: FormActionState, formData: FormData) {
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
  } catch (error: unknown) {
    serverLog('REGISTER_ERROR', error, true)

    return {
      success: false,
      errors: parseAxiosError(error),
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
  } catch (error: unknown) {
    serverLog('VERIFY_USER_ERROR', error, true)

    return {
      success: false,
      errors: parseAxiosError(error),
    }
  }
}

export async function getEncryptedToken() {
  const cookieStore = await cookies()
  const token =
    cookieStore.get('authjs.session-token')?.value || cookieStore.get('__Secure-authjs.session-token')?.value

  return token
}
