import Image from 'next/image'

import ImageBackground from '@/assets/auth/bg-login.png'
import { benefitsList } from '@/lib/placeholder-data/home-page'
import Button from '@/ui/components/Button'
import MainLayout from '@/ui/components/MainLayout'

export default async function Home() {
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
                viewType="Confirm"
              >
                Create a room
              </Button>
              <Button
                className="font-semibold rounded-[20px] text-sm sm:text-[16px]"
                paddingClasses="px-3 py-2 sm:px-4.5 sm:py-2.5 2xl:px-5 2xl:py-3"
                viewType="Cancel"
              >
                Find a room
              </Button>
            </div>
          </div>
          <Image
            src={ImageBackground}
            alt="Sign up background"
            loading="eager"
            className="h-full hidden sm:block w-40 2xl:w-50 rounded-lg"
          />
        </div>
        {/* Block - React to moments together */}
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
          <h2 className="font-semibold text-xl sm:text-3xl 2xl:text-4xl leading-tight mb-2.5">
            React to moments together 🔥
          </h2>
        </div>
      </div>
    </MainLayout>
  )
}
