type PlaylistQueueFooterProps = {
  fetchingMore: boolean
  hasMore: boolean
  videosLength: number
  loading: boolean
  appliedQuery: string
}

export default function PlaylistQueueFooter({
  fetchingMore,
  hasMore,
  videosLength,
  loading,
  appliedQuery,
}: PlaylistQueueFooterProps) {
  return (
    <>
      {fetchingMore && (
        <div className="flex justify-center py-4">
          <div className="size-5 border-2 border-emerald-500/20 border-t-emerald-500 rounded-full animate-spin" />
        </div>
      )}

      {!hasMore && videosLength > 0 && (
        <p className="text-center text-xs text-foreground-disabled py-4 italic">End of playlist</p>
      )}

      {!loading && videosLength === 0 && (
        <p className="text-center text-foreground-faint py-10">
          {appliedQuery.trim() !== '' ? 'No videos found' : 'Playlist is empty'}
        </p>
      )}
    </>
  )
}
