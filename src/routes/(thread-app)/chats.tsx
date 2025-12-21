import { createFileRoute } from '@tanstack/react-router'
import { Chats } from '@/features/chats'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import {
  ResponsiveDialog,
  ResponsiveDialogContent,
} from '@/components/ui/revola'
import { UserAuthForm } from '@/components/UserAuthForm'
import { SignUpForm } from '@/components/SignUpForm'
import { useAuth } from "@/hooks/use-auth";
import useDialogState from '@/hooks/use-dialog-state'
import { MessageSquare, Lock } from 'lucide-react'

export const Route = createFileRoute('/(thread-app)/chats')({
  component: RouteComponent,
})

function RouteComponent() {
  const [authOpen, setAuthOpen] = useDialogState<'signin' | 'signup'>(null)
  const { user: authenticatedUser } = useAuth()
  
  if (!authenticatedUser) {
    return (
      <div className="flex items-center justify-center min-h-[calc(100vh-200px)]">
        <Card className="w-full max-w-md mx-4">
          <CardHeader className="text-center space-y-4">
            <div className="mx-auto w-16 h-16 bg-primary/10 rounded-full flex items-center justify-center">
              <MessageSquare className="h-8 w-8 text-primary" />
            </div>
            <div className="space-y-2">
              <CardTitle className="text-2xl">Welcome to Chats</CardTitle>
              <CardDescription className="text-base">
                Sign in to start chatting with other members of the community
              </CardDescription>
            </div>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex flex-col gap-2">
              <Button 
                onClick={() => setAuthOpen('signin')}
                size="lg"
                className="w-full"
              >
                <Lock className="h-4 w-4 mr-2" />
                Sign In
              </Button>
              <Button 
                variant="outline" 
                onClick={() => setAuthOpen('signup')}
                size="lg"
                className="w-full"
              >
                Create Account
              </Button>
            </div>
            
            <div className="relative">
              <div className="absolute inset-0 flex items-center">
                <span className="w-full border-t" />
              </div>
              <div className="relative flex justify-center text-xs uppercase">
                <span className="bg-background px-2 text-muted-foreground">
                  Features you'll unlock
                </span>
              </div>
            </div>

            <div className="space-y-3 text-sm text-muted-foreground">
              <div className="flex items-start gap-3">
                <div className="mt-0.5 h-5 w-5 rounded-full bg-primary/10 flex items-center justify-center shrink-0">
                  <span className="text-primary text-xs">✓</span>
                </div>
                <p>Real-time messaging with community members</p>
              </div>
              <div className="flex items-start gap-3">
                <div className="mt-0.5 h-5 w-5 rounded-full bg-primary/10 flex items-center justify-center shrink-0">
                  <span className="text-primary text-xs">✓</span>
                </div>
                <p>Create group conversations and discussions</p>
              </div>
              <div className="flex items-start gap-3">
                <div className="mt-0.5 h-5 w-5 rounded-full bg-primary/10 flex items-center justify-center shrink-0">
                  <span className="text-primary text-xs">✓</span>
                </div>
                <p>Share files, images, and rich content</p>
              </div>
            </div>
          </CardContent>
        </Card>

        <ResponsiveDialog
          open={!!authOpen}
          onOpenChange={(isOpen: boolean) => {
            if (!isOpen) setAuthOpen(null)
          }}
        >
          <ResponsiveDialogContent className="sm:max-w-lg p-8">
            {authOpen === 'signin' && <UserAuthForm />}
            {authOpen === 'signup' && <SignUpForm />}
          </ResponsiveDialogContent>
        </ResponsiveDialog>
      </div>
    )
  }
  
  return <Chats />
}
