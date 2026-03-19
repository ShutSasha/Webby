'use client'

import { useState, useTransition, useRef, useEffect } from 'react'

import { useForm } from 'react-hook-form'

import { cn } from '@/lib/utils/utils'
import { useToastStore } from '@/stores/toast-store'
import Button from '@/ui/components/shared/Button'

type RoomSettingsFields = {
  roomName: string
  roomType: 'public' | 'private'
}

const options: { value: 'public' | 'private'; label: string }[] = [
  { value: 'public', label: 'Public Room' },
  { value: 'private', label: 'Private Room' },
]

export default function RoomSettingsContainer() {
  const [isPending, startTransition] = useTransition()
  const addToast = useToastStore(state => state.addToast)

  const [isOpen, setIsOpen] = useState(false)
  const dropdownRef = useRef<HTMLDivElement>(null)

  const {
    register,
    handleSubmit,
    setValue,
    watch,
    formState: { errors },
  } = useForm<RoomSettingsFields>({
    defaultValues: {
      roomName: '',
      roomType: 'public',
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
    startTransition(async () => {
      await new Promise(resolve => setTimeout(resolve, 1000))
      addToast('Room settings updated successfully', 'success')
    })
  }

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col h-full gap-4">
      <div className="flex flex-col gap-1.5">
        <label htmlFor="roomName" className="text-sm text-neutral-500 ml-1">
          Room name:
        </label>
        <input
          {...register('roomName', { required: 'Room name is required' })}
          id="roomName"
          type="text"
          placeholder="Enter Room name"
          className="bg-neutral-900 border border-transparent focus:border-emerald-500/70 placeholder:text-neutral-700
            py-2 px-3 rounded-lg outline-0 text-neutral-200 w-full transition-all duration-300"
        />
        {errors.roomName && <p className="text-red-500 text-xs ml-1">{errors.roomName.message}</p>}
      </div>

      <div className="flex flex-col gap-1.5" ref={dropdownRef}>
        <label className="text-sm text-neutral-500 ml-1">Room type:</label>

        <div className="relative select-none">
          <div
            onClick={() => setIsOpen(!isOpen)}
            className={cn(
              `bg-neutral-900 border border-transparent py-2 px-3 rounded-lg text-neutral-500 w-full cursor-pointer
              transition-all duration-300 flex justify-between items-center`,
              isOpen && 'border-emerald-500/70',
            )}
          >
            <span>{options.find(opt => opt.value === currentRoomType)?.label}</span>
            <svg
              className={cn('w-4 h-4 text-neutral-700 transition-transform', isOpen && 'rotate-180')}
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
                flex flex-col gap-1 shadow-2xl"
            >
              {options.map(option => (
                <div
                  key={option.value}
                  onClick={() => {
                    setValue('roomType', option.value)
                    setIsOpen(false)
                  }}
                  className={cn(
                    'px-4 py-1 rounded-lg transition-colors cursor-pointer text-neutral-300 text-[16px]',
                    currentRoomType === option.value ? 'bg-neutral-800' : 'hover:bg-neutral-800/50',
                  )}
                >
                  {option.label}
                </div>
              ))}
            </div>
          )}
        </div>
      </div>

      <div className="pt-4 mt-auto flex justify-center">
        <Button
          type="submit"
          viewType={!isPending ? 'confirm' : 'loading'}
          disabled={isPending}
          className="rounded-xl px-6 py-2 font-semibold"
        >
          {isPending ? 'Saving...' : 'Save'}
        </Button>
      </div>
    </form>
  )
}
