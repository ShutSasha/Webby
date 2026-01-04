'use client'

import { navMobileItems } from '@/lib/placeholder-data/nav-side'
import { useCommonStore } from '@/stores/common.store'

import MobileDrawer from './MobileDrawer'
import NavElement from './NavElement'

export default function MobileNav() {
  const isMobileNavOpen = useCommonStore(state => state.isMobileNavOpen)
  const toggleMobileNav = useCommonStore(state => state.toggleMobileNav)

  return (
    <nav
      className={`flex w-full md:hidden bg-neutral-900 transition-all duration-300 ease-out text-sm leading-5 sticky
        bottom-0 border-t border-border`}
    >
      <ul className="flex justify-between items-center list-none w-full gap-2 max-[360px]:px-1.5 px-4">
        {navMobileItems.map(item => (
          <li key={item.key} className="min-w-0">
            <NavElement
              href={item.href}
              text={item.text}
              Icon={item.icon}
              iconSize={item.iconSize}
              isLogo={item.key === 'logo'}
              toggleMobileNav={item.key === 'menu' ? toggleMobileNav : undefined}
            />
          </li>
        ))}
      </ul>
      <MobileDrawer isOpen={isMobileNavOpen} onClose={toggleMobileNav}>
        <div>hello</div>
      </MobileDrawer>
    </nav>
  )
}
