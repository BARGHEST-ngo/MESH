import { useState } from 'react'
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription, DialogFooter } from './ui/dialog'
import { Button } from './ui/button'
import { Input } from './ui/input'
import { Label } from './ui/label'
import { Switch } from './ui/switch'
import { useCreatePreAuthKey } from '../api/usePreAuthKeys'
import { Copy, Check, Key, Lock, AlertTriangle, Laptop, Smartphone } from 'lucide-react'
import { encrypt } from '../lib/onboardingCrypto'
import { QRCodeCanvas } from 'qrcode.react'
import { useControlPlaneUrl } from '@/api/useConfig'
import { TAG_LABEL } from '../lib/tags'

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
  const [expirationDays, setExpirationDays] = useState('0')
  const [expirationHours, setExpirationHours] = useState('1')
  const [deviceTag, setDeviceTag] = useState<'analyst' | 'mobile_node'>('mobile_node')
  const [generatedKey, setGeneratedKey] = useState<string | null>(null)
  const [intent, setIntent] = useState<string | null>(null)
  const [pin, setPin] = useState<string | null>(null)
  const [copied, setCopied] = useState(false)
  const { data: configUrl = '' } = useControlPlaneUrl()
  const [controlPlaneURLOverride, setControlPlaneURL] = useState<string | null>(null)
  const controlPlaneURL = controlPlaneURLOverride ?? configUrl
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

      const hours = parseInt(expirationHours) || 0
      const days = parseInt(expirationDays) || 0
      const expiration = new Date()
      expiration.setDate(expiration.getDate() + days)
      expiration.setHours(expiration.getHours() + hours)

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
    setIntent(null)
    setPin(null)
    setReusable(false)
    setEphemeral(false)
    setExpirationDays('0')
    setExpirationHours('1')
    setDeviceTag('mobile_node')
    setCopied(false)
    setUrlError('')
    setIntentError('')
    onOpenChange(false)
  }

  return (
    <Dialog open={open} onOpenChange={handleClose}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Generate auth key</DialogTitle>
          <DialogDescription>
            Register a new device to <span className="text-primary font-semibold">{networkName}</span>
          </DialogDescription>
        </DialogHeader>

        {!generatedKey ? (
          <>
            <div className="space-y-[18px] my-5">
              <div>
                <Label className="mb-2">Device type</Label>
                <div className="flex gap-2">
                  {(['mobile_node', 'analyst'] as const).map((tg) => (
                    <button
                      key={tg}
                      type="button"
                      onClick={() => setDeviceTag(tg)}
                      className={`flex-1 h-11 rounded-[10px] flex items-center justify-center gap-2 text-sm font-semibold border transition-colors ${
                        deviceTag === tg
                          ? 'bg-primary/15 border-primary text-primary'
                          : 'bg-card border-border text-text2 hover:bg-card-hi'
                      }`}
                    >
                      {tg === 'analyst' ? <Laptop size={17} /> : <Smartphone size={17} />}
                      {TAG_LABEL[tg]}
                    </button>
                  ))}
                </div>
              </div>

              {deviceTag === 'mobile_node' && (
                <div>
                  <Label htmlFor="controlPlaneURL" className="mb-2">Control plane URL</Label>
                  <Input
                    id="controlPlaneURL"
                    type="text"
                    value={controlPlaneURL}
                    onChange={(e) => { setControlPlaneURL(e.target.value); setUrlError('') }}
                    placeholder="https://mesh.example.org"
                    className="font-mono"
                    aria-invalid={!!urlError}
                    required
                  />
                  <p className="text-xs text-soft mt-1.5">Public address field devices use to reach you.</p>
                  {urlError && <p className="text-xs text-destructive mt-1">{urlError}</p>}
                  {isLocalhost(controlPlaneURL) && (
                    <p className="text-xs text-warning mt-1">
                      Warning: this looks like a local URL. Field devices won't be able to reach it.
                    </p>
                  )}
                </div>
              )}

              <div className="border border-border rounded-xl overflow-hidden">
                <ToggleRow
                  label="Reusable"
                  sub="Allow multiple devices to use this key"
                  on={reusable}
                  onChange={setReusable}
                />
                <ToggleRow
                  label="Ephemeral"
                  sub="Remove the device when it disconnects"
                  on={ephemeral}
                  onChange={setEphemeral}
                  border
                />
              </div>

              <div className="flex gap-3">
                <div className="flex-1">
                  <Label htmlFor="expirationHr" className="mb-2">Expires in (hours)</Label>
                  <Input id="expirationHr" type="number" value={expirationHours}
                    onChange={(e) => setExpirationHours(e.target.value)} min="0" max="23" className="font-mono" />
                </div>
                <div className="flex-1">
                  <Label htmlFor="expiration" className="mb-2">Expires in (days)</Label>
                  <Input id="expiration" type="number" value={expirationDays}
                    onChange={(e) => setExpirationDays(e.target.value)} min="0" max="365" className="font-mono" />
                </div>
              </div>
            </div>

            {intentError && <p className="text-xs text-destructive mt-1">{intentError}</p>}
            <DialogFooter>
              <Button type="button" variant="ghost" onClick={handleClose}>Cancel</Button>
              <Button onClick={handleGenerate} disabled={createKey.isPending}>
                <Key size={17} />
                {createKey.isPending ? 'Generating…' : 'Generate key'}
              </Button>
            </DialogFooter>
          </>
        ) : (
          <>
            <div className="my-5">
              {intent && pin ? (
                <>
                  <Label className="mb-2">Scan to onboard the device</Label>
                  <div className="qr-download flex justify-center p-[18px] bg-white rounded-[14px] mb-2">
                    <QRCodeCanvas value={intent} size={200} />
                  </div>
                  <p className="text-[12.5px] text-text2 mb-[18px] leading-relaxed">
                    Send this QR image to the field operator over Signal. They scan it to open MESH.
                  </p>

                  <Label className="mb-2">PIN, read aloud over voice only</Label>
                  <div className="flex items-center gap-2.5 px-4 py-3.5 rounded-[11px] bg-inset border border-border mb-2">
                    <Lock size={18} className="text-warning" />
                    <span className="font-mono text-[22px] font-bold tracking-[6px] text-foreground">{pin}</span>
                  </div>
                  <div className="flex items-center gap-1.5 mb-[18px]">
                    <AlertTriangle size={13} className="text-warning" />
                    <span className="text-xs text-warning">Never send the PIN digitally.</span>
                  </div>

                  <Label className="mb-2">Auth key (manual fallback)</Label>
                  <KeyBlock value={generatedKey} />
                </>
              ) : (
                <>
                  <Label className="mb-2">Generated key</Label>
                  <KeyBlock value={generatedKey} />
                  <div className="flex items-center gap-1.5 mt-2.5">
                    <AlertTriangle size={14} className="text-warning" />
                    <span className="text-[12.5px] text-warning">Copy this now. It won't be shown again.</span>
                  </div>
                </>
              )}
            </div>

            <DialogFooter>
              {intent && pin ? (
                <Button variant="secondary" onClick={handleDownloadQR}>Download QR</Button>
              ) : (
                <Button variant="secondary" onClick={handleCopy}>
                  {copied ? <Check size={16} /> : <Copy size={16} />}
                  {copied ? 'Copied' : 'Copy key'}
                </Button>
              )}
              <Button onClick={handleClose}>Done</Button>
            </DialogFooter>
          </>
        )}
      </DialogContent>
    </Dialog>
  )
}

function ToggleRow({
  label,
  sub,
  on,
  onChange,
  border,
}: {
  label: string
  sub: string
  on: boolean
  onChange: (v: boolean) => void
  border?: boolean
}) {
  return (
    <div className={`flex items-center gap-3 px-3.5 py-3.5 bg-card ${border ? 'border-t border-border' : ''}`}>
      <div className="flex-1">
        <div className="text-sm font-semibold text-foreground">{label}</div>
        <div className="text-xs text-text2">{sub}</div>
      </div>
      <Switch checked={on} onCheckedChange={onChange} />
    </div>
  )
}

function KeyBlock({ value }: { value: string | null }) {
  return (
    <div className="font-mono text-[12.5px] text-text2 bg-inset border border-border rounded-[10px] px-3.5 py-3 break-all leading-relaxed">
      {value}
    </div>
  )
}
