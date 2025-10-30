import { createFileRoute, redirect } from '@tanstack/react-router'
import { useState, useMemo } from 'react'
import { Plus, Network, RefreshCw, Filter } from 'lucide-react'
import { Button } from '../components/ui/button'
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
    <div className="90-vh bg-background">
      <div className="container mx-auto p-6 max-w-7xl">
        <div className="mb-8">
          <div className="flex items-center justify-between mb-4">
            <div>
              <h1 className="text-4xl font-bold text-foreground font-mono mb-2">
                &gt; MESH_CONTROL
              </h1>
              <p className="text-muted-foreground">
                Manage your networks and connected clients
              </p>
            </div>
            <div className="flex gap-3">
              <Button
                variant="outline"
                onClick={() => refetch()}
                disabled={isLoading}
                className="flex items-center gap-2"
              >
                <RefreshCw size={16} className={isLoading ? 'animate-spin' : ''} />
                [ REFRESH ]
              </Button>
              <Button
                onClick={() => setCreateDialogOpen(true)}
                className="flex items-center gap-2"
              >
                <Plus size={16} />
                [ NEW NETWORK ]
              </Button>
            </div>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div className="bg-card border border-border rounded-lg p-4">
              <div className="flex items-center gap-3">
                <Network size={24} className="text-primary" />
                <div>
                  <p className="text-2xl font-bold text-foreground">
                    {filteredNetworks.length}
                    {selectedEmail !== 'all' && (
                      <span className="text-base text-muted-foreground"> / {networks.length}</span>
                    )}
                  </p>
                  <p className="text-sm text-muted-foreground">
                    {selectedEmail !== 'all' ? 'Filtered Networks' : 'Networks'}
                  </p>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div className="space-y-4">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-2xl font-bold text-foreground font-mono">
              &gt; NETWORKS
            </h2>

            {uniqueEmails.length > 0 && (
              <div className="flex items-center gap-3">
                <Filter size={16} className="text-muted-foreground" />
                <span className="text-sm text-muted-foreground">Filter by analyst email:</span>
                <Select value={selectedEmail} onValueChange={setSelectedEmail}>
                  <SelectTrigger className="w-[280px]">
                    <SelectValue placeholder="Select email" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="all">
                      <span className="font-mono">[ ALL EMAILS ]</span>
                    </SelectItem>
                    {uniqueEmails.map((email) => (
                      <SelectItem key={email} value={email}>
                        <span className="font-mono">{email}</span>
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
            )}
          </div>

          {isLoading ? (
            <div className="text-center py-12 text-muted-foreground">
              <RefreshCw size={32} className="animate-spin mx-auto mb-4" />
              Loading networks...
            </div>
          ) : networks.length === 0 ? (
            <div className="text-center py-12 border-2 border-dashed border-border rounded-lg">
              <Network size={48} className="mx-auto mb-4 text-muted-foreground" />
              <p className="text-muted-foreground mb-4">No networks found</p>
              <Button onClick={() => setCreateDialogOpen(true)}>
                <Plus size={16} className="mr-2" />
                [ CREATE YOUR FIRST NETWORK ]
              </Button>
            </div>
          ) : filteredNetworks.length === 0 ? (
            <div className="text-center py-12 border-2 border-dashed border-border rounded-lg">
              <Network size={48} className="mx-auto mb-4 text-muted-foreground" />
              <p className="text-muted-foreground mb-4">
                No networks found for email: <span className="font-mono font-semibold">{selectedEmail}</span>
              </p>
              <Button onClick={() => setSelectedEmail('all')} variant="outline">
                [ CLEAR FILTER ]
              </Button>
            </div>
          ) : (
            <div className="space-y-4">
              {filteredNetworks.map((network) => (
                <NetworkCard key={network.id} network={network} />
              ))}
            </div>
          )}
        </div>
      </div>

      <CreateNetworkDialog
        open={createDialogOpen}
        onOpenChange={setCreateDialogOpen}
      />
    </div>
  )
}
