export default function Splitter(props: { text: string }) {
  return (
    <div className="flex items-center gap-2.5 my-4">
      <div className="h-px w-full bg-border" />
      <p className="shrink-0 text-emerald-500 text-[18px]">{props.text}</p>
      <div className="h-px w-full bg-border" />
    </div>
  )
}
