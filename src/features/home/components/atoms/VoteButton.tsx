import { Heart } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { cn } from '@/lib/utils'

interface VoteButtonProps {
    isActive: boolean
    count?: number
    onClick: () => void
    variant?: 'up' | 'down'
    disabled?: boolean
}

export const VoteButton = ({
    isActive,
    count,
    onClick,
    variant = 'up',
    disabled
}: VoteButtonProps) => {
    return (
        <Button
            variant="ghost"
            size="sm"
            className={cn(
                'h-8 px-2 gap-1',
                variant === 'up' && isActive && 'text-primary',
                variant === 'down' && isActive && 'text-destructive'
            )}
            onClick={onClick}
            disabled={disabled}
        >
            <Heart
                className={cn(
                    'h-4 w-4',
                    variant === 'down' && 'rotate-180',
                    isActive && 'fill-current'
                )}
            />
            {count !== undefined && <span className="text-sm">{count}</span>}
        </Button>
    )
}
