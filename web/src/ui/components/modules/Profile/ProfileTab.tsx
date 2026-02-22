'use client'

import { Route } from 'next'
import { usePathname, useSearchParams, useRouter } from 'next/navigation'

type TabCategory = 'Video' | 'Playlist'

type Props = {
  label: TabCategory
}

export default function ProfileTab({ label }: Props) {
  const searchParams = useSearchParams()
  const pathname = usePathname()
  const { replace } = useRouter()

  const currentTab = searchParams.get('tab') ?? 'Video' 
  const isActive = label === currentTab

  const handleTab = (tab: TabCategory) => {
    if (isActive) return

    const params = new URLSearchParams(searchParams)
    params.set('tab', tab)

    replace(`${pathname}?${params.toString()}` as Route, { scroll: false })
  }

  return (
    <div className="flex group" onClick={() => handleTab(label)}>
      <div
        className={`flex items-center py-1 px-7 border border-border text-sm rounded-full cursor-pointer transition-all
          duration-100 ease-in ${isActive ? 'bg-emerald-400' : 'group-hover:bg-neutral-700/40'}`}
      >
        <p className={`transition-all duration-300 font-medium ${isActive ? 'text-neutral-900' : ''}`}>{label}</p>
      </div>
    </div>
  )
}
