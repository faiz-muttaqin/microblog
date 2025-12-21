import { createContext, useContext, useState, useEffect, type ReactNode } from 'react'
import type { Thread, User } from '@/types/thread'
import { api } from '../services/api'
import { toast } from 'sonner'
import { getErrorMessage } from '../services/getErrorMessage'

interface ThreadAppContextType {
  threads: Thread[]
  user: User | undefined
  loading: boolean
  createDialogOpen: boolean
  setCreateDialogOpen: (open: boolean) => void
  loadData: () => Promise<void>
  handleCreateThread: (title: string, body: string, category: string) => Promise<void>
  handleDeleteThread: (threadId: string) => Promise<void>
}

const ThreadAppContext = createContext<ThreadAppContextType | undefined>(undefined)

export const useThreadApp = () => {
  const context = useContext(ThreadAppContext)
  if (!context) {
    throw new Error('useThreadApp must be used within ThreadAppProvider')
  }
  return context
}

interface ThreadAppProviderProps {
  children: ReactNode
}

export const ThreadAppProvider = ({ children }: ThreadAppProviderProps) => {
  const [threads, setThreads] = useState<Thread[]>([])
  const [user, setUser] = useState<User | undefined>()
  const [loading, setLoading] = useState(true)
  const [createDialogOpen, setCreateDialogOpen] = useState(false)

  useEffect(() => {
    loadData()
  }, [])

  const loadData = async () => {
    try {
      setLoading(true)

      // Try to get current user
      try {
        const userResponse = await api.getMe()
        if (userResponse.data) {
          setUser(userResponse.data)
        }
      } catch {
        // User not logged in
      }

      // Load threads
      const threadsResponse = await api.getAllThreads()
      if (threadsResponse.data) {
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

  const value: ThreadAppContextType = {
    threads,
    user,
    loading,
    createDialogOpen,
    setCreateDialogOpen,
    loadData,
    handleCreateThread,
    handleDeleteThread,
  }

  return (
    <ThreadAppContext.Provider value={value}>
      {children}
    </ThreadAppContext.Provider>
  )
}
