import { createFileRoute } from '@tanstack/react-router'
import { useState } from 'react'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { ThreadList } from '@/features/home/components/organisms/ThreadList'
import { useThreadApp } from '@/features/home/context/ThreadAppContext'

export const Route = createFileRoute('/(thread-app)/')({
  component: RouteComponent,
})

function RouteComponent() {
  const { threads, user, loading, handleDeleteThread } = useThreadApp()
  const [categoryFilter, setCategoryFilter] = useState<string>('all')

  const filteredThreads = categoryFilter === 'all'
    ? threads
    : threads.filter(t => t.category === categoryFilter)

  const categories = ['all', ...Array.from(new Set(threads.map(t => t.category)))]

  return (
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
  );
}
