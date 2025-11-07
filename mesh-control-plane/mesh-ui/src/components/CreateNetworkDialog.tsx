import { useState } from 'react'
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription, DialogFooter } from './ui/dialog'
import { Button } from './ui/button'
import { Input } from './ui/input'
import { Label } from './ui/label'
import { useCreateNetwork } from '../api/useNetworks'

interface CreateNetworkDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function CreateNetworkDialog({ open, onOpenChange }: CreateNetworkDialogProps) {
  const [name, setName] = useState('')
  const [displayName, setDisplayName] = useState('')
  const [email, setEmail] = useState('')
  const [error, setError] = useState('')
  
  const createNetwork = useCreateNetwork()

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    
    if (!name.trim()) return

    setError('')
    
    try {
      await createNetwork.mutateAsync({
        name: name.trim(),
        displayName: displayName.trim() || undefined,
        email: email.trim() || undefined,
      })
      
      // Reset form and close dialog
      setName('')
      setDisplayName('')
      setEmail('')
      setError('')
      onOpenChange(false)
    } catch (error: any) {
      console.error('Failed to create network:', error)
      
      // Check if it's a unique constraint error
      const errorMessage = error?.response?.data?.message || error?.message || String(error)
      if (errorMessage.includes('UNIQUE constraint') || errorMessage.includes('already exists')) {
        setError('Network already exists')
      } else {
        setError('Failed to create network')
      }
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <form onSubmit={handleSubmit}>
          <DialogHeader>
            <DialogTitle>&gt; CREATE_NETWORK</DialogTitle>
            <DialogDescription>
              Create a new network namespace for organizing devices.
            </DialogDescription>
          </DialogHeader>

          <div className="space-y-4 my-6">
            <div>
              <Label htmlFor="name">Network Name (required)</Label>
              <Input
                id="name"
                value={name}
                onChange={(e) => {
                  setName(e.target.value)
                  setError('')
                }}
                placeholder="network-1"
                required
                className={`mt-1 ${error ? 'border-red-500 focus-visible:ring-red-500' : ''}`}
              />
              {error ? (
                <p className="text-xs text-red-500 mt-1 font-semibold">
                  {error}
                </p>
              ) : (
                <p className="text-xs text-muted-foreground mt-1">
                  Lowercase, no spaces (use hyphens or underscores)
                </p>
              )}
            </div>

            <div>
              <Label htmlFor="displayName">Display Name (optional)</Label>
              <Input
                id="displayName"
                value={displayName}
                onChange={(e) => setDisplayName(e.target.value)}
                placeholder="Network 1"
                className="mt-1"
              />
            </div>

            <div>
              <Label htmlFor="email">Analyst Email (optional)</Label>
              <Input
                id="email"
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder="admin@example.com"
                className="mt-1"
              />
            </div>
          </div>

          <DialogFooter>
            <Button
              type="button"
              variant="outline"
              onClick={() => onOpenChange(false)}
              disabled={createNetwork.isPending}
            >
              [ CANCEL ]
            </Button>
            <Button type="submit" disabled={createNetwork.isPending || !name.trim()}>
              {createNetwork.isPending ? '[ CREATING... ]' : '[ CREATE ]'}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}

