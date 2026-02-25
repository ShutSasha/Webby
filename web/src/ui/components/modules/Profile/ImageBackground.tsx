import Image from 'next/image'

export default function ImageBackground({ src }: { src: string }) {
  return (
    <Image
      src={src}
      alt=""
      width={700}
      height={394}
      className="w-full rounded-2xl mb-1 group-hover:scale-95 transition-all duration-300 aspect-video object-cover"
      loading="lazy"
      placeholder="blur"
      blurDataURL="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mMUrwcAALMAmGjO2MQAAAAASUVORK5CYII="
    />
  )
}
