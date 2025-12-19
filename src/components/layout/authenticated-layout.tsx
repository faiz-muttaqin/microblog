import { Outlet } from '@tanstack/react-router'
import { SearchProvider } from '@/context/search-provider'
import { AppSidebar } from '@/components/layout/app-sidebar'
import { useSidebar } from '@/components/ui/sidebar'
import { SkipToMain } from '@/components/SkipToMain'
import { Header } from '@/components/layout/header'
import { ConfigDrawer } from '@/components/ConfigDrawer'
// import { TopNav } from '@/components/layout/top-nav'
import { Search } from '@/components/Search'
import { ThemeSwitch } from '@/components/ThemeSwitch'
import { ProfileDropdown } from '@/components/ProfileDropdown'
type AuthenticatedLayoutProps = {
  children?: React.ReactNode
}

export function AuthenticatedLayout({ children }: AuthenticatedLayoutProps) {
  const { state, isMobile } = useSidebar()

  const sidebarSizeVar = state === 'expanded' ? 'var(--sidebar-width)' : 'var(--sidebar-width-icon)'
  const mainStyle: React.CSSProperties = isMobile
    ? { width: '100%' }
    : { width: `calc(100% - ${sidebarSizeVar})` }

  return (
    <SearchProvider>
      <SkipToMain />
      <AppSidebar />

      {children ??
        <div className="absolute right-0" style={mainStyle}>
          <Header>
            {/* <TopNav links={topNav} /> */}
            <div className='ms-auto flex items-center space-x-4'>
              <Search />
              <ThemeSwitch />
              <ConfigDrawer />
              <ProfileDropdown />
            </div>
          </Header>
          <Outlet />
        </div>
      }
    </SearchProvider>
  )
}

// const topNav = [
//   {
//     title: 'Overview',
//     href: 'dashboard/overview',
//     isActive: true,
//     disabled: false,
//   },
//   {
//     title: 'Customers',
//     href: 'dashboard/customers',
//     isActive: false,
//     disabled: true,
//   },
//   {
//     title: 'Products',
//     href: 'dashboard/products',
//     isActive: false,
//     disabled: true,
//   },
//   {
//     title: 'Settings',
//     href: 'dashboard/settings',
//     isActive: false,
//     disabled: true,
//   },
// ]
