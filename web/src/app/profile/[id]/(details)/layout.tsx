type Props = {
  children: React.ReactNode
  params: Promise<{ id: string }>
}

export default async function Layout({ children }: Props) {
  return <div className="bg-surface rounded-[20px] p-5 gap-4">{children}</div>
}
