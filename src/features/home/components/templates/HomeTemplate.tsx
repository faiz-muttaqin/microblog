


import { useState, useEffect } from 'react'
import { Toaster } from '@/components/ui/sonner'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { Header } from '../organisms/Header'
import { Sidebar } from '../organisms/Sidebar'
import { ThreadList } from '../organisms/ThreadList'
import { CreateThreadDialog } from '../organisms/CreateThreadDialog'
import { api } from '../../services/api'
import type { Thread, User } from '@/types/thread'
import { toast } from 'sonner'
import { getErrorMessage } from '../../services/getErrorMessage'


export const HomeTemplate = () => {
  const [threads, setThreads] = useState<Thread[]>([])
  const [user, setUser] = useState<User | undefined>()
  const [loading, setLoading] = useState(true)
  const [createDialogOpen, setCreateDialogOpen] = useState(false)
  const [categoryFilter, setCategoryFilter] = useState<string>('all')

  useEffect(() => {
    loadData()
  }, [])

  const loadData = async () => {
    try {
      setLoading(true)

      // Try to get current user
      try {
        const userResponse = await api.getMe()
        if (userResponse.data){
          setUser(userResponse.data)
        }
      } catch {
        // User not logged in
      }

      // Load threads
      const threadsResponse = await api.getAllThreads()
      if (threadsResponse.data){
        setThreads(threadsResponse.data || [])
      }
    } catch (error: unknown) {
      const message = getErrorMessage(error)
      toast.error(message || 'Failed to load data')
    } finally {
      setLoading(false)
    }
  }

  const handleCreateThread = async (title: string, body: string, category: string) => {
    try {
      await api.createThread(title, body, category)
      toast.success('Thread created successfully!')
      loadData()
    } catch (error: unknown) {
      const message = getErrorMessage(error)
      toast.error(message || 'Failed to create thread')
    }
  }

  

  const handleDeleteThread = async (threadId: string) => {
    if (!confirm('Are you sure you want to delete this thread?')) return

    try {
      await api.deleteThread(threadId)
      toast.success('Thread deleted successfully!')
      loadData()
    } catch (error: unknown) {
      const message = getErrorMessage(error)
      toast.error(message || 'Failed to delete thread')
    }
  }

  const filteredThreads = categoryFilter === 'all'
    ? threads
    : threads.filter(t => t.category === categoryFilter)

  const categories = ['all', ...Array.from(new Set(threads.map(t => t.category)))]

  return (
    <div className="min-h-screen bg-background">
      <Header
        onCreateThread={() => setCreateDialogOpen(true)}
      />

      <div className="container flex">
        <Sidebar />

        <main className="flex-1 p-4 md:p-6">
          <div className="max-w-2xl mx-auto space-y-4">
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

            <ThreadList
              threads={filteredThreads}
              currentUserId={user?.id}
              loading={loading}
              onDelete={handleDeleteThread}
            />
          </div>
        </main>
      </div>

      <CreateThreadDialog
        open={createDialogOpen}
        onClose={() => setCreateDialogOpen(false)}
        onSubmit={handleCreateThread}
      />

      <Toaster />
    </div>
  )
}
