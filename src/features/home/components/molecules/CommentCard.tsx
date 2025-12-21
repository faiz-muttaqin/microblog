import { MoreHorizontal } from 'lucide-react'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { Avatar, AvatarImage, AvatarFallback } from '@/components/ui/avatar'
import { formatDistanceToNow } from 'date-fns'
import { VoteActions } from './VoteActions'
import type { Comment } from '@/types/thread'

interface CommentCardProps {
    comment: Comment
    currentUserId?: string
    onVote?: (commentId: string, voteType: 'up' | 'down' | 'neutral') => void
    onDelete?: (commentId: string) => void
}

export const CommentCard = ({ comment, currentUserId, onVote, onDelete }: CommentCardProps) => {
    const isUpVoted = currentUserId
        ? !!comment.up_voted_by_me || (comment.upVotesBy || []).includes(currentUserId)
        : false
    const isDownVoted = currentUserId
        ? !!comment.down_voted_by_me || (comment.downVotesBy || []).includes(currentUserId)
        : false

    const handleVote = (voteType: 'up' | 'down') => {
        if (!onVote || !currentUserId) return

        // If clicking the same vote, neutralize
        if (voteType === 'up' && isUpVoted) {
            onVote(comment.id, 'neutral')
        } else if (voteType === 'down' && isDownVoted) {
            onVote(comment.id, 'neutral')
        } else {
            onVote(comment.id, voteType)
        }
    }

    const handleDelete = () => {
        if (onDelete) {
            onDelete(comment.id)
        }
    }

    const initials = comment.user?.name
        .split(' ')
        .map((n) => n[0])
        .join('')
        .toUpperCase()
        .slice(0, 2)

    const isOwner = currentUserId === comment.user_id

    return (
        <Card className="hover:bg-accent/5 transition-colors rounded">
            <CardContent className="p-4">
                <div className="flex gap-3">
                    <Avatar className="h-8 w-8 mt-1">
                        {comment.user?.avatar && (
                            <AvatarImage src={comment.user.avatar} alt={comment.user.name} />
                        )}
                        <AvatarFallback>{initials}</AvatarFallback>
                    </Avatar>

                    <div className="flex-1 min-w-0">
                        <div className="flex items-center justify-between mb-1">
                            <div className="flex items-center gap-2">
                                <span className="font-semibold text-sm">{comment.user?.name}</span>
                                <span className="text-xs text-muted-foreground">
                                    {formatDistanceToNow(new Date(comment.createdAt), { addSuffix: true })}
                                </span>
                            </div>

                            {isOwner && onDelete && (
                                <DropdownMenu>
                                    <DropdownMenuTrigger asChild>
                                        <Button variant="ghost" size="sm" className="h-8 w-8 p-0">
                                            <MoreHorizontal className="h-4 w-4" />
                                        </Button>
                                    </DropdownMenuTrigger>
                                    <DropdownMenuContent align="end">
                                        <DropdownMenuItem
                                            onClick={handleDelete}
                                            className="text-destructive focus:text-destructive"
                                        >
                                            Delete
                                        </DropdownMenuItem>
                                    </DropdownMenuContent>
                                </DropdownMenu>
                            )}
                        </div>

                        <p className="text-sm mb-3 whitespace-pre-wrap">{comment.content}</p>

                        <VoteActions
                            upVotesCount={comment.total_up_votes || 0}
                            downVotesCount={comment.total_down_votes || 0}
                            isUpVoted={isUpVoted}
                            isDownVoted={isDownVoted}
                            onUpVote={() => handleVote('up')}
                            onDownVote={() => handleVote('down')}
                            disabled={!currentUserId}
                        />
                    </div>
                </div>
            </CardContent>
        </Card>
    )
}
