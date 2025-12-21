import { createFileRoute } from '@tanstack/react-router'
import { ThreadList } from '@/features/home/components/organisms/ThreadList'
import { useThreadApp } from '@/features/home/context/ThreadAppContext'
import { useMemo } from 'react'

export const Route = createFileRoute('/(thread-app)/trendings')({
  component: RouteComponent,
})

function RouteComponent() {
  const { threads, user, loading, handleDeleteThread } = useThreadApp()

  // Sort by total engagement (votes + comments)
  const trendingThreads = useMemo(() => {
    return [...threads].sort((a, b) => {
      const scoreA = (a.total_up_votes || 0) + (a.total_comments || 0)
      const scoreB = (b.total_up_votes || 0) + (b.total_comments || 0)
      return scoreB - scoreA
    })
  }, [threads])

  return (
    <div className="max-w-2xl mx-auto space-y-4">
      <div className="mb-6">
        <h1 className="text-2xl font-bold">Trending Threads</h1>
        <p className="text-muted-foreground mt-1">Most active and popular discussions</p>
      </div>

      <ThreadList
        threads={trendingThreads}
        currentUserId={user?.id}
        loading={loading}
        onDelete={handleDeleteThread}
      />
    </div>
  )
}
