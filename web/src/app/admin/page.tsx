'use client'

import AdminAreaChart, { ChartDataPoint } from '@/ui/components/modules/Admin/AdminAreaChart'

const mockStats = [
  { label: 'Total Users', value: '24,592', trend: '+12%', isPositive: true },
  { label: 'Monthly Revenue', value: '$4,250', trend: '+5.4%', isPositive: true },
  { label: 'Active Rooms', value: '142', trend: '-2%', isPositive: false },
  { label: 'Pending Reports', value: '28', trend: '+12', isPositive: false },
]

const usersData: ChartDataPoint[] = [
  { label: 'Jul', value: 1200 },
  { label: 'Aug', value: 1800 },
  { label: 'Sep', value: 2100 },
  { label: 'Oct', value: 2400 },
  { label: 'Nov', value: 3200 },
  { label: 'Dec', value: 3800 },
  { label: 'Jan', value: 4100 },
  { label: 'Feb', value: 4800 },
  { label: 'Mar', value: 5400 },
  { label: 'Apr', value: 6200 },
  { label: 'May', value: 7100 },
  { label: 'Jun', value: 8500 },
]

const revenueData: ChartDataPoint[] = [
  { label: 'Jul', value: 800 },
  { label: 'Aug', value: 1100 },
  { label: 'Sep', value: 1050 },
  { label: 'Oct', value: 1400 },
  { label: 'Nov', value: 1800 },
  { label: 'Dec', value: 2400 },
  { label: 'Jan', value: 2200 },
  { label: 'Feb', value: 2800 },
  { label: 'Mar', value: 3100 },
  { label: 'Apr', value: 3500 },
  { label: 'May', value: 3900 },
  { label: 'Jun', value: 4250 },
]

export default function AdminOverviewPage() {
  return (
    <div className="flex flex-col gap-6 animate-in fade-in duration-500 pb-10">
      {/* Quick Stats Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        {mockStats.map((stat, i) => (
          <div key={i} className="bg-[#0A0A0A] border border-neutral-800/60 rounded-2xl p-5 flex flex-col gap-2">
            <span className="text-neutral-500 text-sm font-medium">{stat.label}</span>
            <div className="flex items-end justify-between">
              <span className="text-3xl font-bold text-neutral-100">{stat.value}</span>
              <span className={`text-sm font-semibold mb-1 ${stat.isPositive ? 'text-emerald-500' : 'text-red-500'}`}>
                {stat.trend}
              </span>
            </div>
          </div>
        ))}
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mt-2">
        <AdminAreaChart
          title="User Growth (Last 12 Months)"
          data={usersData}
          colorTheme="emerald"
          valueFormatter={val => val.toLocaleString()}
        />

        <AdminAreaChart
          title="Revenue (Last 12 Months)"
          data={revenueData}
          colorTheme="purple"
          valueFormatter={val => `$${val.toLocaleString()}`}
        />
      </div>
    </div>
  )
}
