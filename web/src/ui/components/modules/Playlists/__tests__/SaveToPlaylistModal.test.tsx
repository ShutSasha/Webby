import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import '@testing-library/jest-dom/vitest'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import { useBulkTogglePlaylistMediaMutation } from '@/lib/hooks/api/playlist/useBulkTogglePlaylistMedia'
import { useSearchUserPlaylistsQuery } from '@/lib/hooks/api/playlist/useSearchUserPlaylists'
import { useToastStore } from '@/stores/toast-store'

import SaveToPlaylistModal from '../SaveToPlaylistModal'

vi.mock('@/lib/hooks/api/playlist/useSearchUserPlaylists', () => ({
  useSearchUserPlaylistsQuery: vi.fn(),
}))

vi.mock('@/lib/hooks/api/playlist/useBulkTogglePlaylistMedia', () => ({
  useBulkTogglePlaylistMediaMutation: vi.fn(),
}))

vi.mock('@/lib/hooks/useInfiniteScroll', () => ({
  useInfiniteScroll: () => vi.fn(),
}))

vi.mock('@/stores/toast-store', () => ({
  useToastStore: vi.fn(),
}))

const mockPlaylists = [
  {
    playlistId: 'playlist-1',
    name: 'Chill Vibes',
    countOfVideos: 5,
    playlistCover: '/mock-cover.jpg',
    isVideoAdded: false,
  },
  {
    playlistId: 'playlist-2',
    name: 'Rock Classics',
    countOfVideos: 10,
    playlistCover: '/mock-cover2.jpg',
    isVideoAdded: true,
  },
]

describe('SaveToPlaylistModal Integration', () => {
  const mockOnClose = vi.fn()
  const mockSaveBulkChanges = vi.fn()
  const mockAddToast = vi.fn()

  beforeEach(() => {
    vi.clearAllMocks()
    ;(useToastStore as unknown as ReturnType<typeof vi.fn>).mockImplementation((selector: any) => {
      return selector({ addToast: mockAddToast })
    })
    ;(useSearchUserPlaylistsQuery as unknown as ReturnType<typeof vi.fn>).mockReturnValue({
      data: { pages: [{ items: mockPlaylists }] },
      isLoading: false,
      isFetchingNextPage: false,
      hasNextPage: false,
      fetchNextPage: vi.fn(),
    })
    ;(useBulkTogglePlaylistMediaMutation as unknown as ReturnType<typeof vi.fn>).mockReturnValue({
      mutate: mockSaveBulkChanges,
      isPending: false,
    })
  })

  it('renders correctly and displays fetched playlists', () => {
    render(
      <SaveToPlaylistModal isOpen={true} onClose={mockOnClose} videoId="video-123" userId="user-1" mediaType="Video" />,
    )

    expect(screen.getByText('Chill Vibes')).toBeInTheDocument()
    expect(screen.getByText('Rock Classics')).toBeInTheDocument()

    expect(screen.getByText('5 videos')).toBeInTheDocument()
    expect(screen.getByText('10 videos')).toBeInTheDocument()

    const saveBtn = screen.getByRole('button', { name: /save/i })
    expect(saveBtn).toBeDisabled()
  })

  it('toggles local selection and updates video count dynamically', async () => {
    const user = userEvent.setup()
    render(
      <SaveToPlaylistModal isOpen={true} onClose={mockOnClose} videoId="video-123" userId="user-1" mediaType="Video" />,
    )

    const chillVibesItem = screen.getByText('Chill Vibes').closest('div')
    await user.click(chillVibesItem!)

    expect(screen.getByText('6 videos')).toBeInTheDocument()

    const saveBtn = screen.getByRole('button', { name: /save/i })
    expect(saveBtn).not.toBeDisabled()

    await user.click(chillVibesItem!)

    expect(screen.getByText('5 videos')).toBeInTheDocument()
    expect(saveBtn).toBeDisabled()
  })

  it('calls save mutation with correct data and triggers success toast', async () => {
    const user = userEvent.setup()

    mockSaveBulkChanges.mockImplementation((_: any, options: any) => {
      options.onSuccess()
    })

    render(
      <SaveToPlaylistModal isOpen={true} onClose={mockOnClose} videoId="video-123" userId="user-1" mediaType="Video" />,
    )

    await user.click(screen.getByText('Chill Vibes'))
    await user.click(screen.getByText('Rock Classics'))

    const saveBtn = screen.getByRole('button', { name: /save/i })
    await user.click(saveBtn)

    expect(mockSaveBulkChanges).toHaveBeenCalledWith(
      {
        playlistIds: ['playlist-1', 'playlist-2'],
        mediaId: 'video-123',
        mediaType: 'Video',
      },
      expect.any(Object),
    )

    expect(mockAddToast).toHaveBeenCalledWith('Playlists updated successfully', 'success')
  })

  it('shows loading state correctly', () => {
    ;(useSearchUserPlaylistsQuery as unknown as ReturnType<typeof vi.fn>).mockReturnValue({
      data: undefined,
      isLoading: true,
      isFetchingNextPage: false,
      playlists: [],
    })

    render(
      <SaveToPlaylistModal isOpen={true} onClose={mockOnClose} videoId="video-123" userId="user-1" mediaType="Video" />,
    )

    expect(screen.queryByText('Chill Vibes')).not.toBeInTheDocument()
  })

  it('resets local state when modal is closed and reopened', async () => {
    const user = userEvent.setup()

    const { rerender } = render(
      <SaveToPlaylistModal isOpen={true} onClose={mockOnClose} videoId="video-123" userId="user-1" mediaType="Video" />,
    )

    const chillVibesItem = screen.getByText('Chill Vibes').closest('div')
    await user.click(chillVibesItem!)

    const saveBtn = screen.getByRole('button', { name: /save/i })
    expect(saveBtn).not.toBeDisabled()

    rerender(
      <SaveToPlaylistModal
        isOpen={false}
        onClose={mockOnClose}
        videoId="video-123"
        userId="user-1"
        mediaType="Video"
      />,
    )

    rerender(
      <SaveToPlaylistModal isOpen={true} onClose={mockOnClose} videoId="video-123" userId="user-1" mediaType="Video" />,
    )

    const newSaveBtn = screen.getByRole('button', { name: /save/i })
    expect(newSaveBtn).toBeDisabled()

    expect(screen.getByText('5 videos')).toBeInTheDocument()
  })
})
