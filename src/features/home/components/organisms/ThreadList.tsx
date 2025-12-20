import { Loader2 } from 'lucide-react'
import { ThreadCard } from '../molecules/ThreadCard'
import type { Thread } from '@/types/thread'

interface ThreadListProps {
    threads: Thread[]
    currentUserId?: string
    loading?: boolean
    onVote?: (threadId: string, voteType: 'up' | 'down' | 'neutral') => void
    onDelete?: (threadId: string) => void
}

export const ThreadList = ({
    threads,
    currentUserId,
    loading,
    onVote,
    onDelete
}: ThreadListProps) => {
    if (loading) {
        return (
            <div className="flex justify-center py-12">
                <Loader2 className="h-8 w-8 animate-spin text-primary" />
            </div>
        )
    }
    if (threads.length === 0) {
        return (
            <div className="text-center py-12">
                <p className="text-muted-foreground">No threads found.</p>
                <p className="text-sm text-muted-foreground mt-1">Be the first to create one!</p>
            </div>
        )
    }

    return (
        <div className="space-y-4">
            {threads.map((thread) => (
                <ThreadCard
                    key={thread.id}
                    thread={thread}
                    currentUserId={currentUserId}
                    onVote={onVote}
                    onDelete={onDelete}
                />
            ))}
        </div>
    )
}
