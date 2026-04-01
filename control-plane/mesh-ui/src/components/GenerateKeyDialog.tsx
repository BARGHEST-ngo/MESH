import { useState } from 'react'
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription, DialogFooter } from './ui/dialog'
import { Button } from './ui/button'
import { Input } from './ui/input'
import { Label } from './ui/label'
import { useCreatePreAuthKey } from '../api/usePreAuthKeys'
import { Copy, Check } from 'lucide-react'
import { encrypt } from '../lib/onboardingCrypto'
import { QRCodeCanvas } from 'qrcode.react'

interface GenerateKeyDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  networkName: string
}	

function isLocalhost(url: string): boolean {
  if (!url) return false
  const lower = url.toLowerCase()
  return (
    lower.includes('localhost') ||
    lower.includes('127.0.0.1') ||
    lower.includes('::1') ||
    lower.startsWith('http://192.168.') ||
    lower.startsWith('http://10.')
  )
}

export function GenerateKeyDialog({ open, onOpenChange, networkName }: GenerateKeyDialogProps) {
  const [reusable, setReusable] = useState(false)
  const [ephemeral, setEphemeral] = useState(false)
  const [expirationDays, setExpirationDays] = useState('1')
  const [deviceTag, setDeviceTag] = useState<'analyst' | 'mobile_node'>('analyst')
  const [generatedKey, setGeneratedKey] = useState<string | null>(null)
  const [intent, setIntent] = useState<string | null>(null)
  const [pin, setPin] = useState<string | null>(null)
  const [copied, setCopied] = useState(false)
  const [controlPlaneURL, setControlPlaneURL] = useState('') 
  const [urlError, setUrlError] = useState('')
  const [intentError, setIntentError] = useState('')
  const createKey = useCreatePreAuthKey()

  const handleGenerate = async () => {
    try {
      setUrlError('')
      if (deviceTag === 'mobile_node') {
	      const trimmedURL = controlPlaneURL.trim()
	      if (!trimmedURL) {
		      setUrlError('Control Plane URL is required for mobile nodes')
		      console.error('Control Plane URL is required for mobile nodes')
		      return
        }
        try {
          const parsed = new URL(trimmedURL)
          if (parsed.protocol !== 'https:') {
            setUrlError('Invalid URL format. Must start with https://')
            console.error('Invalid URL format')
            return
          }
        } catch {
          setUrlError('Invalid URL format. Must start with https://')
          console.error('Invalid URL format')
          return
        }
      }

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

      const authKey = result.preAuthKey?.key

      setIntentError('')

      if (authKey && deviceTag === 'mobile_node') {
        try {
          const { uri, pin } = await encrypt(controlPlaneURL.trim(), authKey)
          setPin(pin || null)
          setIntent(uri || null)
          setGeneratedKey(authKey)
          return
        } catch {
          setIntentError('Unable to generate encrypted intent')
          console.error('Unable to generate encrypted intent')
          return
        }
      }

      setGeneratedKey(authKey || null)

    } catch (error) {
      console.error('Failed to generate key:', error)
      setIntentError('Failed to generate key')
    }
  }

  const handleCopy = async () => {
    if (generatedKey) {
      await navigator.clipboard.writeText(generatedKey)
      setCopied(true)
      setTimeout(() => setCopied(false), 2000)
    }
  }

  const handleDownloadQR = () => {
    const canvas = document.querySelector('.qr-download canvas') as HTMLCanvasElement | null
    if (!canvas) return
    const a = document.createElement('a')
    a.href = canvas.toDataURL('image/png')
    a.download = 'mesh-join-qr.png'
    a.click()
  }

  const handleClose = () => {
    setGeneratedKey(null)
    setReusable(false)
    setEphemeral(false)
    setExpirationDays('1')
    setDeviceTag('analyst')
    setCopied(false)
    setUrlError('')
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
                <Input
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
                <Input
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
                    onClick={() => setDeviceTag('mobile_node')}
                    className={`flex-1 px-4 py-2 rounded border transition-colors ${
                      deviceTag === 'mobile_node'
                        ? 'bg-primary text-primary-foreground border-primary'
                        : 'bg-secondary border-border hover:bg-secondary/80'
                    }`}
                  >
                    [ MOBILE_NODE ]
                  </button>
                </div>
              </div>

              
              {deviceTag === 'mobile_node' && (
              <div>
                <Label htmlFor="controlPlaneURL">ControlPlaneURL</Label>
                <Input
                  id="controlPlaneURL"
                  type="text"
                  value={controlPlaneURL}
                  onChange={(e) => { setControlPlaneURL(e.target.value); setUrlError('') }}
                  placeholder="https://publicControlplaneURL.example.com"
                  className={`mt-2 ${urlError ? 'border-red-500 focus-visible:ring-red-500' : ''}`}
                  required
                  />
              <p className="text-xs text-muted-foreground mt-2">
                  Public URL that field devices will use to connect to you
              </p>
              {isLocalhost(controlPlaneURL) && (
                <p className="text-xs text-yellow-400 mt-1">
                Warning: This appears to be a local URL. Field devices won't be able to reach it.
                </p>
              )} 
              </div>
              )}

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

            {intentError && (
              <p className="text-xs text-red-500 mt-1">{intentError}</p>
            )}
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
              {intent && pin ? (
                <>
                  <Label>Scan QR Code</Label>
                  <div className="qr-download mt-2 flex justify-center p-4 bg-white rounded">
                    <QRCodeCanvas value={intent} size={220} />
                  </div>
                  <p className="text-xs text-muted-foreground mt-2">
                    Send this QR code image to the field operator via Signal. They scan it with Google Lens to open MESH.
                  </p>
                  <Label className="mt-4 block">PIN</Label>
                  <div className="mt-2 p-4 bg-secondary border border-border rounded font-mono text-lg tracking-widest">
                    {pin}
                  </div>
                  <p className="text-xs text-yellow-400 mt-2">
                    Read the PIN to the field operator over voice — do not send it digitally.
                  </p>
                  <Label className="mt-4 block">Auth Key (manual join)</Label>
                  <p className="text-xs text-muted-foreground mb-1">Fallback if QR is not usable.</p>
                  <div className="mt-1 p-4 bg-secondary border border-border rounded font-mono text-sm break-all">
                    {generatedKey}
                  </div>
                </>
              ) : (
                <>
                  <Label>Generated Key</Label>
                  <div className="mt-2 p-4 bg-secondary border border-border rounded font-mono text-sm break-all">
                    {generatedKey}
                  </div>
                  <p className="text-xs text-yellow-400 mt-2">
                    Copy this key now! You won't be able to see it again.
                  </p>
                </>
              )}
            </div>

            <DialogFooter>
              {intent && pin ? (
                <Button onClick={handleDownloadQR} className="flex items-center gap-2">
                  [ DOWNLOAD QR ]
                </Button>
              ) : (
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
              )}
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

