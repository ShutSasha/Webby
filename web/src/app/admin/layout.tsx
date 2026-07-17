import { isValidRole } from '@/lib/actions/auth.actions'
import MainContainer from '@/ui/components/layouts/MainContainer'
import MainLayout from '@/ui/components/layouts/MainLayout'
import AdminNavTabs from '@/ui/components/modules/Admin/AdminNavTabs'
import ForceLogout from '@/ui/components/shared/ForceLogout'

export const dynamic = 'force-dynamic'

export default async function AdminLayout({ children }: { children: React.ReactNode }) {
  const userHasAccessToAdminPage = await isValidRole()

  if (!userHasAccessToAdminPage.data) {
    return <ForceLogout />
  }

  return (
    <MainLayout>
      <MainContainer>
        <div className="flex flex-col gap-6 w-full h-full">
          <div>
            <h1 className="text-3xl font-bold text-foreground-strong tracking-tight">Admin Control Panel</h1>
            <p className="text-foreground-muted mt-1">Manage platform settings, users, and content moderation.</p>
          </div>

          <AdminNavTabs />

          <div className="flex-1 w-full overflow-y-auto custom-scrollbar">{children}</div>
        </div>
      </MainContainer>
    </MainLayout>
  )
}
