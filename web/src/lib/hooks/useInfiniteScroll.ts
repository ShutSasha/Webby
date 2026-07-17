import { useCallback, useRef } from 'react'

type UseInfiniteScrollProps = {
  isLoading: boolean
  isFetchingNextPage: boolean
  hasNextPage: boolean | undefined
  fetchNextPage: () => void
}

export const useInfiniteScroll = ({
  isLoading,
  isFetchingNextPage,
  hasNextPage,
  fetchNextPage,
}: UseInfiniteScrollProps) => {
  const observer = useRef<IntersectionObserver | null>(null)

  const lastElementRef = useCallback(
    (node: HTMLElement | null) => {
      if (isLoading || isFetchingNextPage) return

      if (observer.current) {
        observer.current.disconnect()
      }

      observer.current = new IntersectionObserver(entries => {
        if (entries[0].isIntersecting && hasNextPage) {
          fetchNextPage()
        }
      })

      if (node) {
        observer.current.observe(node)
      }
    },
    [isLoading, isFetchingNextPage, hasNextPage, fetchNextPage],
  )

  return lastElementRef
}
