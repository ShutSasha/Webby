import { useCallback, useEffect, useRef, useState } from 'react'

import { useDebouncedCallback } from 'use-debounce'

export type InfiniteSearchResponse<T> = {
  success: boolean
  data: { items: T[] } | null | any
}

type UseInfiniteSearchProps<T> = {
  fetchFn: (query: string, page: number, pageSize: number) => Promise<InfiniteSearchResponse<T>>
  pageSize?: number
  enabled?: boolean
  debounceMs?: number
}

export function useInfiniteSearch<T>({
  fetchFn,
  pageSize = 20,
  enabled = true,
  debounceMs = 400,
}: UseInfiniteSearchProps<T>) {
  const [items, setItems] = useState<T[]>([])
  const [loading, setLoading] = useState(true)
  const [fetchingMore, setFetchingMore] = useState(false)
  const [page, setPage] = useState(1)
  const [hasMore, setHasMore] = useState(true)
  const [appliedQuery, setAppliedQuery] = useState('')

  const observer = useRef<IntersectionObserver | null>(null)

  const lastElementRef = useCallback(
    (node: HTMLDivElement) => {
      if (loading || fetchingMore) return
      if (observer.current) observer.current.disconnect()

      observer.current = new IntersectionObserver(entries => {
        if (entries[0].isIntersecting && hasMore) {
          setPage(prevPage => prevPage + 1)
        }
      })

      if (node) observer.current.observe(node)
    },
    [loading, fetchingMore, hasMore],
  )

  const loadData = useCallback(
    async (searchQuery: string, targetPage: number, isInitial: boolean) => {
      if (!enabled) return

      if (isInitial) setLoading(true)
      else setFetchingMore(true)

      try {
        const response = await fetchFn(searchQuery, targetPage, pageSize)

        if (response?.success && response.data) {
          const newItems = response.data.items
          setItems(prev => (isInitial ? newItems : [...prev, ...newItems]))
          setHasMore(newItems.length === pageSize)
        }
      } catch (error) {
        console.error('INFINITE_SEARCH_HOOK_ERROR', error)
      } finally {
        setLoading(false)
        setFetchingMore(false)
      }
    },
    [fetchFn, pageSize, enabled],
  )

  const debouncedSearch = useDebouncedCallback((value: string) => {
    setPage(1)
    setItems([])
    setHasMore(true)
    setAppliedQuery(value)
    loadData(value, 1, true)
  }, debounceMs)

  const handleSearchChange = (value: string) => {
    debouncedSearch(value)
  }

  const reset = useCallback(() => {
    setItems([])
    setPage(1)
    setAppliedQuery('')
    setHasMore(true)
    setLoading(false)
  }, [])

  useEffect(() => {
    if (enabled && items.length === 0 && appliedQuery === '') {
      loadData('', 1, true)
    }
  }, [enabled, loadData, items.length, appliedQuery])

  useEffect(() => {
    if (page > 1 && enabled) {
      loadData(appliedQuery, page, false)
    }
  }, [page, appliedQuery, loadData, enabled])

  return {
    items,
    setItems,
    loading,
    fetchingMore,
    hasMore,
    appliedQuery,
    lastElementRef,
    handleSearchChange,
    reset,
  }
}
