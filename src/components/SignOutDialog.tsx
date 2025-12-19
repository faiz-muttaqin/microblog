import { useNavigate, useLocation } from '@tanstack/react-router'
import { ConfirmDialog } from '@/components/ConfirmDialog'
import { useAuth } from '@/hooks/use-auth'
import { toast } from 'sonner'

interface SignOutDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  onAuthDialogChange?: (state: "signin" | "signup" | null) => void
}

export function SignOutDialog({ open, onOpenChange,onAuthDialogChange }: SignOutDialogProps) {
  const navigate = useNavigate()
  const location = useLocation()
  const { signOut } = useAuth()

  const handleSignOut = async () => {
    try {
      await signOut()
      // Preserve current location for redirect after sign-in
      const currentPath = location.href
      let relatedPath = '/sign-in'
      if (!currentPath.startsWith('/dashboard')) {
        // Remove origin for relative path
        relatedPath = location.href
      }
      navigate({
        to: relatedPath,
        search: { redirect: currentPath },
        replace: true,
      })
      onOpenChange(false)
      onAuthDialogChange?.(null)
      toast.success('Signed out successfully')
    } catch (error) {
      toast.error('Failed to sign out')
      console.error('Sign out error:', error)
    }
  }

  return (
    <ConfirmDialog
      open={open}
      onOpenChange={onOpenChange}
      title='Sign out'
      desc='Are you sure you want to sign out? You will need to sign in again to access your account.'
      confirmText='Sign out'
      destructive
      handleConfirm={handleSignOut}
      className='sm:max-w-sm'
    />
  )
}
