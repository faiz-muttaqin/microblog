import { Link, useLocation } from '@tanstack/react-router'
import { Home, Hash, Trophy, User as UserIcon } from 'lucide-react'
import { cn } from '@/lib/utils'

const navItems = [
    { name: 'Home', href: '/', icon: Home },
    { name: 'Threads', href: '/threads', icon: Hash },
    { name: 'Trendings', href: '/trendings', icon: Trophy },
    { name: 'Profile', href: '/dashboard/profile', icon: UserIcon },
]

export const Sidebar = () => {
    const location = useLocation()

    return (
        <aside className="hidden md:flex w-64 flex-col border-r bg-card p-4 sticky top-14 h-[calc(100vh-3.5rem)] overflow-y-auto">
            <nav className="space-y-1">
                {navItems.map((item) => {
                    const Icon = item.icon
                    const isActive = location.pathname === item.href

                    return (
                        <Link
                            key={item.name}
                            to={item.href}
                            className={cn(
                                'flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium transition-colors',
                                isActive
                                    ? 'bg-primary text-primary-foreground'
                                    : 'text-muted-foreground hover:bg-accent hover:text-accent-foreground'
                            )}
                        >
                            <Icon className="h-5 w-5" />
                            {item.name}
                        </Link>
                    )
                })}
            </nav>

            <div className="mt-auto pt-4 border-t">
                <div className="rounded-lg bg-muted p-3">
                    <h3 className="font-semibold text-sm mb-2">Trending Topics</h3>
                    <div className="space-y-2 text-xs text-muted-foreground">
                        <div className="flex justify-between hover:text-foreground cursor-pointer transition-colors">
                            <span>#Technology</span>
                            <span>234</span>
                        </div>
                        <div className="flex justify-between hover:text-foreground cursor-pointer transition-colors">
                            <span>#Programming</span>
                            <span>189</span>
                        </div>
                        <div className="flex justify-between hover:text-foreground cursor-pointer transition-colors">
                            <span>#Design</span>
                            <span>156</span>
                        </div>
                    </div>
                </div>
            </div>
        </aside>
    )
}
