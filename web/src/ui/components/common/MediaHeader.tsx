import MediaButton from './MediaButton'
import PageToggle from '../PageToggle'
import Search from '../Search'

type Props = {
  searchPlaceholder: string
  actionLabel: string
  onActionClick?: () => void
}

export default function MediaHeader({ searchPlaceholder, actionLabel, onActionClick }: Props) {
  return (
    <div className="flex flex-wrap lg:flex-nowrap items-center justify-between gap-4">
      <PageToggle />

      <Search
        placeholder={searchPlaceholder}
        containerClassName="order-1 lg:order-2 w-full lg:flex-1 lg:max-w-[534px]"
      />

      {/* <MediaButton actionLabel={actionLabel} onActionClick={onActionClick} /> */}
      <MediaButton actionLabel={actionLabel} onActionClick={onActionClick}>
        <p className="text-neutral-300">text</p>
      </MediaButton>
    </div>
  )
}
