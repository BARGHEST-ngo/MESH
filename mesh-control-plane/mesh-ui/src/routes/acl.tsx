import { createFileRoute } from '@tanstack/react-router'
import { useState, useEffect } from 'react'
import { usePolicy, useSetPolicy } from '../api/usePolicy'
import { Button } from '../components/ui/button'
import { RefreshCw, Save } from 'lucide-react'
import { parseHuJSON, toHuJSON } from '../api/usePolicy'
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

  console.log('Current malformedError state:', malformedError)

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
      console.log('Error object:', error)
      if (isAxiosError(error)) {
        console.log('Is axios error, status:', error.response?.status)
        console.log('Response data:', error.response?.data)
        if (error.response?.status === 500) {
          const message = error.response?.data?.message || error.response?.data?.error || error.message || 'Server error occurred'
          console.log('Setting malformed error:', message)
          setMalformedError(message)
        }
      } else if (error instanceof Error) {
        console.log('Setting error message:', error.message)
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

  return (
    <div className="p-6 max-w-7xl mx-auto">
      <div className="mb-6">
        <h1 className="text-3xl font-bold mb-2 font-mono">&gt; CUSTOMIZE_ACL</h1>
        <p className="text-muted-foreground">
          Edit your ACL policy configuration. Custom rules will be preserved when networks are created or updated.
        </p>
      </div>

      <div className="flex gap-3 mb-4">
        <Button
          onClick={handleRefresh}
          disabled={isLoading}
          variant="outline"
          className="font-mono"
        >
          <RefreshCw size={16} className={isLoading ? 'animate-spin' : ''} />
          <span>REFRESH</span>
        </Button>
        <Button
          onClick={handleSave}
          disabled={!hasChanges || setPolicy.isPending}
          className="font-mono"
        >
          <Save size={16} />
          <span>{setPolicy.isPending ? 'SAVING...' : 'SAVE CHANGES'}</span>
        </Button>
        {hasChanges && (
          <span className="text-sm text-amber-600 dark:text-amber-400 flex items-center">
            !! Unsaved changes
          </span>
        )}
      </div>

      <div className="border border-border rounded-lg overflow-hidden">
        <textarea
          value={editedPolicy}
          onChange={(e) => handlePolicyChange(e.target.value)}
          className="w-full h-[600px] p-4 font-mono text-sm bg-card text-foreground resize-none focus:outline-none focus:ring-2 focus:ring-ring"
          placeholder="Loading policy..."
          disabled={isLoading}
          spellCheck={false}
        />
      </div>

      {setPolicy.isError && (
        <div className="mt-4 p-4 bg-destructive/10 border border-destructive rounded-lg">
          <p className="text-destructive font-semibold">Error saving policy</p>
          <p className="text-sm text-muted-foreground mt-1">
            {setPolicy.error instanceof Error ? setPolicy.error.message : 'Unknown error occurred'}
          </p>
        </div>
      )}

      {malformedError && (
        <div className="mt-4 p-4 bg-destructive/10 border border-destructive rounded-lg">
          <p className="text-destructive font-semibold">Malformed Policy</p>
          <p className="text-sm text-muted-foreground mt-1">
            {malformedError}
          </p>
        </div>
      )}

      <div className="mt-6 p-4 bg-muted/50 rounded-lg">
        <h2 className="font-semibold mb-2 font-mono">&gt; SMART_MERGE_INFO</h2>
        <p className="text-sm text-muted-foreground">
          When you create or update networks, the system automatically regenerates network isolation rules 
	          (pattern: <code className="bg-background px-1 py-0.5 rounded">tag:net-&lt;network&gt; → tag:net-&lt;network&gt;:*</code>) 
	          while preserving all your custom ACL rules and tag ownership.
        </p>
        <p className="text-sm text-muted-foreground mt-2">
	          Each node registered to a network automatically receives a per-network tag (in addition to its role tag),
	          so isolation works correctly even when devices are tagged.
	          <br />
          This means you can safely add custom rules here (like cross-network access or specific IP/port restrictions) 
          and they won't be lost when networks are created or updated.
          If you have any questions regarding this, refer to the docs at https://docs.meshforensics.org
        </p>
      </div>
    </div>
  )
}

