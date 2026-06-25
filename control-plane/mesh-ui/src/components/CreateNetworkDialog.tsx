import { useState } from 'react'
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription, DialogFooter } from './ui/dialog'
import { Button } from './ui/button'
import { Input } from './ui/input'
import { Label } from './ui/label'
import { useCreateNetwork } from '../api/useNetworks'
import { Shield } from 'lucide-react'
import { isValidNetworkName, networkNameText } from '../lib/aclPolicyGenerator'


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
  const trimmedName = name.trim()
  const isInvalid = trimmedName.length > 0 && !isValidNetworkName(trimmedName)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!trimmedName) return
    if (!isValidNetworkName(trimmedName)) {
      setError(networkNameText)
      return
    }

    setError('')

    try {
      await createNetwork.mutateAsync({
        name: trimmedName,
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
            <DialogTitle>Create network</DialogTitle>
            <DialogDescription>
              Create a new network namespace for organizing devices.
            </DialogDescription>
          </DialogHeader>

          <div className="space-y-4 my-5">
            <div>
              <Label htmlFor="name" className="mb-2">Network name</Label>
              <Input
                id="name"
                value={name}
                onChange={(e) => {
                  setName(e.target.value)
                  setError('')
                }}
                placeholder="network-1"
                required
                className="font-mono"
                aria-invalid={!!error || isInvalid}
              />
              {error ? (
                <p className="text-xs text-destructive mt-1.5 font-semibold">
                  {error}
                </p>
              ) : (
                <p className={`text-xs mt-1.5 ${isInvalid ? 'text-destructive font-semibold' : 'text-soft'}`}>
                  {networkNameText}
                </p>
              )}
            </div>

            <div>
              <Label htmlFor="displayName" className="mb-2">Display name</Label>
              <Input
                id="displayName"
                value={displayName}
                onChange={(e) => setDisplayName(e.target.value)}
                placeholder="Network 1"
                required
              />
            </div>

            <div>
              <Label htmlFor="email" className="mb-2">Analyst email</Label>
              <Input
                id="email"
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder="admin@example.com"
                required
              />
            </div>

            <div className="flex items-start gap-3 p-3.5 bg-primary/10 border border-primary/20 rounded-xl">
              <Shield aria-hidden className="mt-0.5 h-4 w-4 text-primary shrink-0" />
              <div className="text-sm">
                <p className="font-semibold text-foreground">Network isolation enabled</p>
                <p className="text-text2 mt-1 leading-relaxed">
                  Devices in this network can only reach other devices in the same network.
                  Cross-network traffic is blocked by ACL policy.
                </p>
              </div>
            </div>
          </div>

          <DialogFooter>
            <Button
              type="button"
              variant="ghost"
              onClick={() => onOpenChange(false)}
              disabled={createNetwork.isPending}
            >
              Cancel
            </Button>
            <Button type="submit" disabled={createNetwork.isPending || !trimmedName || isInvalid}>
              {createNetwork.isPending ? 'Creating…' : 'Create network'}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
