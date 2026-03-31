import { StudioVideoRow } from '@/ui/components/modules/Studio/StudioVideoRow'
import AnimatedTabs from '@/ui/components/shared/AnimatedTabs'

type Props = {
  searchParams: Promise<{ [key: string]: string | string[] | undefined }>
}

export default async function StudioPage({ searchParams }: Props) {
  const { tab } = await searchParams
  const currentTab = tab === 'playlists' ? 'playlists' : 'videos'

  const studioTabs = [
    { label: 'Videos', href: '?tab=videos', isActive: currentTab === 'videos' },
    { label: 'Playlists', href: '?tab=playlists', isActive: currentTab === 'playlists' },
  ]

  return (
    <>
      <h1 className="text-2xl font-bold text-neutral-100 tracking-tight">{`User's`} content</h1>

      <div className="border-b border-neutral-800/60 pb-4">
        <AnimatedTabs tabs={studioTabs} layoutId="studio-tabs" className="gap-4" />
      </div>

      <div className="flex flex-col flex-1 overflow-y-auto -mx-2 px-2">
        <div
          className="flex items-center text-xs font-medium text-neutral-400 border-b border-neutral-800 pb-3 px-4 sticky
            top-0 bg-neutral-900 z-10"
        >
          <div className="flex-1 min-w-[300px]">Video</div>
          <div className="w-28 text-center shrink-0">Visibility</div>
          <div className="w-32 text-center shrink-0">Date</div>
          <div className="w-24 text-center shrink-0">Views</div>
          <div className="w-12 shrink-0"></div>
        </div>

        <div className="flex flex-col">
          {[1, 2, 3, 4, 5].map(item => (
            <StudioVideoRow key={item} id={item.toString()} />
          ))}
        </div>
      </div>
    </>
  )
}
