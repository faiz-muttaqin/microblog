import { formatDistanceToNow } from 'date-fns'
import type { User } from '../../types'
import { ProfileDropdown } from '@/components/ProfileDropdown'

interface UserInfoProps {
    user?: User
    createdAt: string
}

export const UserInfo = ({ user, createdAt }: UserInfoProps) => {
    if (!user) {
        return null
    }

    return (
        <div className="flex items-start gap-3">
            <ProfileDropdown/>
            <div className="flex-1 min-w-0">
                <div className="flex items-center gap-2">
                    <span className="font-semibold text-sm truncate">{user.name}</span>
                    <span className="text-xs text-muted-foreground">
                        {formatDistanceToNow(new Date(createdAt), { addSuffix: true })}
                    </span>
                </div>
            </div>
        </div>
    )
}
