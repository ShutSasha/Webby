import { InfiniteData } from '@tanstack/react-query'

import { PaginatedData } from '@/types/general.types'

export const addPaginatedCacheItem = <T>(
  oldData: InfiniteData<PaginatedData<T>> | undefined,
  newItem: T,
): InfiniteData<PaginatedData<T>> | undefined => {
  if (!oldData?.pages || oldData.pages.length === 0) return oldData

  const newPages = [...oldData.pages]

  newPages[0] = {
    ...newPages[0],
    items: [newItem, ...(newPages[0].items || [])],
    totalCount: (newPages[0].totalCount || 0) + 1,
  }

  return {
    ...oldData,
    pages: newPages,
  }
}

export const removePaginatedCacheItem = <T>(
  oldData: InfiniteData<PaginatedData<T>> | undefined,
  itemId: string,
  idKey: keyof T = 'id' as keyof T,
): InfiniteData<PaginatedData<T>> | undefined => {
  if (!oldData?.pages) return oldData

  return {
    ...oldData,
    pages: oldData.pages.map(page => ({
      ...page,
      items: page.items?.filter(item => String(item[idKey]) !== String(itemId)) || [],
    })),
  }
}

export const updatePaginatedCacheItem = <T>(
  oldData: InfiniteData<PaginatedData<T>> | undefined,
  updatedItem: T,
  idKey: keyof T = 'id' as keyof T,
): InfiniteData<PaginatedData<T>> | undefined => {
  if (!oldData?.pages) return oldData

  return {
    ...oldData,
    pages: oldData.pages.map(page => ({
      ...page,
      items:
        page.items?.map(item =>
          String(item[idKey]) === String(updatedItem[idKey]) ? { ...item, ...updatedItem } : item,
        ) || [],
    })),
  }
}
