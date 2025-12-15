import { HTMLInputTypeAttribute } from 'react'

type Props = {
  name?: string
  type?: HTMLInputTypeAttribute
  value?: string
  className?: string
  defaultValue?: string
  placeholder?: string
  onChange: (value: string) => void
}

export default function Input({ name, type, value, defaultValue, placeholder, className, onChange }: Props) {
  return (
    <input
      value={value}
      name={name}
      autoComplete="off"
      className={`${className} focus:border-emerald-500 focus:ring-emerald-500 ring-[0.3px] ring-transparent border
        border-border block rounded-lg text-sm placeholder:text-[#FCFFFF]/16 focus:outline-none`}
      placeholder={placeholder}
      onChange={e => {
        onChange(e.target.value)
      }}
      defaultValue={defaultValue}
      type={type}
    />
  )
}
