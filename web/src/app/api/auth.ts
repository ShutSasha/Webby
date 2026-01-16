'use server'

import { Route } from 'next'
import { redirect } from 'next/dist/client/components/navigation'
import { AuthError } from 'next-auth'

import $api from '@/app/api'

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
    // await $api.post('/auth/registration', {
    //   email,
    //   username,
    //   password,
    // })
  } catch (error: any) {
    console.error('Registration error:', error)

    const serverErrors = error.response?.data?.errors || {}
    const message = error.response?.data?.message || 'Помилка реєстрації'

    return {
      success: false,
      message,
      errors: serverErrors,
    }
  }

  redirect(`/sign-up/email-verify?email=${encodeURIComponent(email)}` satisfies Route)
}
