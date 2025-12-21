import { createFileRoute } from '@tanstack/react-router'
import { ThreadList } from '@/features/home/components/organisms/ThreadList'
import { useThreadApp } from '@/features/home/context/ThreadAppContext'

export const Route = createFileRoute('/(thread-app)/threads/')({
  component: RouteComponent,
})

function RouteComponent() {
  const { threads, user, loading, handleDeleteThread } = useThreadApp()

  return (
    <div className="max-w-2xl mx-auto space-y-4">
      <div className="mb-6">
        <h1 className="text-2xl font-bold">All Threads</h1>
        <p className="text-muted-foreground mt-1">Browse all community threads</p>
      </div>

      <ThreadList
        threads={threads}
        currentUserId={user?.id}
        loading={loading}
        onDelete={handleDeleteThread}
      />
    </div>
  )
}
