'use client'

import { Route } from 'next'
import { usePathname, useSearchParams, useRouter } from 'next/navigation'
import { useDebouncedCallback } from 'use-debounce'

import SearchIcon from '@/assets/icons/ic_search.svg'

type Props = {
  placeholder: string
  inputClassName?: string
  iconClassName?: string
  containerClassName?: string
}

const DEFAULT_ICON_CLASSNAME = 'left-4 h-6 w-6 text-foreground-disabled'
const DEFAULT_INPUT_CLASSNAME = 'pl-12 py-3 rounded-2xl font-medium text-[16px] leading-5 border border-border w-full '

export default function Search({
  placeholder,
  inputClassName = DEFAULT_INPUT_CLASSNAME,
  iconClassName = DEFAULT_ICON_CLASSNAME,
  containerClassName = '',
}: Props) {
  const searchParams = useSearchParams()
  const pathname = usePathname()
  const { replace } = useRouter()

  const handleSearch = useDebouncedCallback(term => {
    const params = new URLSearchParams(searchParams)

    if (term) {
      params.set('query', term)
    } else {
      params.delete('query')
    }

    const url = `${pathname}?${params.toString()}`
    replace(url as Route)
  }, 300)

  return (
    <div className={`relative ${containerClassName}`}>
      <SearchIcon className={`${iconClassName} absolute top-1/2 -translate-y-1/2 stroke-[1.5px]`} aria-hidden="true" />
      <label htmlFor="search" className="sr-only">
        Search
      </label>
      <input
        id="search"
        name="search"
        autoComplete="off"
        className={`${inputClassName} focus:border-emerald-500 focus:ring-emerald-500 ring-[0.1px] ring-transparent
          block placeholder:text-foreground-disabled focus:outline-none `}
        placeholder={placeholder}
        onChange={e => {
          handleSearch(e.target.value)
        }}
        defaultValue={searchParams.get('query')?.toString()}
      />
    </div>
  )
}
