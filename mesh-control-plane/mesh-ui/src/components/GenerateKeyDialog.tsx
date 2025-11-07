import { useState } from 'react'
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription, DialogFooter } from './ui/dialog'
import { Button } from './ui/button'
import { Input } from './ui/input'
import { Label } from './ui/label'
import { useCreatePreAuthKey } from '../api/usePreAuthKeys'
import { Copy, Check } from 'lucide-react'

interface GenerateKeyDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  networkName: string
}

export function GenerateKeyDialog({ open, onOpenChange, networkName }: GenerateKeyDialogProps) {
  const [reusable, setReusable] = useState(false)
  const [ephemeral, setEphemeral] = useState(false)
  const [expirationDays, setExpirationDays] = useState('1')
  const [deviceTag, setDeviceTag] = useState<'analyst' | 'victim'>('analyst')
  const [generatedKey, setGeneratedKey] = useState<string | null>(null)
  const [copied, setCopied] = useState(false)
  
  const createKey = useCreatePreAuthKey()

  const handleGenerate = async () => {
    try {
      const days = parseInt(expirationDays) || 1
      const expiration = new Date()
      expiration.setDate(expiration.getDate() + days)

      const result = await createKey.mutateAsync({
        user: networkName,
        reusable,
        ephemeral,
        expiration: expiration.toISOString(),
        aclTags: [`tag:${deviceTag}`],
      })
      
      setGeneratedKey(result.preAuthKey?.key || null)
    } catch (error) {
      console.error('Failed to generate key:', error)
    }
  }

  const handleCopy = async () => {
    if (generatedKey) {
      await navigator.clipboard.writeText(generatedKey)
      setCopied(true)
      setTimeout(() => setCopied(false), 2000)
    }
  }

  const handleClose = () => {
    setGeneratedKey(null)
    setReusable(false)
    setEphemeral(false)
    setExpirationDays('1')
    setDeviceTag('analyst')
    setCopied(false)
    onOpenChange(false)
  }

  return (
    <Dialog open={open} onOpenChange={handleClose}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>&gt; GENERATE_AUTH_KEY</DialogTitle>
          <DialogDescription>
            Create a pre-authentication key to register new devices to <span className="text-primary font-semibold">{networkName}</span>
          </DialogDescription>
        </DialogHeader>

        {!generatedKey ? (
          <>
            <div className="space-y-4 my-6">
              <div className="flex items-center justify-between">
                <div>
                  <Label htmlFor="reusable">Reusable</Label>
                  <p className="text-xs text-muted-foreground">Allow multiple devices to use this key</p>
                </div>
                <input
                  id="reusable"
                  type="checkbox"
                  checked={reusable}
                  onChange={(e) => setReusable(e.target.checked)}
                  className="h-4 w-4"
                />
              </div>

              <div className="flex items-center justify-between">
                <div>
                  <Label htmlFor="ephemeral">Ephemeral</Label>
                  <p className="text-xs text-muted-foreground">Devices will be deleted when they disconnect</p>
                </div>
                <input
                  id="ephemeral"
                  type="checkbox"
                  checked={ephemeral}
                  onChange={(e) => setEphemeral(e.target.checked)}
                  className="h-4 w-4"
                />
              </div>

              <div>
                <Label htmlFor="deviceTag">Device Tag</Label>
                <div className="flex gap-2">
                  <button
                    type="button"
                    onClick={() => setDeviceTag('analyst')}
                    className={`flex-1 px-4 py-2 rounded border transition-colors ${
                      deviceTag === 'analyst'
                        ? 'bg-primary text-primary-foreground border-primary'
                        : 'bg-secondary border-border hover:bg-secondary/80'
                    }`}
                  >
                    [ ANALYST ]
                  </button>
                  <button
                    type="button"
                    onClick={() => setDeviceTag('victim')}
                    className={`flex-1 px-4 py-2 rounded border transition-colors ${
                      deviceTag === 'victim'
                        ? 'bg-primary text-primary-foreground border-primary'
                        : 'bg-secondary border-border hover:bg-secondary/80'
                    }`}
                  >
                    [ VICTIM ]
                  </button>
                </div>
              </div>

              <div>
                <Label htmlFor="expiration">Expiration (days)</Label>
                <Input
                  id="expiration"
                  type="number"
                  value={expirationDays}
                  onChange={(e) => setExpirationDays(e.target.value)}
                  min="1"
                  max="365"
                  className="mt-1"
                />
              </div>
            </div>

            <DialogFooter>
              <Button type="button" variant="outline" onClick={handleClose}>
                [ CANCEL ]
              </Button>
              <Button onClick={handleGenerate} disabled={createKey.isPending}>
                {createKey.isPending ? '[ GENERATING... ]' : '[ GENERATE ]'}
              </Button>
            </DialogFooter>
          </>
        ) : (
          <>
            <div className="my-6">
              <Label>Generated Key</Label>
              <div className="mt-2 p-4 bg-secondary border border-border rounded font-mono text-sm break-all">
                {generatedKey}
              </div>
              <p className="text-xs text-yellow-400 mt-2">
                Copy this key now! You won't be able to see it again.
              </p>
            </div>

            <DialogFooter>
              <Button onClick={handleCopy} className="flex items-center gap-2">
                {copied ? (
                  <>
                    <Check size={16} />
                    [ COPIED ]
                  </>
                ) : (
                  <>
                    <Copy size={16} />
                    [ COPY KEY ]
                  </>
                )}
              </Button>
              <Button onClick={handleClose} variant="outline">
                [ CLOSE ]
              </Button>
            </DialogFooter>
          </>
        )}
      </DialogContent>
    </Dialog>
  )
}

