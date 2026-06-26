import { createFileRoute, redirect } from '@tanstack/react-router'
import { useState, useMemo } from 'react'
import { Plus, Network, RefreshCw, Filter } from 'lucide-react'
import { Button } from '../components/ui/button'
import { Card } from '../components/ui/card'
import { useNetworks } from '../api/useNetworks'
import { NetworkCard } from '../components/NetworkCard'
import { CreateNetworkDialog } from '../components/CreateNetworkDialog'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '../components/ui/select'

export const Route = createFileRoute('/')({
  beforeLoad: ({ context }) => {
    if (!context.auth.isAuthenticated && !context.auth.isLoading) {
      throw redirect({ to: '/login' })
    }
  },
  component: Index,
})

function Index() {
  const [createDialogOpen, setCreateDialogOpen] = useState(false)
  const [selectedEmail, setSelectedEmail] = useState<string>('all')

  const { data: networks = [], isLoading: networksLoading, refetch } = useNetworks()

  const uniqueEmails = useMemo(() => {
    const emails = networks
      .map(n => n.email)
      .filter((email): email is string => Boolean(email))
    return Array.from(new Set(emails)).sort()
  }, [networks])

  const filteredNetworks = useMemo(() => {
    if (selectedEmail === 'all') {
      return networks
    }
    return networks.filter(n => n.email === selectedEmail)
  }, [networks, selectedEmail])

  const isLoading = networksLoading

  return (
    <div className="px-8 py-7 max-w-[1100px] mx-auto">
      <div className="flex items-start justify-between gap-4 mb-6">
        <div>
          <h1 className="text-[26px] font-bold text-foreground tracking-tight mb-1.5">Networks</h1>
          <p className="text-sm text-text2 leading-relaxed">
            Manage your forensic networks and the clients connected to them.
          </p>
        </div>
        <div className="flex gap-2.5 shrink-0">
          <Button variant="secondary" onClick={() => refetch()} disabled={isLoading}>
            <RefreshCw size={17} className={isLoading ? 'animate-spin' : ''} />
            Refresh
          </Button>
          <Button onClick={() => setCreateDialogOpen(true)}>
            <Plus size={17} />
            New network
          </Button>
        </div>
      </div>

      <Card className="p-[18px] mb-6">
        <div className="flex items-center gap-3.5">
          <div className="w-[42px] h-[42px] rounded-[11px] bg-primary/15 flex items-center justify-center shrink-0">
            <Network size={20} className="text-primary" />
          </div>
          <div>
            <div className="text-2xl font-bold text-foreground leading-none">
              {filteredNetworks.length}
              {selectedEmail !== 'all' && (
                <span className="text-base font-semibold text-soft"> / {networks.length}</span>
              )}
            </div>
            <div className="text-[12.5px] text-text2 mt-1">
              {selectedEmail !== 'all' ? 'Filtered networks' : 'Networks'}
            </div>
          </div>
        </div>
      </Card>

      <div className="flex items-center justify-between gap-4 mb-3.5">
        <h2 className="text-[15px] font-semibold text-foreground tracking-tight">All networks</h2>

        {uniqueEmails.length > 0 && (
          <div className="flex items-center gap-2.5">
            <Filter size={15} className="text-soft" />
            <span className="text-[13px] text-text2">Filter by analyst</span>
            <Select value={selectedEmail} onValueChange={setSelectedEmail}>
              <SelectTrigger className="w-[260px]">
                <SelectValue placeholder="Select analyst" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All analysts</SelectItem>
                {uniqueEmails.map((email) => (
                  <SelectItem key={email} value={email}>
                    <span className="font-mono text-[13px]">{email}</span>
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
        )}
      </div>

      {isLoading ? (
        <div className="flex flex-col items-center justify-center py-16 text-text2">
          <RefreshCw size={28} className="animate-spin mb-3" />
          Loading networks…
        </div>
      ) : networks.length === 0 ? (
        <div className="flex flex-col items-center text-center py-16 border border-dashed border-border rounded-[var(--radius-card)]">
          <div className="w-12 h-12 rounded-xl bg-inset flex items-center justify-center mb-4">
            <Network size={24} className="text-soft" />
          </div>
          <p className="text-foreground font-semibold mb-1">No networks yet</p>
          <p className="text-sm text-text2 mb-5 max-w-[320px]">
            Create your first network to start onboarding forensic devices.
          </p>
          <Button onClick={() => setCreateDialogOpen(true)}>
            <Plus size={17} />
            Create your first network
          </Button>
        </div>
      ) : filteredNetworks.length === 0 ? (
        <div className="flex flex-col items-center text-center py-16 border border-dashed border-border rounded-[var(--radius-card)]">
          <div className="w-12 h-12 rounded-xl bg-inset flex items-center justify-center mb-4">
            <Filter size={22} className="text-soft" />
          </div>
          <p className="text-text2 mb-5">
            No networks for <span className="font-mono font-semibold text-foreground">{selectedEmail}</span>
          </p>
          <Button onClick={() => setSelectedEmail('all')} variant="secondary">
            Clear filter
          </Button>
        </div>
      ) : (
        <div className="flex flex-col gap-3.5">
          {filteredNetworks.map((network) => (
            <NetworkCard key={network.id} network={network} />
          ))}
        </div>
      )}

      <CreateNetworkDialog
        open={createDialogOpen}
        onOpenChange={setCreateDialogOpen}
      />
    </div>
  )
}
