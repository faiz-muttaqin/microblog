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
import { CardInfoHeader } from './CardInfoHeader'
import { VoteActions } from './VoteActions'
import type { Thread } from '@/types/thread'
import { api } from '../../services/api'
import { toast } from 'sonner'
import { getErrorMessage } from '../../services/getErrorMessage'
import { Link } from '@tanstack/react-router'

interface ThreadCardProps {
    thread: Thread
    currentUserId?: string
    // onVote?: (threadId: string, voteType: 'up' | 'down' | 'neutral') => void
    onDelete?: (threadId: string) => void
}

export const ThreadCard = ({ thread, currentUserId, onDelete }: ThreadCardProps) => {
    // Derive initial vote state from props, then keep local state for UI updates
    const initialIsUpVoted = currentUserId
        ? !!thread.up_voted_by_me || (thread.upVotesBy || []).includes(currentUserId)
        : false
    const initialIsDownVoted = currentUserId
        ? !!thread.down_voted_by_me || (thread.downVotesBy || []).includes(currentUserId)
        : false

    const [isUpVoted, setIsUpVoted] = useState<boolean>(initialIsUpVoted)
    const [isDownVoted, setIsDownVoted] = useState<boolean>(initialIsDownVoted)

    const [totalUpVotes, setTotalUpVotes] = useState<number>(thread.total_up_votes ?? (thread.upVotesBy || []).length);
    const [totalDownVotes, setTotalDownVotes] = useState<number>(thread.total_down_votes ?? (thread.downVotesBy || []).length);
    const [totalComments, setTotalComments] = useState<number>(thread.total_comments ?? 0);
    const onVote = async (threadId: string, voteType: 'up' | 'down' | 'neutral') => {
        if (!currentUserId) return

        const prevIsUp = isUpVoted
        const prevIsDown = isDownVoted
        const prevUpCount = totalUpVotes
        const prevDownCount = totalDownVotes

        // Optimistic UI update
        if (voteType === 'up') {
            setIsUpVoted(true)
            setIsDownVoted(false)
            setTotalUpVotes((v) => v + (prevIsUp ? 0 : 1))
            if (prevIsDown) setTotalDownVotes((v) => v - 1)
        } else if (voteType === 'down') {
            setIsUpVoted(false)
            setIsDownVoted(true)
            setTotalDownVotes((v) => v + (prevIsDown ? 0 : 1))
            if (prevIsUp) setTotalUpVotes((v) => v - 1)
        } else {
            setIsUpVoted(false)
            setIsDownVoted(false)
            if (prevIsUp) setTotalUpVotes((v) => v - 1)
            if (prevIsDown) setTotalDownVotes((v) => v - 1)
        }

        try {
            let result
            if (voteType === 'up') {
                result = await api.upVoteThread(threadId)
            } else if (voteType === 'down') {
                result = await api.downVoteThread(threadId)
            } else {
                result = await api.neutralVoteThread(threadId)
            }

            if (result && result.data) {
                if (result.data.total_up_votes !== undefined) {
                    setTotalUpVotes(result.data.total_up_votes)
                }
                if (result.data.total_down_votes !== undefined) {
                    setTotalDownVotes(result.data.total_down_votes)
                }
                if (result.data.total_comments !== undefined) {
                    setTotalComments(result.data.total_comments)
                }
            }
        } catch (error: unknown) {
            // Revert optimistic updates on error
            setIsUpVoted(prevIsUp)
            setIsDownVoted(prevIsDown)
            setTotalUpVotes(prevUpCount)
            setTotalDownVotes(prevDownCount)
            const message = getErrorMessage(error)
            toast.error(message || 'Failed to vote')
        }
    }
    const handleUpVote = () => {
        if (!currentUserId) return

        if (isUpVoted) {
            onVote(thread.id, 'neutral')
        } else {
            onVote(thread.id, 'up')
        }
    }

    const handleDownVote = () => {
        if (!currentUserId) return

        if (isDownVoted) {
            onVote(thread.id, 'neutral')
        } else {
            onVote(thread.id, 'down')
        }
    }

    const isOwner = currentUserId === thread.user_id
    const threadDetailLink = `/threads/${thread.id}`
    return (
        <Card className="hover:bg-accent/50 transition-colors gap-0 rounded">
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
                <Link to={threadDetailLink}>
                    <div className="space-y-2 cursor-pointer">
                        <h3 className="font-bold text-lg hover:underline">{thread.title}</h3>
                        <p className="text-sm text-muted-foreground line-clamp-3">{thread.body}</p>
                    </div>
                </Link>

                <div className="flex items-center gap-6 pt-2">
                    <VoteActions
                        upVotesCount={totalUpVotes}
                        downVotesCount={totalDownVotes}
                        isUpVoted={isUpVoted}
                        isDownVoted={isDownVoted}
                        onUpVote={handleUpVote}
                        onDownVote={handleDownVote}
                        disabled={!currentUserId}
                    />

                    <a href={`/threads/${thread.id}`}>
                        <Button variant="ghost" size="sm" className="h-8 px-2 gap-1">
                            <MessageSquare className="h-4 w-4" />
                            <span className="text-sm">{totalComments}</span>
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
