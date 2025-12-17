'use client'
import { FC, HTMLInputAutoCompleteAttribute, HTMLInputTypeAttribute, SVGProps, useState } from 'react'

import EyeIcon from '@/assets/auth/ic_eye.svg'
import EyeOffIcon from '@/assets/auth/ic_eye_off.svg'

import Input from './Input'

type Props = {
  Icon: FC<SVGProps<SVGSVGElement>>
  name?: string
  type?: HTMLInputTypeAttribute
  value?: string
  className?: string
  defaultValue?: string
  placeholder?: string
  autoComplete?: HTMLInputAutoCompleteAttribute
  onChange: (value: string) => void
}

export default function AuthInput({
  Icon,
  name,
  type,
  value,
  className,
  defaultValue,
  placeholder,
  autoComplete,
  onChange,
}: Props) {
  const [showPassword, setShowPassword] = useState<boolean>(false)

  const isPasswordType = type === 'password'

  const currentType = isPasswordType ? (showPassword ? 'text' : 'password') : type

  const togglePasswordVisibility = () => {
    setShowPassword(prev => !prev)
  }

  return (
    <div className="relative">
      <Icon className={'absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-neutral-400'} />

      <Input
        name={name}
        type={currentType}
        value={value}
        onChange={onChange}
        className={`py-2.5 pl-9 pr-10 w-full ${className}`}
        placeholder={placeholder}
        defaultValue={defaultValue}
        autoComplete={autoComplete}
      />

      {isPasswordType && (
        <button
          type="button"
          onClick={togglePasswordVisibility}
          className="absolute right-3 top-1/2 transform -translate-y-1/2 text-neutral-500 hover:text-emerald-500
            transition-colors cursor-pointer"
        >
          {showPassword ? <EyeOffIcon className="w-5 h-5" /> : <EyeIcon className="w-5 h-5" />}
        </button>
      )}
    </div>
  )
}
