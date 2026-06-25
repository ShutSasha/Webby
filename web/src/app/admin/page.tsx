import {
  getAbsoluteStatsAction,
  getRegistrationsStatsAction,
  getSubscriptionsStatsAction,
} from '@/lib/actions/admin.actions'
import AdminAreaChart, { ChartDataPoint } from '@/ui/components/modules/Admin/AdminAreaChart'

const MONTHS_ORDER = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec']

export default async function AdminOverviewPage() {
  const [absoluteRes, registrationsRes, subscriptionsRes] = await Promise.all([
    getAbsoluteStatsAction(),
    getRegistrationsStatsAction(),
    getSubscriptionsStatsAction(),
  ])

  const absStats = absoluteRes.data
  const regStats = registrationsRes.data
  const subStats = subscriptionsRes.data

  const statsGrid = [
    {
      label: 'Total Users',
      value: absStats?.totalRegistrations?.toLocaleString() || '0',
    },
    {
      label: 'Monthly Revenue',
      value: `€${absStats?.monthRevenue?.toLocaleString() || '0'}`,
    },
    {
      label: 'Active Rooms',
      value: absStats?.totalRooms?.toLocaleString() || '0',
    },
    {
      label: 'Pending Reports',
      value: absStats?.activeReportsCount?.toLocaleString() || '0',
    },
  ]

  const usersData: ChartDataPoint[] = MONTHS_ORDER.map(month => ({
    label: month,
    value: regStats?.[month] || 0,
  }))

  const revenueData: ChartDataPoint[] = MONTHS_ORDER.map(month => ({
    label: month,
    value: subStats?.[month] || 0,
  }))

  return (
    <div className="flex flex-col gap-6 animate-in fade-in duration-500 pb-10">
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        {statsGrid.map((stat, i) => (
          <div key={i} className="bg-surface border border-border rounded-2xl p-5 flex flex-col gap-2">
            <span className="text-foreground-faint text-sm font-medium">{stat.label}</span>
            <div className="flex items-end justify-between">
              <span className="text-3xl font-bold text-foreground-secondary">{stat.value}</span>
            </div>
          </div>
        ))}
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mt-2">
        <AdminAreaChart title="User Growth (Per Month)" data={usersData} colorTheme="emerald" formatType="number" />

        <AdminAreaChart title="Revenue (Per Month)" data={revenueData} colorTheme="purple" formatType="currency" />
      </div>
    </div>
  )
}
