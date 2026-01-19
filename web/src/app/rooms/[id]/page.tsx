import MainLayout from '@/ui/components/MainLayout'
import CustomPlayer from '@/ui/components/modules/Player/CustomPlayer'

export default function RoomPage() {
  return (
    <MainLayout>
      <div className="flex flex-col w-full bg-neutral-900 rounded-[20px] p-5 gap-4 box-border">
        <div className="flex gap-5">
          {/* TODO: remove hardcoded videoUrl */}

          {/* Short video (3 minutes) https://www.youtube.com/watch?v=r1753DCSO4M&list=RDHURv04Ogd0g&index=37 */}
          {/* Longer video (1 hour) https://www.youtube.com/watch?v=Zmrj90wYt4c */}
          <CustomPlayer videoUrl="http://commondatastorage.googleapis.com/gtv-videos-bucket/sample/WhatCarCanYouGetForAGrand.mp4" />

          <div className="hidden xl:block xl:w-[300px] 2xl:w-[340px] flex-none bg-amber-700 rounded-2xl" />
        </div>
        <p>some footer content</p>
        <p>some footer content</p>
        <p>some footer content</p>
        <p>some footer content</p>
      </div>
    </MainLayout>
  )
}
