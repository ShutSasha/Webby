'use server'

import { AuthError } from 'next-auth'
import $api from '@/app/lib/api'
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
  const username = formData.get('username') as string
  const password = formData.get('password') as string
  const confirmPassword = formData.get('confirmPassword') as string

  if (password !== confirmPassword) {
    return { error: 'Паролі не співпадають' }
  }

  try {
    await $api.post('/auth/registration', {
      username,
      password,
      confirmPassword,
    })

    return { success: true }
  } catch (error: any) {
    console.error('Registration error:', error)

    const msg = error.response?.data?.message || 'Помилка реєстрації'
    return { error: Array.isArray(msg) ? msg.join(', ') : msg }
  }
}
