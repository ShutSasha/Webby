'use client'

import { motion } from 'framer-motion'
import { Route } from 'next'
import Link from 'next/link'

import { cn } from '@/lib/utils/general.utils'

export type TabItem = {
  label: string
  href: string
  isActive: boolean
}

type Props = {
  tabs: TabItem[]
  layoutId?: string
  className?: string
  linkClassName?: string
}

export default function AnimatedTabs({ tabs, layoutId = 'active-pill', className, linkClassName }: Props) {
  return (
    <div className={cn('flex items-center gap-2', className)}>
      {tabs.map(tab => (
        <Link
          key={tab.href}
          href={tab.href as Route}
          scroll={false}
          className={cn(
            'relative px-5 py-2 z-10 block text-sm font-semibold transition-all duration-300 rounded-full',
            tab.isActive
              ? 'text-neutral-900'
              : 'text-neutral-400 hover:text-neutral-100 bg-neutral-800/40 hover:bg-neutral-800/80',
            linkClassName,
          )}
        >
          {tab.isActive && (
            <motion.div
              layoutId={layoutId}
              className="absolute inset-0 bg-neutral-100 rounded-full -z-10"
              transition={{ type: 'spring', stiffness: 400, damping: 30 }}
            />
          )}
          {tab.label}
        </Link>
      ))}
    </div>
  )
}
