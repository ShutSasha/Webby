import { FC, HTMLInputTypeAttribute, SVGProps } from 'react'

import Input from './Input'

type Props = {
  Icon: FC<SVGProps<SVGSVGElement>>
  name?: string
  type?: HTMLInputTypeAttribute
  value?: string
  className?: string
  defaultValue?: string
  placeholder?: string
  onChange: (value: string) => void
}

export default function AuthInput({ Icon, name, type, value, className, defaultValue, placeholder, onChange }: Props) {
  return (
    <div className="relative">
      <Icon className={'absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4'} />
      <Input
        name={name}
        type={type}
        value={value}
        onChange={onChange}
        className={`py-2.5 pl-9 w-full ${className}`}
        placeholder={placeholder}
        defaultValue={defaultValue}
      />
    </div>
  )
}
