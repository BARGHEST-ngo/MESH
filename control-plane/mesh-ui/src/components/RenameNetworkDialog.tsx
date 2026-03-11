import { useState, useEffect } from 'react'
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription, DialogFooter } from './ui/dialog'
import { Button } from './ui/button'
import { Input } from './ui/input'
import { Label } from './ui/label'
import { useRenameNetwork } from '../api/useNetworks'
import type { V1User } from '../api/openapi/types.gen'

interface RenameNetworkDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  network: V1User | null
}

export function RenameNetworkDialog({ open, onOpenChange, network }: RenameNetworkDialogProps) {
  const [newName, setNewName] = useState('')
  
  const renameNetwork = useRenameNetwork()

  useEffect(() => {
    if (network) {
      setNewName(network.name || '')
    }
  }, [network])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    
    if (!network?.name || !newName.trim() || newName === network.name) return

    try {
      await renameNetwork.mutateAsync({
        user: network,
        newName: newName.trim(),
      })
      
      onOpenChange(false)
    } catch (error) {
      console.error('Failed to rename network:', error)
      alert('Failed to rename network. The new name might already be in use.')
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <form onSubmit={handleSubmit}>
          <DialogHeader>
            <DialogTitle>&gt; RENAME_NETWORK</DialogTitle>
            <DialogDescription>
              Change the name of network: <span className="text-primary font-semibold">{network?.name}</span>
            </DialogDescription>
          </DialogHeader>

          <div className="space-y-4 my-6">
            <div>
              <Label htmlFor="newName">New Network Name</Label>
              <Input
                id="newName"
                value={newName}
                onChange={(e) => setNewName(e.target.value)}
                placeholder="new-network-name"
                required
                className="mt-1"
              />
              <p className="text-xs text-muted-foreground mt-1">
                Lowercase, no spaces (use hyphens or underscores)
              </p>
            </div>
          </div>

          <DialogFooter>
            <Button
              type="button"
              variant="outline"
              onClick={() => onOpenChange(false)}
              disabled={renameNetwork.isPending}
            >
              [ CANCEL ]
            </Button>
            <Button 
              type="submit" 
              disabled={renameNetwork.isPending || !newName.trim() || newName === network?.name}
            >
              {renameNetwork.isPending ? '[ RENAMING... ]' : '[ RENAME ]'}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}

