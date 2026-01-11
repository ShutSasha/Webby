'use client'

import { Route } from 'next'
import { usePathname, useSearchParams, useRouter } from 'next/navigation'
import { useDebouncedCallback } from 'use-debounce'

import SearchIcon from '@/assets/icons/ic_search.svg'

type Props = {
  placeholder: string
}

export default function Search({ placeholder }: Props) {
  const searchParams = useSearchParams()
  const pathname = usePathname()
  const { replace } = useRouter()

  const handleSearch = useDebouncedCallback(term => {
    const params = new URLSearchParams(searchParams)
    params.set('page', '1')

    if (term) {
      params.set('query', term)
    } else {
      params.delete('query')
    }

    const url = `${pathname}?${params.toString()}`
    replace(url as Route)
  }, 300)

  return (
    <div className="relative flex">
      <label htmlFor="search" className="sr-only">
        Search
      </label>
      <input
        name="search"
        autoComplete="off"
        className="focus:ring-primary-800/60 block h-full w-full rounded-lg border-0 pl-9 text-sm
          placeholder:text-[#FCFFFF]/16 focus:ring-[1px] focus:outline-none"
        placeholder={placeholder}
        onChange={e => {
          handleSearch(e.target.value)
        }}
        defaultValue={searchParams.get('query')?.toString()}
      />
      <SearchIcon className="absolute top-1/2 left-3 h-4 w-4 -translate-y-1/2 text-neutral-50" />
    </div>
  )
}
