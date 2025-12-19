import { useState } from 'react'
import { MessageSquare, Share, MoreHorizontal } from 'lucide-react'
import { Card, CardContent, CardHeader } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { UserInfo } from './UserInfo'
import { VoteActions } from './VoteActions'
import { CategoryBadge } from '../atoms/CategoryBadge'
import type { Thread } from '../../types'
import { Link } from '@tanstack/react-router'

interface ThreadCardProps {
    thread: Thread
    currentUserId?: string
    onVote?: (threadId: string, voteType: 'up' | 'down' | 'neutral') => void
    onDelete?: (threadId: string) => void
}

export const ThreadCard = ({ thread, currentUserId, onVote, onDelete }: ThreadCardProps) => {
    const [isUpVoted, setIsUpVoted] = useState(
        currentUserId ? (thread.upVotesBy || []).includes(currentUserId) : false
    )
    const [isDownVoted, setIsDownVoted] = useState(
        currentUserId ? (thread.downVotesBy || []).includes(currentUserId) : false
    )

    const handleUpVote = () => {
        if (!onVote || !currentUserId) return

        if (isUpVoted) {
            onVote(thread.id, 'neutral')
            setIsUpVoted(false)
        } else {
            onVote(thread.id, 'up')
            setIsUpVoted(true)
            setIsDownVoted(false)
        }
    }

    const handleDownVote = () => {
        if (!onVote || !currentUserId) return

        if (isDownVoted) {
            onVote(thread.id, 'neutral')
            setIsDownVoted(false)
        } else {
            onVote(thread.id, 'down')
            setIsDownVoted(true)
            setIsUpVoted(false)
        }
    }

    const isOwner = currentUserId === thread.owner_id

    return (
        <Card className="hover:bg-accent/50 transition-colors">
            <CardHeader className="pb-3">
                <div className="flex items-start justify-between gap-3">
                    <div className="flex-1 min-w-0">
                        <UserInfo user={thread.owner} createdAt={thread.createdAt} />
                        <CategoryBadge category={thread.category} className="mt-2" />
                    </div>

                    {isOwner && onDelete && (
                        <DropdownMenu>
                            <DropdownMenuTrigger asChild>
                                <Button variant="ghost" size="icon" className="h-8 w-8">
                                    <MoreHorizontal className="h-4 w-4" />
                                </Button>
                            </DropdownMenuTrigger>
                            <DropdownMenuContent align="end">
                                <DropdownMenuItem
                                    className="text-destructive"
                                    onClick={() => onDelete(thread.id)}
                                >
                                    Delete Thread
                                </DropdownMenuItem>
                            </DropdownMenuContent>
                        </DropdownMenu>
                    )}
                </div>
            </CardHeader>

            <CardContent className="space-y-3">
                <a href={`/threads/${thread.id}`}>
                    <div className="space-y-2 cursor-pointer">
                        <h3 className="font-bold text-lg hover:underline">{thread.title}</h3>
                        <p className="text-sm text-muted-foreground line-clamp-3">{thread.body}</p>
                    </div>
                </a>

                <div className="flex items-center gap-6 pt-2">
                    <VoteActions
                        upVotesCount={(thread.upVotesBy || []).length}
                        downVotesCount={(thread.downVotesBy || []).length}
                        isUpVoted={isUpVoted}
                        isDownVoted={isDownVoted}
                        onUpVote={handleUpVote}
                        onDownVote={handleDownVote}
                        disabled={!currentUserId}
                    />

                    <a href={`/threads/${thread.id}`}>
                        <Button variant="ghost" size="sm" className="h-8 px-2 gap-1">
                            <MessageSquare className="h-4 w-4" />
                            <span className="text-sm">{thread.totalComments}</span>
                        </Button>
                    </a>

                    <Button variant="ghost" size="sm" className="h-8 px-2">
                        <Share className="h-4 w-4" />
                    </Button>
                </div>
            </CardContent>
        </Card>
    )
}
