import { createFileRoute, Link } from '@tanstack/react-router'
import { useState, useEffect } from 'react'
import { ArrowLeft, MessageSquare, Share, MoreHorizontal } from 'lucide-react'
import { Card, CardContent, CardHeader } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { Skeleton } from '@/components/ui/skeleton'
import { Separator } from '@/components/ui/separator'
import { toast } from 'sonner'
import { api } from '@/features/home/services/api'
import { getErrorMessage } from '@/features/home/services/getErrorMessage'
import { CardInfoHeader } from '@/features/home/components/molecules/CardInfoHeader'
import { VoteActions } from '@/features/home/components/molecules/VoteActions'
import { CommentCard } from '@/features/home/components/molecules/CommentCard'
import { CommentForm } from '@/features/home/components/molecules/CommentForm'
import type { ThreadDetail, Comment, User } from '@/types/thread'

export const Route = createFileRoute('/(thread-app)/threads/$id')({
  component: RouteComponent,
})

function RouteComponent() {
  const { id } = Route.useParams()
  const [thread, setThread] = useState<ThreadDetail | null>(null)
  const [comments, setComments] = useState<Comment[]>([])
  const [user, setUser] = useState<User | undefined>()
  const [loading, setLoading] = useState(true)
  const [commentsLoading, setCommentsLoading] = useState(true)

  useEffect(() => {
    loadData()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [id])

  const loadData = async () => {
    try {
      setLoading(true)
      setCommentsLoading(true)

      // Get current user
      try {
        const userResponse = await api.getMe()
        if (userResponse.data) {
          setUser(userResponse.data)
        }
      } catch {
        // User not logged in
      }

      // Load thread detail
      const threadResponse = await api.getThreadDetail(id)
      if (threadResponse.data) {
        setThread(threadResponse.data)
      }

      // Load comments
      const commentsResponse = await api.getThreadAllComments(id)
      if (commentsResponse.data) {
        setComments(commentsResponse.data || [])
      }
    } catch (error: unknown) {
      const message = getErrorMessage(error)
      toast.error(message || 'Failed to load thread')
    } finally {
      setLoading(false)
      setCommentsLoading(false)
    }
  }

  const handleThreadVote = async (voteType: 'up' | 'down' | 'neutral') => {
    if (!thread) return

    try {
      let response
      if (voteType === 'up') {
        response = await api.upVoteThread(thread.id)
      } else if (voteType === 'down') {
        response = await api.downVoteThread(thread.id)
      } else {
        response = await api.neutralVoteThread(thread.id)
      }

      if (response.data) {
        // Update thread with new vote counts
        setThread({
          ...thread,
          total_up_votes: response.data.total_up_votes || thread.total_up_votes,
          total_down_votes: response.data.total_down_votes || thread.total_down_votes,
          up_voted_by_me: response.data.up_voted_by_me || false,
          down_voted_by_me: response.data.down_voted_by_me || false,
        })
      }
    } catch (error: unknown) {
      const message = getErrorMessage(error)
      toast.error(message || 'Failed to vote')
    }
  }

  const handleCommentVote = async (commentId: string, voteType: 'up' | 'down' | 'neutral') => {
    if (!thread) return

    try {
      let response
      if (voteType === 'up') {
        response = await api.upVoteComment(thread.id, commentId)
      } else if (voteType === 'down') {
        response = await api.downVoteComment(thread.id, commentId)
      } else {
        response = await api.neutralVoteComment(thread.id, commentId)
      }

      if (response.data) {
        // Update comment with new vote counts
        const voteData = response.data
        setComments(comments.map(comment =>
          comment.id === commentId
            ? {
                ...comment,
                total_up_votes: voteData.total_up_votes || comment.total_up_votes,
                total_down_votes: voteData.total_down_votes || comment.total_down_votes,
                up_voted_by_me: voteData.up_voted_by_me || false,
                down_voted_by_me: voteData.down_voted_by_me || false,
              }
            : comment
        ))
      }
    } catch (error: unknown) {
      const message = getErrorMessage(error)
      toast.error(message || 'Failed to vote')
    }
  }

  const handleCreateComment = async (content: string) => {
    if (!thread) return

    try {
      await api.createComment(thread.id, content)
      toast.success('Comment posted successfully!')
      // Reload comments
      const commentsResponse = await api.getThreadAllComments(thread.id)
      if (commentsResponse.data) {
        setComments(commentsResponse.data || [])
      }
    } catch (error: unknown) {
      const message = getErrorMessage(error)
      toast.error(message || 'Failed to post comment')
      throw error
    }
  }

  const handleDeleteThread = async () => {
    if (!thread || !confirm('Are you sure you want to delete this thread?')) return

    try {
      await api.deleteThread(thread.id)
      toast.success('Thread deleted successfully!')
      // Navigate back
      window.history.back()
    } catch (error: unknown) {
      const message = getErrorMessage(error)
      toast.error(message || 'Failed to delete thread')
    }
  }

  const handleThreadVoteClick = (voteType: 'up' | 'down') => {
    if (!thread || !user) return

    const isUpVoted = thread.up_voted_by_me || (thread.upVotesBy || []).includes(user.id)
    const isDownVoted = thread.down_voted_by_me || (thread.downVotesBy || []).includes(user.id)

    // If clicking the same vote, neutralize
    if (voteType === 'up' && isUpVoted) {
      handleThreadVote('neutral')
    } else if (voteType === 'down' && isDownVoted) {
      handleThreadVote('neutral')
    } else {
      handleThreadVote(voteType)
    }
  }

  if (loading) {
    return (
      <div className="max-w-2xl mx-auto space-y-4">
        <Skeleton className="h-10 w-32" />
        <Card>
          <CardHeader>
            <Skeleton className="h-20 w-full" />
          </CardHeader>
          <CardContent>
            <Skeleton className="h-40 w-full" />
          </CardContent>
        </Card>
      </div>
    )
  }

  if (!thread) {
    return (
      <div className="max-w-2xl mx-auto text-center py-12">
        <p className="text-muted-foreground">Thread not found</p>
        <Link to="/threads">
          <Button variant="link" className="mt-4">
            <ArrowLeft className="h-4 w-4 mr-2" />
            Back to Threads
          </Button>
        </Link>
      </div>
    )
  }

  const isUpVoted = user
    ? !!thread.up_voted_by_me || (thread.upVotesBy || []).includes(user.id)
    : false
  const isDownVoted = user
    ? !!thread.down_voted_by_me || (thread.downVotesBy || []).includes(user.id)
    : false
  const isOwner = user?.id === thread.user_id

  return (
    <div className="max-w-2xl mx-auto space-y-4">
      <Link to="/threads">
        <Button variant="ghost" size="sm">
          <ArrowLeft className="h-4 w-4 mr-2" />
          Back to Threads
        </Button>
      </Link>

      {/* Single Unified Card */}
      <Card className="hover:bg-accent/5 transition-colors gap-0 rounded">
        {/* Thread Section */}
        <CardHeader>
          <div className="flex items-start justify-between gap-4">
            <CardInfoHeader
              user={thread.user}
              createdAt={thread.created_at}
              category={thread.category}
            />
            {isOwner && (
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button variant="ghost" size="sm" className="h-8 w-8 p-0">
                    <MoreHorizontal className="h-4 w-4" />
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end">
                  <DropdownMenuItem
                    onClick={handleDeleteThread}
                    className="text-destructive focus:text-destructive"
                  >
                    Delete
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            )}
          </div>

          <div className="mt-3">
            <h1 className="text-2xl font-bold mb-2">{thread.title}</h1>
          </div>
        </CardHeader>

        <CardContent className="space-y-4">
          <p className="whitespace-pre-wrap">{thread.body}</p>

          <div className="flex items-center gap-4 pt-2">
            <VoteActions
              upVotesCount={thread.total_up_votes || 0}
              downVotesCount={thread.total_down_votes || 0}
              isUpVoted={isUpVoted}
              isDownVoted={isDownVoted}
              onUpVote={() => handleThreadVoteClick('up')}
              onDownVote={() => handleThreadVoteClick('down')}
              disabled={!user}
            />

            <div className="flex items-center gap-1 text-muted-foreground">
              <MessageSquare className="h-4 w-4" />
              <span className="text-sm">{comments.length}</span>
            </div>

            <Button variant="ghost" size="sm">
              <Share className="h-4 w-4 mr-2" />
              Share
            </Button>
          </div>

          {/* Separator */}
          <Separator className="my-6" />

          {/* Comments Section */}
          <div className="space-y-4">
            <h2 className="text-xl font-semibold">
              Comments ({comments.length})
            </h2>

            {/* Comment Form with Icon Tray */}
            {user ? (
              <div className="space-y-2">
                <CommentForm onSubmit={handleCreateComment} />
                
              </div>
            ) : (
              <div className="p-4 text-center text-muted-foreground bg-muted/30 rounded-md">
                Please log in to post a comment
              </div>
            )}

            {/* Comments List */}
            {commentsLoading ? (
              <div className="space-y-3 pt-4">
                {[...Array(3)].map((_, i) => (
                  <div key={i} className="p-4 bg-muted/30 rounded-md">
                    <Skeleton className="h-20 w-full" />
                  </div>
                ))}
              </div>
            ) : comments.length > 0 ? (
              <div className="space-y-3 pt-4">
                {comments.map((comment) => (
                  <CommentCard
                    key={comment.id}
                    comment={comment}
                    currentUserId={user?.id}
                    onVote={handleCommentVote}
                  />
                ))}
              </div>
            ) : (
              <div className="p-8 text-center text-muted-foreground bg-muted/30 rounded-md">
                No comments yet. Be the first to comment!
              </div>
            )}
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
