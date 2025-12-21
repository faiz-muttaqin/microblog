import { useState } from 'react'
import { MoreHorizontal } from 'lucide-react'
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
import { api } from '../../services/api'
import { toast } from 'sonner'
import { getErrorMessage } from '../../services/getErrorMessage'
interface CommentCardProps {
    threadId: string
    comment: Comment
    currentUserId?: string
    onDelete?: (commentId: string) => void
}

export const CommentCard = ({threadId, comment, currentUserId, onDelete }: CommentCardProps) => {
    // Derive initial vote state from props, then keep local state for UI updates
    const initialIsUpVoted = currentUserId
        ? !!comment.up_voted_by_me || (comment.upVotesBy || []).includes(currentUserId)
        : false
    const initialIsDownVoted = currentUserId
        ? !!comment.down_voted_by_me || (comment.downVotesBy || []).includes(currentUserId)
        : false

    const [isUpVoted, setIsUpVoted] = useState<boolean>(initialIsUpVoted)
    const [isDownVoted, setIsDownVoted] = useState<boolean>(initialIsDownVoted)

    const [totalUpVotes, setTotalUpVotes] = useState<number>(comment.total_up_votes ?? (comment.upVotesBy || []).length);
    const [totalDownVotes, setTotalDownVotes] = useState<number>(comment.total_down_votes ?? (comment.downVotesBy || []).length);
    const onVote = async (threadId: string,commentId: string, voteType: 'up' | 'down' | 'neutral') => {
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
                result = await api.upVoteComment(threadId,commentId)
            } else if (voteType === 'down') {
                result = await api.downVoteComment(threadId,commentId)
            } else {
                result = await api.neutralVoteComment(threadId,commentId)
            }

            if (result && result.data) {
                if (result.data.total_up_votes !== undefined) {
                    setTotalUpVotes(result.data.total_up_votes)
                }
                if (result.data.total_down_votes !== undefined) {
                    setTotalDownVotes(result.data.total_down_votes)
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
    // const onVote = async (commentId: string, voteType: 'up' | 'down' | 'neutral') => {

    //     // Find the comment
    //     if (!comment) return

    //     const prevIsUp = !!comment.up_voted_by_me || (comment.upVotesBy || []).includes(currentUserId)
    //     const prevIsDown = !!comment.down_voted_by_me || (comment.downVotesBy || []).includes(currentUserId)
    //     const prevUpCount = comment.total_up_votes ?? (comment.upVotesBy || []).length
    //     const prevDownCount = comment.total_down_votes ?? (comment.downVotesBy || []).length

    //     // Optimistic UI update
    //     const optimisticUpdate = (c: Comment) => {
    //         if (c.id !== commentId) return c

    //         let newUpCount = prevUpCount
    //         let newDownCount = prevDownCount
    //         let newIsUp = prevIsUp
    //         let newIsDown = prevIsDown

    //         if (voteType === 'up') {
    //             newIsUp = true
    //             newIsDown = false
    //             newUpCount = prevUpCount + (prevIsUp ? 0 : 1)
    //             if (prevIsDown) newDownCount = prevDownCount - 1
    //         } else if (voteType === 'down') {
    //             newIsUp = false
    //             newIsDown = true
    //             newDownCount = prevDownCount + (prevIsDown ? 0 : 1)
    //             if (prevIsUp) newUpCount = prevUpCount - 1
    //         } else {
    //             newIsUp = false
    //             newIsDown = false
    //             if (prevIsUp) newUpCount = prevUpCount - 1
    //             if (prevIsDown) newDownCount = prevDownCount - 1
    //         }

    //         return {
    //             ...c,
    //             total_up_votes: newUpCount,
    //             total_down_votes: newDownCount,
    //             up_voted_by_me: newIsUp,
    //             down_voted_by_me: newIsDown,
    //         }
    //     }

    //     setComments(comments.map(optimisticUpdate))

    //     try {
    //         let response
    //         if (voteType === 'up') {
    //             response = await api.upVoteComment(thread.id, commentId)
    //         } else if (voteType === 'down') {
    //             response = await api.downVoteComment(thread.id, commentId)
    //         } else {
    //             response = await api.neutralVoteComment(thread.id, commentId)
    //         }

    //         if (response.data) {
    //             // Update comment with server response
    //             const voteData = response.data
    //             setComments(comments.map(c =>
    //                 c.id === commentId
    //                     ? {
    //                         ...c,
    //                         total_up_votes: voteData.total_up_votes ?? c.total_up_votes,
    //                         total_down_votes: voteData.total_down_votes ?? c.total_down_votes,
    //                         up_voted_by_me: voteData.up_voted_by_me ?? false,
    //                         down_voted_by_me: voteData.down_voted_by_me ?? false,
    //                     }
    //                     : c
    //             ))
    //         }
    //     } catch (error: unknown) {
    //         // Revert optimistic updates on error
    //         setComments(comments.map(c =>
    //             c.id === commentId
    //                 ? {
    //                     ...c,
    //                     total_up_votes: prevUpCount,
    //                     total_down_votes: prevDownCount,
    //                     up_voted_by_me: prevIsUp,
    //                     down_voted_by_me: prevIsDown,
    //                 }
    //                 : c
    //         ))
    //         const message = getErrorMessage(error)
    //         toast.error(message || 'Failed to vote')
    //     }
    // }
    const handleUpVote = () => {
        if (!onVote || !currentUserId) return
        // If already upvoted, neutralize; otherwise upvote
        if (isUpVoted) {
            onVote(threadId, comment.id, 'neutral')
        } else {
            onVote(threadId, comment.id, 'up')
        }
    }

    const handleDownVote = () => {
        if (!onVote || !currentUserId) return
        // If already downvoted, neutralize; otherwise downvote
        if (isDownVoted) {
            onVote(threadId, comment.id, 'neutral')
        } else {
            onVote(threadId, comment.id, 'down')
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
        <div className="p-4 bg-muted/30 rounded-md hover:bg-muted/50 transition-colors">
            <div className="flex gap-3">
                <Avatar className="h-8 w-8 mt-1">
                    {comment.user?.avatar && (
                        <AvatarImage src={comment.user.avatar} alt={comment.user.username} />
                    )}
                    <AvatarFallback>{initials}</AvatarFallback>
                </Avatar>

                <div className="flex-1 min-w-0">
                    <div className="flex items-center justify-between mb-1">
                        <div className="flex items-center gap-2">
                            <span className="font-semibold text-sm">{comment.user?.username}</span>
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
                        upVotesCount={totalUpVotes}
                        downVotesCount={totalDownVotes}
                        isUpVoted={isUpVoted}
                        isDownVoted={isDownVoted}
                        onUpVote={handleUpVote}
                        onDownVote={handleDownVote}
                        disabled={!currentUserId}
                    />
                </div>
            </div>
        </div>
    )
}
