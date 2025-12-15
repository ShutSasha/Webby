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
      className={`${className} focus:ring-emerald-500 border-border block rounded-lg border text-sm
        placeholder:text-[#FCFFFF]/16 focus:ring-[1px] focus:outline-none`}
      placeholder={placeholder}
      onChange={e => {
        onChange(e.target.value)
      }}
      defaultValue={defaultValue}
      type={type}
    />
  )
}
