'use server'

import { Route } from 'next'
import { redirect } from 'next/dist/client/components/navigation'
import { AuthError } from 'next-auth'

import $api from '@/app/api'
import { serverLog } from '@/lib/utils/utils'

import { signIn } from '../../../auth'

export async function authenticate(prevState: string | undefined, formData: FormData) {
  try {
    await signIn('credentials', formData)
  } catch (error) {
    if (error instanceof AuthError) {
      switch (error.type) {
        case 'CredentialsSignin':
          return 'Невірні дані входу (Username або Password).'
        default:
          return 'Щось пішло не так.'
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
      errors: { repeatPassword: 'Паролі не співпадають' },
    }
  }

  try {
    const result = await $api.post('/auth/sign-up', {
      email,
      username,
      password,
    })

    console.log('Registration result:', result.data)
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
