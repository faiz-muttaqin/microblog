import { createFileRoute } from '@tanstack/react-router'
import { useState } from 'react'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { Input } from '@/components/ui/input'
import { ThreadList } from '@/features/home/components/organisms/ThreadList'
import { useThreadApp } from '@/features/home/context/ThreadAppContext'
import { useAuth } from '@/hooks/use-auth'
import { Button } from '@/components/ui/button'
import { PenSquare, ImageIcon, Link2, Smile, TrendingUp, Users, ExternalLink } from 'lucide-react'
import useDialogState from '@/hooks/use-dialog-state'
import { Avatar, AvatarImage, AvatarFallback } from '@/components/ui/avatar'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Separator } from '@/components/ui/separator'

export const Route = createFileRoute('/(thread-app)/')({
  component: RouteComponent,
})

  // Dummy data for sidebar
  const trendingTopics = [
    { tag: 'WebDevelopment', posts: 1243, trend: '+12%' },
    { tag: 'AI', posts: 892, trend: '+45%' },
    { tag: 'React', posts: 756, trend: '+8%' },
    { tag: 'TypeScript', posts: 634, trend: '+23%' },
    { tag: 'OpenSource', posts: 521, trend: '+15%' },
  ]

  const suggestedUsers = [
    { id: '1', name: 'Sarah Chen', username: '@sarahchen', avatar: '', bio: 'Full-stack developer • Tech blogger', followers: 2341 },
    { id: '2', name: 'Alex Kumar', username: '@alexk', avatar: '', bio: 'Open source enthusiast • JavaScript', followers: 1893 },
    { id: '3', name: 'Maria Garcia', username: '@mariadev', avatar: '', bio: 'UI/UX Designer • Figma expert', followers: 3102 },
  ]

function RouteComponent() {
  const { setCreateDialogOpen } = useThreadApp()
  const [, setAuthOpen] = useDialogState<'signin' | 'signup'>(null)
  const { user: authenticatedUser } = useAuth()
  const { threads, user, loading, handleDeleteThread } = useThreadApp()
  const [categoryFilter, setCategoryFilter] = useState<string>('all')

  const filteredThreads = categoryFilter === 'all'
    ? threads
    : threads.filter(t => t.category === categoryFilter)

  const categories = ['all', ...Array.from(new Set(threads.map(t => t.category)))]

  return (
    <div className="container mx-auto px-4">
      <div className="flex gap-6">
        {/* Main Content */}
        <div className="flex-1 max-w-2xl space-y-4">
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold">Latest Threads</h1>
        <Select value={categoryFilter} onValueChange={setCategoryFilter}>
          <SelectTrigger className="w-45">
            <SelectValue placeholder="Filter by category" />
          </SelectTrigger>
          <SelectContent>
            {categories.map((cat) => (
              <SelectItem key={cat} value={cat}>
                {cat === 'all' ? 'All Categories' : cat.charAt(0).toUpperCase() + cat.slice(1)}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>

      {/* Create Thread Input Box */}
      <div className="bg-card border rounded-lg p-4 shadow-sm hover:shadow-md transition-shadow">
        <div className="flex items-start gap-3">
          <div className="shrink-0 mt-1">
            {authenticatedUser ? (
              <Avatar className='h-10 w-10'>
                {authenticatedUser.avatar && <AvatarImage src={authenticatedUser.avatar} alt={authenticatedUser.username || "User"} />}
                <AvatarFallback>{authenticatedUser.username?.charAt(0).toUpperCase() || 'U'}</AvatarFallback>
              </Avatar>
            ) : (
              <div className="h-10 w-10 rounded-full bg-primary/10 flex items-center justify-center">
                <PenSquare className="h-5 w-5 text-muted-foreground" />
              </div>
            )}
          </div>

          <div className="flex-1 space-y-3">
            <Input
              placeholder={authenticatedUser ? "What's on your mind? Start a new thread..." : "Sign in to create a thread..."}
              className="border-0 bg-muted/50 focus-visible:ring-1 focus-visible:ring-ring"
              onClick={() => authenticatedUser ? setCreateDialogOpen(true) : setAuthOpen('signin')}
              readOnly
            />

            <div className="flex items-center justify-between">
              <div className="flex items-center gap-1">
                <Button
                  variant="ghost"
                  size="sm"
                  className="h-9 px-3 text-muted-foreground hover:text-foreground"
                  onClick={() => authenticatedUser ? setCreateDialogOpen(true) : setAuthOpen('signin')}
                >
                  <ImageIcon className="h-4 w-4" />
                </Button>
                <Button
                  variant="ghost"
                  size="sm"
                  className="h-9 px-3 text-muted-foreground hover:text-foreground"
                  onClick={() => authenticatedUser ? setCreateDialogOpen(true) : setAuthOpen('signin')}
                >
                  <Link2 className="h-4 w-4" />
                </Button>
                <Button
                  variant="ghost"
                  size="sm"
                  className="h-9 px-3 text-muted-foreground hover:text-foreground"
                  onClick={() => authenticatedUser ? setCreateDialogOpen(true) : setAuthOpen('signin')}
                >
                  <Smile className="h-4 w-4" />
                </Button>
              </div>

              <Button
                size="sm"
                onClick={() => authenticatedUser ? setCreateDialogOpen(true) : setAuthOpen('signin')}
                className="font-semibold"
              >
                {authenticatedUser ? (
                  <>
                    <PenSquare className="h-4 w-4 mr-2" />
                    Post Thread
                  </>
                ) : (
                  'Sign In to Post'
                )}
              </Button>
            </div>
          </div>
        </div>
      </div>

      <ThreadList
        threads={filteredThreads}
        currentUserId={user?.id}
        loading={loading}
        onDelete={handleDeleteThread}
      />
        </div>

        {/* Right Sidebar */}
        <aside className="hidden lg:block w-80 space-y-4 sticky top-4 h-fit">
          {/* What's Happening Card */}
          <Card>
            <CardHeader className="pb-3">
              <CardTitle className="text-lg flex items-center gap-2">
                <TrendingUp className="h-5 w-5 text-primary" />
                What's Happening
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-1">
              {trendingTopics.map((topic, index) => (
                <div key={index}>
                  <button className="w-full text-left p-3 rounded-md hover:bg-muted/50 transition-colors">
                    <div className="flex items-start justify-between gap-2">
                      <div className="flex-1 min-w-0">
                        <div className="flex items-center gap-2">
                          <span className="font-semibold text-sm">#{topic.tag}</span>
                          <Badge variant="secondary" className="text-xs">
                            {topic.trend}
                          </Badge>
                        </div>
                        <p className="text-xs text-muted-foreground mt-0.5">
                          {topic.posts.toLocaleString()} posts
                        </p>
                      </div>
                    </div>
                  </button>
                  {index < trendingTopics.length - 1 && <Separator className="my-1" />}
                </div>
              ))}
              <Button variant="ghost" size="sm" className="w-full mt-2">
                Show more
              </Button>
            </CardContent>
          </Card>

          {/* Who to Follow Card */}
          <Card>
            <CardHeader className="pb-3">
              <CardTitle className="text-lg flex items-center gap-2">
                <Users className="h-5 w-5 text-primary" />
                Who to Follow
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-1">
              {suggestedUsers.map((suggestedUser, index) => (
                <div key={suggestedUser.id}>
                  <div className="p-3 rounded-md hover:bg-muted/50 transition-colors">
                    <div className="flex items-start gap-3">
                      <Avatar className="h-10 w-10">
                        {suggestedUser.avatar && <AvatarImage src={suggestedUser.avatar} alt={suggestedUser.name} />}
                        <AvatarFallback>{suggestedUser.name.charAt(0)}</AvatarFallback>
                      </Avatar>
                      <div className="flex-1 min-w-0">
                        <div className="flex items-start justify-between gap-2">
                          <div className="flex-1 min-w-0">
                            <p className="font-semibold text-sm truncate">{suggestedUser.name}</p>
                            <p className="text-xs text-muted-foreground truncate">{suggestedUser.username}</p>
                          </div>
                          <Button size="sm" variant="outline" className="h-7 text-xs">
                            Follow
                          </Button>
                        </div>
                        <p className="text-xs text-muted-foreground mt-1 line-clamp-2">
                          {suggestedUser.bio}
                        </p>
                        <p className="text-xs text-muted-foreground mt-1">
                          {suggestedUser.followers.toLocaleString()} followers
                        </p>
                      </div>
                    </div>
                  </div>
                  {index < suggestedUsers.length - 1 && <Separator className="my-1" />}
                </div>
              ))}
              <Button variant="ghost" size="sm" className="w-full mt-2">
                Show more
              </Button>
            </CardContent>
          </Card>

          {/* Footer Links */}
          <Card className="bg-muted/30">
            <CardContent className="p-4">
              <div className="flex flex-wrap gap-2 text-xs text-muted-foreground">
                <button className="hover:underline">Terms</button>
                <span>•</span>
                <button className="hover:underline">Privacy</button>
                <span>•</span>
                <button className="hover:underline">Help</button>
                <span>•</span>
                <button className="hover:underline flex items-center gap-1">
                  About <ExternalLink className="h-3 w-3" />
                </button>
              </div>
              <p className="text-xs text-muted-foreground mt-2">
                © 2025 Microblog. All rights reserved.
              </p>
            </CardContent>
          </Card>
        </aside>
      </div>
    </div>
  );
}
