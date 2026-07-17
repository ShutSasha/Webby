import Link from 'next/link'

import MainLayout from '@/ui/components/layouts/MainLayout'
import Button from '@/ui/components/shared/Button'
import EmptyState from '@/ui/components/shared/EmptyState'

export default function NotFound() {
  return (
    <MainLayout>
      <div className="flex-1 flex flex-col items-center justify-center min-h-[60vh] px-4">
        <EmptyState
          title="Page Not Found"
          description="We couldn't find the page you were looking for. It might have been moved or doesn't exist."
          className="py-0"
          disableFlex
        />

        <div className="mt-8">
          <Link href="/">
            <Button viewType="confirm" paddingClasses="px-6 py-3" className="font-semibold rounded-[20px]">
              Return Home
            </Button>
          </Link>
        </div>
      </div>
    </MainLayout>
  )
}
