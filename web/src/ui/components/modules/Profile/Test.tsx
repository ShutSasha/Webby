export default async function Test() {
  await new Promise(r => setTimeout(r, 3000))

  return <div>TestTestTestTest</div>
}

export function TestSkeleton() {
  return <div className="bg-neutral-700 animate-pulse h-4 w-20"></div>
}
