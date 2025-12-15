'use client'
import { useState } from 'react'

import Input from '@/components/Input'

export default function SignUpForm() {
  const [email, setEmail] = useState<string>('')

  const handleEmailChange = (value: string) => {
    setEmail(value)
  }

  return (
    <form className="flex w-full flex-col">
      <h1 className="text-white text-xl font-bold text-center mb-4">Sign up</h1>
      <Input
        name="email"
        type="email"
        value={email}
        onChange={handleEmailChange}
        className="py-2.5 px-3 w-full"
        placeholder="Email"
      />
    </form>
  )
}
