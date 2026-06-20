import MainContainer from '@/ui/components/layouts/MainContainer'
import MainLayout from '@/ui/components/layouts/MainLayout'
import AdminNavTabs from '@/ui/components/modules/Admin/AdminNavTabs'

export default function AdminLayout({ children }: { children: React.ReactNode }) {
  return (
    <MainLayout>
      <MainContainer>
        <div className="flex flex-col gap-6 w-full h-full">
          <div>
            <h1 className="text-3xl font-bold text-white tracking-tight">Admin Control Panel</h1>
            <p className="text-neutral-400 mt-1">Manage platform settings, users, and content moderation.</p>
          </div>

          <AdminNavTabs />

          <div className="flex-1 w-full overflow-y-auto custom-scrollbar">{children}</div>
        </div>
      </MainContainer>
    </MainLayout>
  )
}
