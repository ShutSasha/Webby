'use client'

import { useMemo } from 'react'

import { useGetComplaintsQuery } from '@/lib/hooks/api/admin/useGetComplaints'
import { useInfiniteScroll } from '@/lib/hooks/useInfiniteScroll'
import ComplaintCard from '@/ui/components/modules/Admin/ComplaintCard'

export default function AdminReportsPage() {
  const { data, isLoading, isFetchingNextPage, hasNextPage, fetchNextPage } = useGetComplaintsQuery(10)

  const complaints = useMemo(() => {
    return data?.pages.flatMap(page => page?.items || []) || []
  }, [data])

  const lastElementRef = useInfiniteScroll({
    isLoading,
    isFetchingNextPage,
    hasNextPage,
    fetchNextPage,
  })

  return (
    <div className="flex flex-col gap-4 animate-in fade-in duration-500 pb-10">
      {isLoading && complaints.length === 0 ? (
        <div className="flex justify-center py-20">
          <div className="size-10 border-4 border-emerald-500/20 border-t-emerald-500 rounded-full animate-spin" />
        </div>
      ) : complaints.length === 0 ? (
        <div className="flex flex-col items-center justify-center py-20 text-foreground-faint">
          <p>No active complaints to review right now. Great job!</p>
        </div>
      ) : (
        <>
          {complaints.map((complaint, index) => {
            const isLast = complaints.length === index + 1

            const card = <ComplaintCard key={complaint.id} complaint={complaint} />

            if (isLast) {
              return (
                <div ref={lastElementRef} key={`last-${complaint.id}`}>
                  {card}
                </div>
              )
            }

            return card
          })}
        </>
      )}

      {isFetchingNextPage && (
        <div className="flex justify-center py-6">
          <div className="size-6 border-4 border-emerald-500/20 border-t-emerald-500 rounded-full animate-spin" />
        </div>
      )}
    </div>
  )
}
