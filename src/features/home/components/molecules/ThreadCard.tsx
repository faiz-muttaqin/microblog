// import { useState, useEffect } from 'react'
import { MessageSquare, Share, MoreHorizontal } from 'lucide-react'
import { Card, CardContent, CardHeader } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { CardInfoHeader } from './CardInfoHeader'
import { VoteActions } from './VoteActions'
import type { Thread } from '@/types/thread'

interface ThreadCardProps {
    thread: Thread
    currentUserId?: string
    onVote?: (threadId: string, voteType: 'up' | 'down' | 'neutral') => void
    onDelete?: (threadId: string) => void
}

export const ThreadCard = ({ thread, currentUserId, onVote, onDelete }: ThreadCardProps) => {
    // Derive vote state from props to avoid synchronous state updates in effects
    const isUpVoted = currentUserId
        ? !!thread.up_voted_by_me || (thread.upVotesBy || []).includes(currentUserId)
        : false
    const isDownVoted = currentUserId
        ? !!thread.down_voted_by_me || (thread.downVotesBy || []).includes(currentUserId)
        : false

    const handleUpVote = () => {
        if (!onVote || !currentUserId) return

        if (isUpVoted) {
            onVote(thread.id, 'neutral')
        } else {
            onVote(thread.id, 'up')
        }
    }

    const handleDownVote = () => {
        if (!onVote || !currentUserId) return

        if (isDownVoted) {
            onVote(thread.id, 'neutral')
        } else {
            onVote(thread.id, 'down')
        }
    }

    const isOwner = currentUserId === thread.user_id

    return (
        <Card className="hover:bg-accent/50 transition-colors gap-0">
            <CardHeader className="">
                <div className="flex items-start justify-between gap-3">
                    <div className="flex-1 min-w-0">
                        <CardInfoHeader user={thread.user} createdAt={thread.created_at} category={thread.category} />
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

            <CardContent className="">
                <a href={`/threads/${thread.id}`}>
                    <div className="space-y-2 cursor-pointer">
                        <h3 className="font-bold text-lg hover:underline">{thread.title}</h3>
                        <p className="text-sm text-muted-foreground line-clamp-3">{thread.body}</p>
                    </div>
                </a>

                <div className="flex items-center gap-6 pt-2">
                    <VoteActions
                        upVotesCount={thread.total_up_votes ?? (thread.upVotesBy || []).length}
                        downVotesCount={thread.total_down_votes ?? (thread.downVotesBy || []).length}
                        isUpVoted={isUpVoted}
                        isDownVoted={isDownVoted}
                        onUpVote={handleUpVote}
                        onDownVote={handleDownVote}
                        disabled={!currentUserId}
                    />

                    <a href={`/threads/${thread.id}`}>
                        <Button variant="ghost" size="sm" className="h-8 px-2 gap-1">
                            <MessageSquare className="h-4 w-4" />
                            <span className="text-sm">{thread.total_comments || 0}</span>
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
