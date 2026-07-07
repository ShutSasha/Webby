'use client'

import { useId, useCallback } from 'react'

import { Area, AreaChart, CartesianGrid, ResponsiveContainer, Tooltip, XAxis, YAxis } from 'recharts'

export type ChartDataPoint = {
  label: string
  value: number
}

type Props = {
  title: string
  data: ChartDataPoint[]
  colorTheme?: 'emerald' | 'purple' | 'blue'
  formatType?: 'number' | 'currency'
}

const THEMES = {
  emerald: { stroke: '#10b981', stop: '#10b981' },
  purple: { stroke: '#a855f7', stop: '#a855f7' },
  blue: { stroke: '#3b82f6', stop: '#3b82f6' },
}

const CustomTooltip = ({ active, payload, label, formatter, colorTheme }: any) => {
  if (active && payload && payload.length) {
    const colorClass =
      colorTheme === 'purple' ? 'text-purple-400' : colorTheme === 'blue' ? 'text-blue-400' : 'text-emerald-400'

    return (
      <div className="bg-surface border border-border/80 p-4 rounded-2xl shadow-2xl backdrop-blur-md">
        <p className="text-foreground-faint text-xs font-medium mb-1 uppercase tracking-wider">{label}</p>
        <p className={`font-bold text-2xl ${colorClass}`}>
          {formatter ? formatter(payload[0].value) : payload[0].value}
        </p>
      </div>
    )
  }
  return null
}

export default function AdminAreaChart({ title, data, colorTheme = 'emerald', formatType = 'number' }: Props) {
  const chartId = useId()
  const theme = THEMES[colorTheme]

  const formatValue = useCallback(
    (value: number) => {
      if (formatType === 'currency') return `€${value.toLocaleString()}`
      return value.toLocaleString()
    },
    [formatType],
  )

  return (
    <div className="bg-surface border border-border rounded-2xl p-6 flex flex-col w-full h-full">
      <h3 className="text-lg font-bold text-foreground-secondary mb-6">{title}</h3>

      <div className="w-full h-72">
        <ResponsiveContainer width="100%" height="100%" minWidth={0} minHeight={0}>
          <AreaChart data={data} margin={{ top: 10, right: 0, left: -20, bottom: 0 }}>
            <defs>
              <linearGradient id={`colorGradient-${chartId}`} x1="0" y1="0" x2="0" y2="1">
                <stop offset="5%" stopColor={theme.stop} stopOpacity={0.3} />
                <stop offset="95%" stopColor={theme.stop} stopOpacity={0} />
              </linearGradient>
            </defs>

            <CartesianGrid strokeDasharray="3 3" stroke="#262626" vertical={false} />

            <XAxis dataKey="label" stroke="#737373" fontSize={12} tickLine={false} axisLine={false} dy={10} />
            <YAxis stroke="#737373" fontSize={12} tickLine={false} axisLine={false} tickFormatter={formatValue} />

            <Tooltip
              content={<CustomTooltip formatter={formatValue} colorTheme={colorTheme} />}
              cursor={{ stroke: '#404040', strokeWidth: 1, strokeDasharray: '4 4' }}
            />

            <Area
              type="monotone"
              dataKey="value"
              stroke={theme.stroke}
              strokeWidth={3}
              fillOpacity={1}
              fill={`url(#colorGradient-${chartId})`}
              activeDot={{ r: 6, fill: theme.stroke, stroke: '#0A0A0A', strokeWidth: 3 }}
            />
          </AreaChart>
        </ResponsiveContainer>
      </div>
    </div>
  )
}
