import { useState } from 'react'
import { useNavigate } from '@tanstack/react-router'
import { Button } from './ui/button'
import { Input } from './ui/input'
import { Label } from './ui/label'
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
      setError('Please enter an authentication key')
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
    <div className="flex items-center justify-center min-h-screen bg-background">
      <div className="w-full max-w-md p-8 space-y-6 bg-card rounded border border-border shadow-2xl">
        <div className="text-center">
          <h1 className="text-3xl font-bold text-foreground mb-2 font-mono">
            &gt; MESH_CONTROL_PLANE
          </h1>
        </div>

        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="authKey" className="font-mono">&gt; AUTH_KEY:</Label>
            <Input
              id="authKey"
              type="password"
              placeholder="••••••••••••••••"
              value={authKey}
              onChange={(e) => setAuthKey(e.target.value)}
              disabled={loginMutation.isPending}
              className="w-full font-mono"
              aria-invalid={!!error}
            />
          </div>

          {error && (
            <div className="p-3 text-sm text-destructive bg-destructive/10 border border-destructive rounded font-mono">
              [ ERROR ] {error}
            </div>
          )}

          <Button
            type="submit"
            className="w-full font-mono"
            disabled={loginMutation.isPending}
          >
            {loginMutation.isPending ? '[ AUTHENTICATING... ]' : '[ ACCESS SYSTEM ]'}
          </Button>
        </form>

        <div className="text-center text-xs text-muted-foreground font-mono space-y-1">
          <div>Run the following command from the terminal to generate auth key:</div>
          <code className="block bg-muted px-2 py-1 rounded text-[10px]">
            docker exec -it headscale headscale apikeys create --expiration 3h
          </code>
        </div>
      </div>
    </div>
  )
}

