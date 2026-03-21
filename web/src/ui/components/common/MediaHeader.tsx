import { Suspense } from 'react'

import Search from '../Search'

type Props = {
  searchPlaceholder: string
  pageToggle: React.ReactNode
  actionSlot?: React.ReactNode
}

export default function MediaHeader({ searchPlaceholder, pageToggle, actionSlot }: Props) {
  return (
    <div className="flex flex-wrap lg:flex-nowrap items-center justify-between gap-4">
      {pageToggle}
      <Suspense>
        <Search
          placeholder={searchPlaceholder}
          containerClassName="order-1 lg:order-2 w-full lg:flex-1 lg:max-w-[534px]"
        />
      </Suspense>

      {/* Button here */}
      {actionSlot}
    </div>
  )
}
