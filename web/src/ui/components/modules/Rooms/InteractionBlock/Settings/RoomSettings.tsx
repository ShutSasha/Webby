'use client'

import { useState, useRef, useEffect } from 'react'

import { useParams } from 'next/navigation'
import { useForm } from 'react-hook-form'

import { useGetRoomByIdQuery } from '@/lib/hooks/api/room/useGetRoomByIdQuery'
import { useUpdateRoom } from '@/lib/hooks/api/room/useUpdateRoom'
import { cn } from '@/lib/utils/general.utils'
import Button from '@/ui/components/shared/Button'

import RoomSettingsSkeleton from './RoomSettingsSkeleton'

type RoomSettingsFields = {
  roomName: string
  roomType: 'public' | 'private'
}

const options: { value: 'public' | 'private'; label: string }[] = [
  { value: 'public', label: 'Public Room' },
  { value: 'private', label: 'Private Room' },
]

export default function RoomSettings() {
  const params = useParams()
  const roomId = params?.id as string

  const [isOpen, setIsOpen] = useState(false)
  const dropdownRef = useRef<HTMLDivElement>(null)

  const { data: room, isLoading: isLoadingRoom } = useGetRoomByIdQuery(roomId)
  const { mutate: updateRoom, isPending } = useUpdateRoom()

  const {
    register,
    handleSubmit,
    setValue,
    watch,
    formState: { errors },
  } = useForm<RoomSettingsFields>({
    values: {
      roomName: room?.name || '',
      roomType: room?.isPrivate ? 'private' : 'public',
    },
  })

  const currentRoomType = watch('roomType')

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
        setIsOpen(false)
      }
    }
    document.addEventListener('mousedown', handleClickOutside)
    return () => document.removeEventListener('mousedown', handleClickOutside)
  }, [])

  const onSubmit = (data: RoomSettingsFields) => {
    if (!roomId) return

    const formData = new FormData()
    formData.append('name', data.roomName)
    formData.append('isPrivate', data.roomType === 'private' ? 'true' : 'false')

    updateRoom({ roomId, formData })
  }

  if (isLoadingRoom) {
    return <RoomSettingsSkeleton />
  }

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col h-full px-1 animate-in fade-in duration-300">
      <div className="flex flex-col gap-5">
        <h3 className="text-[11px] font-bold text-neutral-500 uppercase tracking-wider pl-1">Room Details</h3>

        <div className="flex flex-col gap-1.5">
          <label htmlFor="roomName" className="text-[13px] font-medium text-neutral-300 pl-1">
            Room name
          </label>
          <input
            {...register('roomName', { required: 'Room name is required' })}
            id="roomName"
            type="text"
            placeholder="Enter Room name"
            className="bg-neutral-900/60 border border-neutral-800/80 focus:border-emerald-500/50
              hover:border-neutral-700 placeholder:text-neutral-600 py-2.5 px-3.5 rounded-xl outline-none
              text-neutral-200 w-full transition-all text-sm shadow-sm shadow-black/20"
          />
          {errors.roomName && <p className="text-red-500 text-xs mt-0.5 pl-1">{errors.roomName.message}</p>}
        </div>

        <div className="flex flex-col gap-1.5" ref={dropdownRef}>
          <label className="text-[13px] font-medium text-neutral-300 pl-1">Room privacy</label>

          <div className="relative select-none">
            <div
              onClick={() => setIsOpen(!isOpen)}
              className={cn(
                `bg-neutral-900/60 border border-neutral-800/80 py-2.5 px-3.5 rounded-xl text-neutral-300 w-full
                cursor-pointer transition-all flex justify-between items-center hover:border-neutral-700 text-sm
                shadow-sm shadow-black/20`,
                isOpen && 'border-emerald-500/50 bg-neutral-900',
              )}
            >
              <span>{options.find(opt => opt.value === currentRoomType)?.label}</span>
              <svg
                className={cn('w-4 h-4 text-neutral-500 transition-transform duration-200', isOpen && 'rotate-180')}
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M19 9l-7 7-7-7" />
              </svg>
            </div>

            {isOpen && (
              <div
                className="absolute left-0 right-0 mt-2 p-1.5 bg-neutral-900 border border-neutral-800 rounded-xl z-50
                  flex flex-col gap-1 shadow-xl animate-in fade-in zoom-in-95 duration-200"
              >
                {options.map(option => (
                  <div
                    key={option.value}
                    onClick={() => {
                      setValue('roomType', option.value)
                      setIsOpen(false)
                    }}
                    className={cn(
                      'px-3 py-2 rounded-lg transition-colors cursor-pointer text-sm font-medium',
                      currentRoomType === option.value
                        ? 'bg-neutral-800 text-neutral-100'
                        : 'text-neutral-400 hover:bg-neutral-800/50 hover:text-neutral-200',
                    )}
                  >
                    {option.label}
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>
      </div>

      <div className="mt-8 pt-5 border-t border-neutral-800/50">
        <Button
          type="submit"
          viewType={!isPending ? 'confirm' : 'loading'}
          disabled={isPending}
          className="rounded-xl py-2.5 font-semibold w-full text-sm shadow-lg shadow-emerald-500/10
            hover:shadow-emerald-500/20"
        >
          {isPending ? 'Saving...' : 'Save changes'}
        </Button>
      </div>
    </form>
  )
}
