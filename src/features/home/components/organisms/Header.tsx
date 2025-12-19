import { Link } from '@tanstack/react-router'
import { Hash, Moon, Sun, PenSquare } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { useTheme } from '@/context/theme-provider'
import { ProfileDropdown } from '@/components/ProfileDropdown'
import { auth } from '@/lib/firebase'
interface HeaderProps {
    onCreateThread?: () => void
}

export const Header = ({ onCreateThread }: HeaderProps) => {
    const { theme, setTheme } = useTheme()
    const firebaseUser = auth.currentUser

    return (
        <header className="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
            <div className="container flex h-14 items-center justify-between">
                <Link to="/" className="flex items-center gap-2">
                    <Hash className="h-6 w-6 text-primary" />
                    <span className="font-bold text-xl">Microblog</span>
                </Link>

                <div className="flex items-center gap-2">
                    <Button
                        variant="ghost"
                        size="icon"
                        onClick={() => setTheme(theme === 'dark' ? 'light' : 'dark')}
                    >
                        {theme === 'dark' ? <Sun className="h-5 w-5" /> : <Moon className="h-5 w-5" />}
                    </Button>

                    {firebaseUser && (
                        <Button onClick={onCreateThread} size="sm">
                            <PenSquare className="h-4 w-4 mr-2" />
                            <span className="hidden sm:inline">New Thread</span>
                        </Button>
                    )}
                    <ProfileDropdown />

                </div>
            </div>
        </header>
    )
}
