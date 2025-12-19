import { Badge } from '@/components/ui/badge'
import { cn } from '@/lib/utils'

interface CategoryBadgeProps {
  category: string
  className?: string
}

export const CategoryBadge = ({ category, className }: CategoryBadgeProps) => {
  return (
    <Badge variant="secondary" className={cn('capitalize', className)}>
      {category}
    </Badge>
  )
}
