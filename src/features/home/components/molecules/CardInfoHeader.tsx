import { formatDistanceToNow } from 'date-fns'
import type { User } from '@/types/thread'
import { CategoryBadge } from '../atoms/CategoryBadge'
import { Button } from '@/components/ui/button'
import { Avatar, AvatarImage, AvatarFallback } from '@/components/ui/avatar'

interface CardInfoHeaderProps {
    user?: User
    createdAt: string
    category: string
}

export const CardInfoHeader = ({ user, createdAt, category }: CardInfoHeaderProps) => {
    if (!user) {
        return null
    }
    // Get initials for avatar fallback
    const initials = user.name
        .split(' ')
        .map((n) => n[0])
        .join('')
        .toUpperCase()
        .slice(0, 2)
    return (
        <div className="flex items-center gap-4">
            <Button variant='ghost' className='relative h-8 w-8 rounded-full'>
                <Avatar className='h-8 w-8'>
                    {user.avatar && <AvatarImage src={user.avatar} alt={user.name} />}
                    <AvatarFallback>{initials}</AvatarFallback>
                </Avatar>
            </Button>
            <div className="flex-1 min-w-0">
                <div className="flex items-center gap-2">
                    <span className="font-semibold text-sm truncate">{user.name}</span>
                    <span className="text-xs text-muted-foreground">
                        {formatDistanceToNow(new Date(createdAt), { addSuffix: true })}
                    </span>
                </div>
                <CategoryBadge category={category} className="" />
            </div>
        </div>
    )
}
