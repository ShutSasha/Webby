import { ComponentPropsWithoutRef, forwardRef } from 'react'

type InputProps = ComponentPropsWithoutRef<'input'>

const Input = forwardRef<HTMLInputElement, InputProps>(({ className, autoComplete = 'off', ...props }, ref) => {
  return (
    <input
      ref={ref}
      autoComplete={autoComplete}
      className={` ${className} focus:border-emerald-500 focus:ring-emerald-500 ring-[0.3px] ring-transparent border
        border-border block rounded-lg text-sm placeholder:text-foreground-strong/20 focus:outline-none `}
      {...props}
    />
  )
})

Input.displayName = 'Input'
export default Input
