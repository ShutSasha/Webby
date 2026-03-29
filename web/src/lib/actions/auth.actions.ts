'use server'

import { Route } from 'next'
import { redirect } from 'next/dist/client/components/navigation'
import { AuthError } from 'next-auth'

import $api from '@/lib/api/api.config'
import { parseAxiosError, serverLog } from '@/lib/utils/utils'
import { LoginRes } from '@/types/auth'
import { FormActionState } from '@/types/general'
import { signIn } from '@/workspace/auth'

export type ReturnAuthError = {
  success: boolean
  errors: Record<string, string>
}

export async function login(email: string, password: string): Promise<LoginRes | ReturnAuthError> {
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

export async function resendVerifyCode({ email }: { email: string }) {
  try {
    await $api.post('/auth/resend-verification-code', { email })
  } catch (error) {
    serverLog('Error resending verification code:', error)
  }
}

export async function changePassworddAction(_prevState: FormActionState, formData: FormData) {
  const currentPassword = formData.get('currentPassword') as string
  const newPassword = formData.get('newPassword') as string
  const confirmPassword = formData.get('confirmPassword') as string

  if (newPassword !== confirmPassword) {
    return {
      success: false,
      errors: { confirmPassword: 'Passwords do not match' },
    }
  }

  try {
    await $api.patch('/auth/change-password', {
      currentPassword,
      newPassword,
    })

    return { success: true, errors: null, timestamp: Date.now() }
  } catch (error: unknown) {
    serverLog('CHANGE_PASSWORD_ERROR', error, true)
    return { success: false, errors: parseAxiosError(error) }
  }
}

export async function resetPasswordAction(_prevState: FormActionState, formData: FormData) {
  const email = formData.get('email') as string
  const newPassword = formData.get('newPassword') as string
  const confirmPassword = formData.get('confirmPassword') as string

  if (newPassword !== confirmPassword) {
    return {
      success: false,
      errors: { confirmPassword: 'Passwords do not match' },
    }
  }

  try {
    await $api.patch('/auth/reset-password', {
      email,
      newPassword,
    })

    return { success: true, errors: null, timestamp: Date.now() }
  } catch (error: unknown) {
    serverLog('RESET_PASSWORD_ERROR', error, true)

    return {
      success: false,
      errors: parseAxiosError(error),
    }
  }
}
