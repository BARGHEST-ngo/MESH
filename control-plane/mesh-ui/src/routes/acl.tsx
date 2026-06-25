import { createFileRoute } from '@tanstack/react-router'
import { useState, useEffect } from 'react'
import { RefreshCw, Save, Shield, Network } from 'lucide-react'
import { usePolicy, useSetPolicy } from '../api/usePolicy'
import { parseHuJSON, toHuJSON } from '../api/usePolicy'
import { Button } from '../components/ui/button'
import { Card } from '../components/ui/card'
import type { V1SetPolicyResponse } from '../api/openapi/types.gen'
import apiClient from '../api/client'
import { isAxiosError } from 'axios'

export const Route = createFileRoute('/acl')({
  component: ACLPage,
})

//(ovi)TODO: we need to handle when a policy hasn't been created
// currently it will display loading policy... but because policy hasn't been created
// user should be displayed a better message
function ACLPage() {
  const { data: policyData, isLoading, refetch } = usePolicy()
  const setPolicy = useSetPolicy()
  const [editedPolicy, setEditedPolicy] = useState<string>('')
  const [hasChanges, setHasChanges] = useState(false)
  const [malformedError, setMalformedError] = useState<string | null>(null)

  useEffect(() => {
    if (policyData && !hasChanges && !editedPolicy) {
      setEditedPolicy(policyData.raw)
    }
  }, [policyData, hasChanges, editedPolicy])

  const handlePolicyChange = (value: string) => {
    setEditedPolicy(value)
    setHasChanges(value !== policyData?.raw)
  }

  const handleSave = async () => {
    try {
      if (!editedPolicy || editedPolicy.trim() === '') {
        throw new Error('Policy cannot be empty')
      }
      const parsed = parseHuJSON(editedPolicy)
      //(Ovi) TODO: we should define a policy schema here to avoid arbitrary json inputs
      if (!parsed || typeof parsed !== 'object') {
        throw new Error('Invalid policy format')
      }
      await apiClient.put<V1SetPolicyResponse>('/policy', {
        policy: toHuJSON(parsed),
      })
      setHasChanges(false)
      setMalformedError(null)
    } catch (error) {
      console.error('Failed to save policy:', error)
      if (isAxiosError(error)) {
        if (error.response?.status === 500) {
          const message = error.response?.data?.message || error.response?.data?.error || error.message || 'Server error occurred'
          setMalformedError(message)
        }
      } else if (error instanceof Error) {
        setMalformedError(error.message)
      }
    }
  }

  const handleRefresh = () => {
    refetch()
    if (policyData) {
      setEditedPolicy(policyData.raw)
      setHasChanges(false)
    }
  }

  const lineCount = Math.max(editedPolicy.split('\n').length, 1)

  return (
    <div className="px-8 py-7 max-w-[1100px] mx-auto">
      <div className="mb-6">
        <h1 className="text-[26px] font-bold text-foreground tracking-tight mb-1.5">Access policy</h1>
        <p className="text-sm text-text2 max-w-[640px] leading-relaxed">
          Define who can reach whom across your networks. Your custom rules are preserved automatically
          when networks are created or updated.
        </p>
      </div>

      <div className="flex items-center gap-3 mb-3.5">
        <Button variant="secondary" onClick={handleRefresh} disabled={isLoading}>
          <RefreshCw size={17} className={isLoading ? 'animate-spin' : ''} />
          Refresh
        </Button>
        <Button onClick={handleSave} disabled={!hasChanges || setPolicy.isPending}>
          <Save size={17} />
          {setPolicy.isPending ? 'Saving…' : 'Save changes'}
        </Button>
        {hasChanges && (
          <span className="inline-flex items-center gap-2 text-[13px] font-semibold text-warning">
            <span className="w-[7px] h-[7px] rounded-full bg-warning" />
            Unsaved changes
          </span>
        )}
      </div>

      <Card className="overflow-hidden p-0">
        <div className="flex items-center gap-2 px-3.5 py-2.5 border-b border-border bg-card-hi">
          <Shield size={16} className="text-primary" />
          <span className="font-mono text-[12.5px] text-text2">policy.hujson</span>
        </div>
        <div className="flex max-h-[460px] overflow-auto">
          <div className="py-4 pl-4 pr-2.5 text-right select-none bg-card border-r border-border">
            {Array.from({ length: lineCount }, (_, i) => (
              <div key={i} className="font-mono text-[12.5px] leading-5 text-soft">{i + 1}</div>
            ))}
          </div>
          <textarea
            value={editedPolicy}
            onChange={(e) => handlePolicyChange(e.target.value)}
            className="flex-1 min-h-[420px] p-4 bg-card text-text2 font-mono text-[12.5px] leading-5 border-none outline-none resize-y box-border"
            placeholder="Loading policy…"
            disabled={isLoading}
            spellCheck={false}
          />
        </div>
      </Card>

      {setPolicy.isError && (
        <div className="mt-4 p-4 bg-destructive-dim border border-destructive/30 rounded-xl">
          <p className="text-destructive font-semibold">Error saving policy</p>
          <p className="text-sm text-text2 mt-1">
            {setPolicy.error instanceof Error ? setPolicy.error.message : 'Unknown error occurred'}
          </p>
        </div>
      )}

      {malformedError && (
        <div className="mt-4 p-4 bg-destructive-dim border border-destructive/30 rounded-xl">
          <p className="text-destructive font-semibold">Malformed policy</p>
          <p className="text-sm text-text2 mt-1">{malformedError}</p>
        </div>
      )}

      <Card className="mt-[18px] p-[18px] bg-bg-raised">
        <div className="flex gap-3">
          <div className="w-[34px] h-[34px] rounded-[9px] bg-primary/15 flex items-center justify-center shrink-0">
            <Network size={18} className="text-primary" />
          </div>
          <div>
            <div className="text-[14.5px] font-semibold text-foreground mb-1">Smart merge keeps your rules safe</div>
            <p className="text-[13px] text-text2 leading-relaxed max-w-[720px]">
              Each node gets a per-network tag automatically, so network isolation always works, even with
              custom tags. You can safely add cross-network or port rules here; they're never overwritten.
              See the docs at <span className="text-primary">docs.meshforensics.org</span>.
            </p>
          </div>
        </div>
      </Card>
    </div>
  )
}
