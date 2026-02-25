import Image from 'next/image'

interface VideoData {
  id: number
  src: string
  title: string
  views: string
  date: string
}

const videos: VideoData[] = [
  {
    id: 1,
    src: 'https://i.ibb.co/bgB69GMR/thumb-1920-569355.png',
    title: 'Your Name - Cinematic Scenery',
    views: '2.4M',
    date: '12/01/2026',
  },
  {
    id: 2,
    src: 'https://i.ibb.co/5h0v3Mq7/thumb-1920-415519.jpg',
    title: 'A Silent Voice: Emotional OST',
    views: '850K',
    date: '05/02/2026',
  },
  {
    id: 3,
    src: 'https://i.ibb.co/ynq0NhJG/thumb-1920-135625.jpg',
    title: 'Fantasy World Exploration',
    views: '125K',
    date: '20/12/2025',
  },
  {
    id: 4,
    src: 'https://i.ibb.co/p6JQWBDg/fantasy-scene-anime-style.jpg',
    title: 'Night Sky Aesthetic',
    views: '45K',
    date: '15/02/2026',
  },
  {
    id: 5,
    src: 'https://i.ibb.co/d4fJc7FX/anime-moon-landscape.jpg',
    title: 'Moonlight Serenade',
    views: '310K',
    date: '10/01/2026',
  },
  {
    id: 6,
    src: 'https://i.ibb.co/WWJxpyRv/anime-moon-landscape-1.jpg',
    title: 'Lofi Beats to Study to',
    views: '1.2M',
    date: '28/11/2025',
  },
  {
    id: 7,
    src: 'https://i.ibb.co/gM4DRDtV/thumb-1920-1331367.png',
    title: 'Gachiakuta Manga Review',
    views: '89K',
    date: '02/02/2026',
  },
  {
    id: 8,
    src: 'https://i.ibb.co/wF9rywMX/thumb-1920-1363137.png',
    title: 'Cyberpunk Cityscape',
    views: '215K',
    date: '14/01/2026',
  },
  {
    id: 9,
    src: 'https://i.ibb.co/twvv1WLG/thumb-1920-1311951.jpg',
    title: 'Summer Clouds Timelapse',
    views: '12K',
    date: '25/12/2025',
  },
  {
    id: 10,
    src: 'https://i.ibb.co/5XN85kT8/thumb-1920-737474.png',
    title: 'Next.js 15 Deep Dive',
    views: '340K',
    date: '01/02/2026',
  },
  {
    id: 11,
    src: 'https://i.ibb.co/pjP3rBDP/thumb-1920-736462.png',
    title: 'Retina Display Optimization',
    views: '5K',
    date: '18/02/2026',
  },
  {
    id: 12,
    src: 'https://i.ibb.co/sv7Lj07c/thumb-1920-614743.png',
    title: 'Web Development Journey',
    views: '95K',
    date: '11/12/2025',
  },
  {
    id: 13,
    src: 'https://i.ibb.co/FShJxq8/anime-style-clouds.jpg',
    title: 'Purple Skies Art Process',
    views: '28K',
    date: '22/01/2026',
  },
  {
    id: 14,
    src: 'https://i.ibb.co/60Ns8j8r/cd4af4dd04fcfba0a358cfdee5c039f7.jpg',
    title: 'Anime Landscape Speedpaint',
    views: '1.3K',
    date: '03/11/2025',
  },
  {
    id: 15,
    src: 'https://i.ibb.co/4RLdNrBC/785ca39a2a95c19e66b01b3e0615d32c.jpg',
    title: 'Secret Garden - Relaxing Mix',
    views: '2.1M',
    date: '05/01/2026',
  },
]

export default async function UserVideos() {
  // await, sync videos - Lates/Popular
  await new Promise(r => setTimeout(r, 3000))

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 2xl:grid-cols-4 gap-4">
      {videos.map(video => (
        <div key={video.id} className="group cursor-pointer">
          <Image
            className="w-full rounded-2xl mb-1 group-hover:scale-95 transition-all duration-300 aspect-video
              object-cover"
            src={video.src}
            alt={video.title}
            width={700}
            height={394}
            placeholder="blur"
            blurDataURL="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mMUrwcAALMAmGjO2MQAAAAASUVORK5CYII="
            loading="lazy"
          />
          <p className="text-sm font-medium">{video.title}</p>
          <div className="flex gap-2 text-[12px] text-neutral-500">
            <p>{video.views}</p>
            <p>{video.date}</p>
          </div>
        </div>
      ))}
    </div>
  )
}

export function UserVideosSkeleton() {
  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 2xl:grid-cols-4 gap-4">
      {Array.from({ length: 12 }).map((_, index) => (
        <div key={index}>
          <div className="w-full aspect-video rounded-2xl mb-1 bg-neutral-800 animate-pulse" />

          {/* Title */}
          <div className="h-4 w-3/4 bg-neutral-800 rounded-md mb-1.5 animate-pulse mt-1" />

          {/* Meta info (Views and Date) */}
          <div className="flex gap-2 items-center">
            <div className="h-3 w-16 bg-neutral-800 rounded-md animate-pulse" />
            <div className="h-3 w-20 bg-neutral-800 rounded-md animate-pulse" />
          </div>
        </div>
      ))}
    </div>
  )
}
