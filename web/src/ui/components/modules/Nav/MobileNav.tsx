'use client'

import { navMobileElements } from '@/lib/placeholder-data/nav-side'
import { useCommonStore } from '@/stores/common.store'

import MobileDrawer from './MobileDrawer'
import { MobileNavElement } from './NavElements'

export default function MobileNav() {
  const isMobileNavOpen = useCommonStore(state => state.isMobileNavOpen)
  const toggleMobileNav = useCommonStore(state => state.toggleMobileNav)

  return (
    <nav
      className={`flex w-full lg:hidden bg-neutral-900 transition-all duration-300 ease-out text-sm leading-5 sticky
        bottom-0 border-t border-border z-999`}
    >
      <ul className="flex justify-between items-center list-none w-full gap-2 px-4 py-2">
        {navMobileElements.map(item => (
          <li key={item.key} className="min-w-0">
            <MobileNavElement
              href={item.href}
              text={item.text}
              Icon={item.icon}
              iconSize={item.iconSize}
              toggleMobileNav={toggleMobileNav}
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
