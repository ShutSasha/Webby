import { HTMLInputAutoCompleteAttribute, HTMLInputTypeAttribute } from 'react'

type Props = {
  name?: string
  type?: HTMLInputTypeAttribute
  value?: string
  className?: string
  defaultValue?: string
  placeholder?: string
  autoComplete?: HTMLInputAutoCompleteAttribute
  onChange: (value: string) => void
}

export default function Input({
  name,
  type,
  value,
  defaultValue,
  placeholder,
  className,
  autoComplete = 'off',
  onChange,
}: Props) {
  return (
    <input
      value={value}
      name={name}
      autoComplete={autoComplete}
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
