import { Link } from '@tanstack/react-router'
import {  Moon, Sun } from 'lucide-react'
import { GiKite } from "react-icons/gi";
import { Button } from '@/components/ui/button'
import { useTheme } from '@/context/theme-provider'
import { ProfileDropdown } from '@/components/ProfileDropdown'
import { ConfigDrawer } from '@/components/ConfigDrawer'

export const Header = () => {
    const { theme, setTheme } = useTheme()

    return (
        <header className="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-backdrop:bg-background/60">
            <div className="container flex h-14 items-center justify-between">
                <Link to="/" className="flex items-center gap-2">
                    <GiKite className="h-6 w-6 text-primary" />
                    <span className="font-bold text-xl">KITE</span>
                </Link>

                <div className="flex items-center gap-2">
                    <Button
                        variant="ghost"
                        size="icon"
                        onClick={() => setTheme(theme === 'dark' ? 'light' : 'dark')}
                    >
                        {theme === 'dark' ? <Sun className="h-5 w-5" /> : <Moon className="h-5 w-5" />}
                    </Button>
                    <ConfigDrawer />
                    <ProfileDropdown />

                </div>
            </div>
        </header>
    )
}
