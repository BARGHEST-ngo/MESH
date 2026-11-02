import { useState } from 'react'
import { useNavigate } from '@tanstack/react-router'
import { Lock, AlertTriangle } from 'lucide-react'
import { Button } from './ui/button'
import { Input } from './ui/input'
import { Label } from './ui/label'
import { Card } from './ui/card'
import { Wordmark } from './Wordmark'
import { useAuth } from '../lib/auth'
import { useLogin } from '../api/useLogin'

export default function LoginForm() {
  const [authKey, setAuthKey] = useState('')
  const [error, setError] = useState('')
  const navigate = useNavigate()
  const { login: authLogin } = useAuth()
  const loginMutation = useLogin()

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')

    if (!authKey.trim()) {
      setError('Enter your authentication key to continue.')
      return
    }

    try {
      await loginMutation.mutateAsync({ authKey })

      const success = await authLogin(authKey)

      if (success) {
        navigate({ to: '/' })
      } else {
        setError('Authentication failed')
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Authentication failed')
    }
  }

  return (
    <div className="flex items-center justify-center min-h-screen bg-background p-6">
      <div className="w-full max-w-[400px]">
        <div className="flex flex-col items-center mb-7">
          <Wordmark height={40} className="text-foreground" />
          <div className="text-[15px] font-semibold text-text2 tracking-[0.3px] mt-3.5">Control Plane</div>
          <p className="text-sm text-text2 mt-1.5">Sign in to manage your forensic networks</p>
        </div>

        <Card className="p-6">
          <form onSubmit={handleSubmit}>
            <Label htmlFor="authKey" className="mb-2">Authentication key</Label>
            <Input
              id="authKey"
              type="password"
              placeholder="••••••••••••••••"
              value={authKey}
              onChange={(e) => {
                setAuthKey(e.target.value)
                setError('')
              }}
              disabled={loginMutation.isPending}
              className="font-mono"
              aria-invalid={!!error}
            />

            {error && (
              <div className="flex items-center gap-2 mt-3 px-3 py-2.5 rounded-[9px] bg-destructive-dim border border-destructive/30">
                <AlertTriangle size={15} className="text-destructive shrink-0" />
                <span className="text-sm text-foreground">{error}</span>
              </div>
            )}

            <Button type="submit" size="lg" className="w-full mt-4" disabled={loginMutation.isPending}>
              <Lock size={17} />
              {loginMutation.isPending ? 'Signing in…' : 'Sign in'}
            </Button>
          </form>

          <div className="mt-5 pt-[18px] border-t border-border">
            <div className="text-[12.5px] text-text2 mb-2">
              Need a key? Generate one from your server terminal:
            </div>
            <div className="font-mono text-[12.5px] text-text2 bg-inset border border-border rounded-[9px] px-3.5 py-2.5">
              <span className="text-soft">$</span> task apikey
            </div>
          </div>
        </Card>

        <div className="text-center mt-[18px] text-xs text-soft">
          Self-hosted · the control plane never carries forensic traffic
        </div>
      </div>
    </div>
  )
}
