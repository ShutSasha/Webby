import Image from 'next/image'

import PlusIcon from '@/assets/icons/ic_plus_create.svg'
import MainLayout from '@/ui/components/MainLayout'
import CustomPlayer from '@/ui/components/modules/Player/CustomPlayer'
import ComplaintButton from '@/ui/components/modules/Profile/ComplaintButton'
import FollowButton from '@/ui/components/modules/Profile/FollowButton'
import AsideVideoCard from '@/ui/components/modules/Videos/AsideVideoCard'
import VideoDescription from '@/ui/components/modules/Videos/VideoDescription'
import { BLUR_DATA_URLS } from '@/ui/images'

export default function VideoPage() {
  return (
    <MainLayout>
      <div className="flex flex-col w-full bg-neutral-900 rounded-[20px] p-5 gap-4 box-border">
        <div className="flex gap-5">
          <div className="flex-1 min-w-0">
            <CustomPlayer videoUrl="https://www.youtube.com/watch?v=Caqv5St2axA&list=RDCaqv5St2axA&start_radio=1" />

            <div className="flex items-center justify-between mt-3 mb-2">
              <p className="text-neutral-300 text-[20px] font-bold">VIDEO TITLE VIDEO TITLE VIDEO TITLE</p>
              <div className="flex gap-3 items-center">
                <button
                  className="flex items-center gap-2 px-4 py-2 bg-neutral-800 rounded-full transition-all duration-300
                    hover:bg-neutral-700/40 cursor-pointer active:scale-90 group border border-transparent"
                >
                  <PlusIcon className="size-4 text-emerald-500" />
                  <p className="text-sm leading-3.5">Add to playlist</p>
                </button>
                {/* TODO: add real users' id */}
                <ComplaintButton authorId="id" targetId="id" />
              </div>
            </div>

            <div className="flex items-center gap-3">
              <Image
                src="https://i.ibb.co/bgB69GMR/thumb-1920-569355.png"
                className="size-9 object-cover rounded-full"
                alt=""
                width={50}
                height={50}
                loading="lazy"
                placeholder="blur"
                blurDataURL={BLUR_DATA_URLS['neutral900']}
              />
              <p className="text-[16px] font-medium">guuuntersteam</p>
              {/* TODO: change to real props */}
              <FollowButton targetUserId="id" initialIsFollowing={false} />
            </div>

            <VideoDescription text={`${[...new Array(50)].map(_ => 'description ')}`} />
          </div>

          <div className="w-[418px] flex flex-col gap-3">
            {[...new Array(20)].map((_, index) => (
              <AsideVideoCard key={index} />
            ))}
          </div>
        </div>
      </div>
    </MainLayout>
  )
}
