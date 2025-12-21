import { createFileRoute } from '@tanstack/react-router'
import { ThreadList } from '@/features/home/components/organisms/ThreadList'
import { useThreadApp } from '@/features/home/context/ThreadAppContext'
import { useMemo } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'

export const Route = createFileRoute('/(thread-app)/profile')({
  component: RouteComponent,
})

function RouteComponent() {
  const { threads, user, loading, handleDeleteThread } = useThreadApp()

  const userThreads = useMemo(() => {
    if (!user) return []
    return threads.filter(t => t.user_id === user.id)
  }, [threads, user])

  if (!user) {
    return (
      <div className="max-w-2xl mx-auto">
        <Card>
          <CardHeader>
            <CardTitle>Not Logged In</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-muted-foreground mb-4">Please log in to view your profile</p>
            <Button asChild>
              <a href="/login">Log In</a>
            </Button>
          </CardContent>
        </Card>
      </div>
    )
  }

  return (
    <div className="max-w-2xl mx-auto space-y-6">
      <Card>
        <CardHeader>
          <CardTitle className="text-2xl">{user.username}</CardTitle>
          <p className="text-muted-foreground">{user.email}</p>
        </CardHeader>
        <CardContent>
          <div className="flex gap-6 text-sm">
            <div>
              <span className="font-semibold">{userThreads.length}</span>
              <span className="text-muted-foreground ml-1">Threads</span>
            </div>
          </div>
        </CardContent>
      </Card>

      <div>
        <h2 className="text-xl font-bold mb-4">Your Threads</h2>
        <ThreadList
          threads={userThreads}
          currentUserId={user.id}
          loading={loading}
          onDelete={handleDeleteThread}
        />
      </div>
    </div>
  )
}
