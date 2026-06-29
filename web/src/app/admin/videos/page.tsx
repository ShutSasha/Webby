import ModerationVideosTable from '@/ui/components/modules/Admin/videos/ModerationVideosTable'
import Search from '@/ui/components/shared/Search'

type Props = {
  searchParams: Promise<{ query?: string }>
}

export default async function ModerationVideosPage({ searchParams }: Props) {
  const params = await searchParams
  const query = params?.query || ''

  return (
    <div className="flex flex-col gap-6 animate-in fade-in duration-500 pb-10">
      <div className="flex items-center gap-4 w-full">
        <Search
          placeholder="Search videos by title..."
          containerClassName="flex-1"
          inputClassName="w-full bg-surface border border-border rounded-xl pl-12 pr-4 py-3 outline-none focus:border-neutral-600
          text-foreground-tertiary placeholder:text-foreground-disabled transition-colors"
          iconClassName="left-4 size-5 text-foreground-faint"
          delayMs={400}
        />
      </div>

      <ModerationVideosTable searchQuery={query} />
    </div>
  )
}
