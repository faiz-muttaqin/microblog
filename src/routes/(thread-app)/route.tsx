import { createFileRoute, Outlet } from '@tanstack/react-router'
import { Toaster } from '@/components/ui/sonner'
import { Header } from '@/features/home/components/organisms/Header'
import { Sidebar } from '@/features/home/components/organisms/Sidebar'
import { CreateThreadDialog } from '@/features/home/components/organisms/CreateThreadDialog'
import { ThreadAppProvider, useThreadApp } from '@/features/home/context/ThreadAppContext'

export const Route = createFileRoute('/(thread-app)')({
  component: RouteComponent,
})

function RouteComponent() {
  return (
    <ThreadAppProvider>
      <LayoutContent />
    </ThreadAppProvider>
  )
}

function LayoutContent() {
  const { createDialogOpen, setCreateDialogOpen, handleCreateThread } = useThreadApp()

  return (
    <div className="min-h-screen bg-background">
      <Header />

      <div className="container flex">
        <Sidebar />

        <main className="flex-1 p-4 md:p-6">
          <Outlet />
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
