import { Suspense } from 'react'

import Image from 'next/image'
import Link from 'next/link'

import Banner from '@/assets/home-page/banner2.png'
import MainLayout from '@/ui/components/layouts/MainLayout'
import AnimatedTabsContainer from '@/ui/components/modules/Rooms/AnimatedTabsContainer'
import RoomsContainer from '@/ui/components/modules/Rooms/Public/RoomsContainer'
import Button from '@/ui/components/shared/Button'

type Props = {
  searchParams: Promise<{ [key: string]: string | string[] | undefined }>
}

export default async function Home({ searchParams }: Props) {
  const resolvedSearchParams = await searchParams
  const currentTab = (resolvedSearchParams.tab as string) || 'All'
  const categoryParam = currentTab === 'All' ? '' : currentTab

  return (
    <MainLayout>
      <div className="flex flex-col gap-12 w-full max-w-5xl 2xl:max-w-7xl mx-auto pb-20">
        <div
          className="relative w-full bg-surface rounded-3xl p-8 sm:p-12 overflow-hidden flex items-center min-h-80
            transition-colors duration-300"
        >
          <div
            className="absolute right-0 top-0 bottom-0 w-1/2 hidden md:block bg-linear-to-l from-emerald-900/20
              to-transparent pointer-events-none"
          />

          <div className="relative z-10 flex flex-col gap-6 max-w-xl">
            <div>
              <h1
                className="font-bold text-4xl sm:text-5xl 2xl:text-6xl text-foreground-strong mb-3 leading-[1.15]
                  transition-colors duration-300"
              >
                Watch together.
                <br />
                Be connected.
              </h1>
              <h2 className="font-semibold text-lg sm:text-xl text-emerald-500">No downloads</h2>
            </div>
            <div className="flex flex-wrap gap-3 mt-2">
              <Link href="/rooms">
                <Button
                  className="font-semibold rounded-[20px] text-sm sm:text-[16px]"
                  paddingClasses="px-5 py-2.5 sm:px-6 sm:py-3"
                  viewType="confirm"
                >
                  Create a room
                </Button>
              </Link>
              <Link href="/rooms">
                <Button
                  className="font-semibold rounded-[20px] text-sm sm:text-[16px]"
                  paddingClasses="px-5 py-2.5 sm:px-6 sm:py-3"
                  viewType="cancel"
                >
                  Find a room
                </Button>
              </Link>
            </div>
          </div>
        </div>

        <div className="flex flex-col gap-5 bg-surface rounded-3xl p-8 sm:p-12 transition-colors duration-300">
          <h2 className="text-2xl sm:text-3xl font-bold text-foreground-strong transition-colors duration-300">
            Discover active watch rooms:
          </h2>

          <Suspense
            fallback={
              <div className="h-10 w-full animate-pulse bg-skeleton rounded-full transition-colors duration-300" />
            }
          >
            <AnimatedTabsContainer currentTab={currentTab} />
          </Suspense>

          <div className="mt-2">
            <RoomsContainer
              query=""
              category={categoryParam}
              limit={4}
              gridClassName="2xl:grid-cols-4 lg:grid-cols-2"
            />
          </div>
        </div>

        <div
          className="relative w-full bg-neutral-900 rounded-3xl p-8 sm:p-12 overflow-hidden flex items-center
            min-h-[300px] transition-colors duration-300"
        >
          <Image
            src={Banner}
            alt="banner"
            fill
            className="object-cover dark:opacity-40 transition-opacity duration-300"
          />
          <div className="absolute inset-0 bg-overlay z-0 transition-colors duration-300" />

          <div className="relative z-10 flex flex-col gap-6 max-w-2xl">
            <div>
              <h2
                className="font-bold text-3xl sm:text-4xl 2xl:text-5xl text-foreground-strong mb-3 leading-tight
                  uppercase tracking-wide transition-colors duration-300"
              >
                Find your next experience
              </h2>
              <p className="font-medium text-lg sm:text-xl text-emerald-500">
                Explore millions of videos and playlists, ready to share
              </p>
            </div>
            <div className="flex flex-wrap gap-3 mt-2">
              <Link href="/videos">
                <Button
                  className="font-semibold rounded-[20px] text-sm sm:text-[16px]"
                  paddingClasses="px-5 py-2.5 sm:px-6 sm:py-3"
                  viewType="confirm"
                >
                  Explore videos
                </Button>
              </Link>
              <Link href="/playlists">
                <Button
                  className="font-semibold rounded-[20px] text-sm sm:text-[16px]"
                  paddingClasses="px-5 py-2.5 sm:px-6 sm:py-3"
                  viewType="confirm"
                >
                  Explore playlists
                </Button>
              </Link>
            </div>
          </div>
        </div>
      </div>
    </MainLayout>
  )
}
