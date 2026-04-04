'use client'
import { ComponentPropsWithoutRef, FC, forwardRef, SVGProps } from 'react'

import EyeIcon from '@/assets/auth/ic_eye.svg'
import EyeOffIcon from '@/assets/auth/ic_eye_off.svg'

import Input from '../../shared/Input'

interface AuthInputProps extends ComponentPropsWithoutRef<'input'> {
  Icon: FC<SVGProps<SVGSVGElement>>
  showPassword?: boolean
  togglePassword?: () => void
}

const AuthInput = forwardRef<HTMLInputElement, AuthInputProps>(
  ({ Icon, showPassword, togglePassword, type, className, ...props }, ref) => {
    const isPassword = type === 'password'

    const currentType = isPassword ? (showPassword ? 'text' : 'password') : type

    return (
      <div className="relative">
        <Icon
          className="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-neutral-400 pointer-events-none"
        />

        <Input
          {...props}
          ref={ref}
          type={currentType}
          className={`py-2.5 pl-9 w-full ${isPassword ? 'pr-10' : 'pr-3'} ${className}`}
        />

        {isPassword && (
          <button
            type="button"
            onClick={togglePassword}
            className="absolute right-3 top-1/2 transform -translate-y-1/2 text-neutral-500 hover:text-emerald-500
              transition-colors cursor-pointer"
            tabIndex={-1}
          >
            {showPassword ? <EyeOffIcon className="w-5 h-5" /> : <EyeIcon className="w-5 h-5" />}
          </button>
        )}
      </div>
    )
  },
)

AuthInput.displayName = 'AuthInput'
export default AuthInput
