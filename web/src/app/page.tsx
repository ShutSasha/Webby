'use client'

import Image from 'next/image'
import Link from 'next/link'

import ImageBackgroundGroup from '@/assets/auth/bg-login.png'
import ImageBackgroundMan from '@/assets/auth/home-page-man-sign-up.png'
import { benefitsList } from '@/lib/placeholder-data/home-page'
import { clog, serverLog } from '@/lib/utils/utils'
import { useToastStore } from '@/stores/toast-store'
import MainLayout from '@/ui/components/MainLayout'
import Button from '@/ui/components/shared/Button'

import $api from './api'

// TODO: return server component
export default function Home() {
  const addToast = useToastStore(state => state.addToast)

  const handleTest = async () => {
    try {
      const { data } = await $api.get('/auth/check')

      clog('test data response', data)
      addToast('Successss', 'success')
    } catch (error) {
      serverLog('TEST ERROR', error, true)
    }
  }

  return (
    <MainLayout>
      <div className="flex flex-col gap-4 w-full max-w-5xl 2xl:max-w-7xl mx-auto">
        <div className="bg-neutral-900 rounded-[20px] p-5 flex justify-between gap-4">
          <div className="flex flex-col justify-between gap-4">
            <div>
              <h1 className="font-semibold text-xl sm:text-3xl 2xl:text-4xl leading-tight mb-2.5">
                Watch videos together — anytime, anywhere
              </h1>
              <h2 className="font-medium text-lg sm:text-xl 2xl:text-2xl leading-tight text-emerald-500">
                No downloads
              </h2>
            </div>
            <div className="flex gap-3">
              <Button
                className="font-semibold rounded-[20px] text-sm sm:text-[16px]"
                paddingClasses="px-3 py-2 sm:px-4.5 sm:py-2.5 2xl:px-5 2xl:py-3"
                viewType="confirm"
              >
                Create a room
              </Button>
              <Link href="/rooms">
                <Button
                  className="font-semibold rounded-[20px] text-sm sm:text-[16px]"
                  paddingClasses="px-3 py-2 sm:px-4.5 sm:py-2.5 2xl:px-5 2xl:py-3"
                  viewType="cancel"
                >
                  Find a room
                </Button>
              </Link>
            </div>
          </div>
          <Image
            src={ImageBackgroundMan}
            alt=""
            loading="eager"
            className="h-full hidden sm:block w-40 2xl:w-50 rounded-lg"
          />
        </div>
        {/* Block - React to moments together */}
        <button
          onClick={handleTest}
          className="text-[36px] text-red-500 font-black px-4 py-2 rounded-full border border-border hover:text-white
            hover:bg-red-600 transition-all duration-300 cursor-pointer"
        >
          TEST BUTTON
        </button>
        <div className="bg-neutral-900 rounded-[20px] p-5">
          <ul className="flex flex-wrap justify-center sm:flex-row sm:flex-nowrap sm:justify-between mb-5 gap-2">
            {benefitsList.map((benefit, index) => (
              <div key={index} className="flex flex-col items-center gap-2">
                <benefit.icon className="size-10 lg:size-12.5 2xl:size-14.5 text-neutral-100" />
                <h3 className="font-bold text-[16px] uppercase text-center">{benefit.title}</h3>
                <p className="text-sm text-neutral-500 max-w-[200px] text-center">{benefit.description}</p>
              </div>
            ))}
          </ul>
          <div className="flex justify-between gap-4">
            <div className="flex flex-col justify-between gap-4">
              <div className="max-w-[540px]">
                <h2 className="font-semibold text-xl sm:text-2xl 2xl:text-3xl leading-tight mb-2.5">
                  React to moments together 🔥
                </h2>
                <p className="font-medium text-sm sm:text-lg 2xl:text-xl leading-tight text-neutral-500">
                  Find moments of shared joy, even when you’re apart.
                </p>
              </div>
              <div className="flex flex-col gap-4">
                <p className="font-medium text-sm 2xl:text-lg leading-tight text-neutral-100">
                  Ready to join the journey? Sign up or log in below to get started!
                </p>
                <div className="flex gap-3">
                  <Link href="/sign-up">
                    <Button
                      className="font-semibold rounded-[20px] text-sm sm:text-[16px]"
                      paddingClasses="px-3 py-2 sm:px-4.5 sm:py-2.5 2xl:px-5 2xl:py-3"
                      viewType="confirm"
                    >
                      Sign up
                    </Button>
                  </Link>
                  <Link href="/login">
                    <Button
                      className="font-semibold rounded-[20px] text-sm sm:text-[16px]"
                      paddingClasses="px-3 py-2 sm:px-4.5 sm:py-2.5 2xl:px-5 2xl:py-3"
                      viewType="cancel"
                    >
                      Login
                    </Button>
                  </Link>
                </div>
              </div>
            </div>
            <Image
              src={ImageBackgroundGroup}
              alt=""
              loading="eager"
              className="h-full hidden sm:block w-40 2xl:w-50 rounded-lg"
            />
          </div>
        </div>
      </div>
    </MainLayout>
  )
}
