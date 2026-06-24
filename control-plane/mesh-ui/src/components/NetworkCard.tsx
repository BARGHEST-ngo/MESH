import { useState } from 'react'
import { Card } from './ui/card'
import { Button } from './ui/button'
import { Badge } from './ui/badge'
import { Key, ChevronDown, ChevronRight, Wifi, WifiOff, Clock, Trash, Edit, Mail, Globe, ShieldOff } from 'lucide-react'
import { useNodesByNetwork } from '../api/useNodes'
import { useDeleteNetwork } from '../api/useNetworks'
import { useDeleteNode, useExpireNode, useApproveExitNode, useRevokeExitNode, isAdvertisingExitNode, isApprovedExitNode } from '../api/useNodes'
import { GenerateKeyDialog } from './GenerateKeyDialog'
import { RenameNetworkDialog } from './RenameNetworkDialog'
import type { V1User, V1Node } from '../api/openapi/types.gen'

interface NetworkCardProps {
  network: V1User
}

const TAG_LABEL: Record<string, string> = { analyst: 'Analyst', mobile_node: 'Mobile node' }

// Role tags only; per-network "tag:net-..." tags are internal noise.
function roleTags(node: V1Node): string[] {
  return (node.tags || [])
    .filter(t => !t.startsWith('tag:net-'))
    .map(t => t.replace('tag:', ''))
}

export function NetworkCard({ network }: NetworkCardProps) {
  const [expanded, setExpanded] = useState(false)
  const [keyDialogOpen, setKeyDialogOpen] = useState(false)
  const [renameDialogOpen, setRenameDialogOpen] = useState(false)

  const { data: nodes = [], isLoading: nodesLoading } = useNodesByNetwork(network.name)
  const deleteNetwork = useDeleteNetwork()
  const deleteNode = useDeleteNode()
  const expireNode = useExpireNode()
  const approveExitNode = useApproveExitNode()
  const revokeExitNode = useRevokeExitNode()

  const handleDeleteNetwork = async () => {
    if (!network.name) return

    const clientCount = nodes.length
    const confirmMessage = clientCount > 0
      ? `Are you sure you want to delete "${network.name}"? This will remove ${clientCount} client(s).`
      : `Are you sure you want to delete "${network.name}"?`

    if (window.confirm(confirmMessage)) {
      try {
        if (clientCount > 0) {
          await handleDeleteAllNodes()
        }
        await deleteNetwork.mutateAsync(network)
      } catch (error) {
        console.error('Failed to delete network:', error)
        alert('Failed to delete network. Make sure all clients are removed first.')
      }
    }
  }

  const handleDeleteAllNodes = async () => {
    try {
      for (const node of nodes) {
        if (!node.id) continue
        await deleteNode.mutateAsync(node.id)
      }
    } catch (error) {
      console.error('Failed to delete node:', error)
    }
  }

  const handleDeleteNode = async (nodeId: string, nodeName: string) => {
    if (window.confirm(`Are you sure you want to delete client "${nodeName}"?`)) {
      try {
        await deleteNode.mutateAsync(nodeId)
      } catch (error) {
        console.error('Failed to delete node:', error)
      }
    }
  }

  const handleExpireNode = async (nodeId: string, nodeName: string) => {
    if (window.confirm(`Are you sure you want to expire client "${nodeName}"? It will need to re-authenticate.`)) {
      try {
        await expireNode.mutateAsync(nodeId)
      } catch (error) {
        console.error('Failed to expire node:', error)
      }
    }
  }

  const handleApproveExitNode = async (node: V1Node) => {
    const nodeName = node.givenName || node.name || 'Unknown'
    if (window.confirm(`Approve "${nodeName}" as an exit node? All traffic from other nodes using this exit node will route through it.`)) {
      try {
        await approveExitNode.mutateAsync(node)
      } catch (error) {
        console.error('Failed to approve exit node:', error)
      }
    }
  }

  const handleRevokeExitNode = async (node: V1Node) => {
    const nodeName = node.givenName || node.name || 'Unknown'
    if (window.confirm(`Revoke exit node status from "${nodeName}"?`)) {
      try {
        await revokeExitNode.mutateAsync(node)
      } catch (error) {
        console.error('Failed to revoke exit node:', error)
      }
    }
  }

  const onlineCount = nodes.filter(n => n.online).length
  const totalCount = nodes.length

  return (
    <>
      <Card hover className="overflow-hidden">
        <div className="flex items-start justify-between gap-3.5 p-5">
          <button
            onClick={() => setExpanded(!expanded)}
            className="flex items-start gap-3 flex-1 text-left"
          >
            {expanded ? (
              <ChevronDown size={20} className="text-soft mt-0.5" />
            ) : (
              <ChevronRight size={20} className="text-soft mt-0.5" />
            )}
            <div>
              <div className="text-[17px] font-bold text-foreground tracking-tight">
                {network.displayName || network.name}
              </div>
              {network.displayName && (
                <div className="text-[11.5px] text-soft font-mono mt-0.5">{network.name}</div>
              )}
              <div className="flex items-center gap-2 mt-2.5">
                <Badge variant={onlineCount > 0 ? 'success' : 'secondary'}>
                  <span className="relative inline-flex w-1.5 h-1.5">
                    {onlineCount > 0 && (
                      <span className="absolute inset-0 rounded-full bg-success motion-safe:[animation:cpPulse_1.8s_ease-out_infinite]" />
                    )}
                    <span
                      className={`w-1.5 h-1.5 rounded-full relative ${onlineCount > 0 ? 'bg-success' : 'bg-soft'}`}
                    />
                  </span>
                  {onlineCount} / {totalCount} online
                </Badge>
                {network.email && (
                  <span className="inline-flex items-center gap-1.5 text-[12.5px] text-text2">
                    <Mail size={13} className="text-soft" />
                    {network.email}
                  </span>
                )}
              </div>
            </div>
          </button>

          <div className="flex gap-2 shrink-0">
            <Button size="sm" variant="outline" onClick={() => setKeyDialogOpen(true)}>
              <Key size={15} />
              Key
            </Button>
            <Button size="sm" variant="outline" onClick={() => setRenameDialogOpen(true)}>
              <Edit size={15} />
              Rename
            </Button>
            <Button
              size="sm"
              variant="destructive"
              onClick={handleDeleteNetwork}
              disabled={deleteNetwork.isPending}
              aria-label="Delete network"
            >
              <Trash size={15} />
            </Button>
          </div>
        </div>

        {expanded && (
          <div className="px-5 pb-[18px] motion-safe:[animation:cpFade_.2s_ease]">
            {nodesLoading ? (
              <div className="text-center py-8 text-text2">Loading clients…</div>
            ) : nodes.length === 0 ? (
              <div className="text-center py-8 border border-dashed border-border rounded-xl">
                <p className="text-text2 mb-3">No clients registered</p>
                <Button size="sm" onClick={() => setKeyDialogOpen(true)}>
                  <Key size={15} />
                  Generate a key to add clients
                </Button>
              </div>
            ) : (
              <>
                <div className="text-[12.5px] font-semibold text-text2 mt-1 mb-2.5 pl-0.5">Clients</div>
                <div className="flex flex-col gap-2">
                  {nodes.map((node) => {
                    const advertising = isAdvertisingExitNode(node) && !isApprovedExitNode(node)
                    return (
                      <div
                        key={node.id}
                        className="flex items-center gap-3.5 px-3.5 py-3 bg-bg-raised border border-border rounded-xl"
                      >
                        <div
                          className={`w-9 h-9 rounded-[10px] shrink-0 flex items-center justify-center ${
                            node.online ? 'bg-success-dim' : 'bg-inset'
                          }`}
                        >
                          {node.online ? (
                            <Wifi size={18} className="text-success" />
                          ) : (
                            <WifiOff size={18} className="text-soft" />
                          )}
                        </div>

                        <div className="flex-1 min-w-0">
                          <div className="flex items-center gap-2 flex-wrap">
                            <span className="text-sm font-semibold text-foreground">
                              {node.givenName || node.name}
                            </span>
                            <Badge variant={node.online ? 'success' : 'secondary'}>
                              {node.online ? 'Online' : 'Offline'}
                            </Badge>
                            {roleTags(node).map((tag) => (
                              <Badge key={tag} variant={tag === 'analyst' ? 'accent' : 'secondary'}>
                                {TAG_LABEL[tag] || tag}
                              </Badge>
                            ))}
                            {isApprovedExitNode(node) && (
                              <Badge variant="warning">
                                <Globe size={12} />
                                Exit node
                              </Badge>
                            )}
                          </div>
                          <div className="flex items-center gap-3 mt-1 flex-wrap">
                            <span className="font-mono text-[11.5px] text-text2">
                              {node.ipAddresses?.join(', ') || 'No IP'}
                            </span>
                            {node.lastSeen && (
                              <span className="inline-flex items-center gap-1 text-xs text-soft">
                                <Clock size={12} />
                                {new Date(node.lastSeen).toLocaleString()}
                              </span>
                            )}
                            {advertising && (
                              <span className="text-xs font-semibold text-warning">
                                Advertising exit node — not approved
                              </span>
                            )}
                          </div>
                        </div>

                        <div className="flex gap-1.5 shrink-0">
                          {advertising && (
                            <Button
                              size="sm"
                              variant="outline"
                              onClick={() => handleApproveExitNode(node)}
                              disabled={approveExitNode.isPending}
                              className="border-warning/50 text-warning hover:bg-warning-dim"
                            >
                              <Globe size={14} />
                              Approve exit
                            </Button>
                          )}
                          {isApprovedExitNode(node) && (
                            <Button
                              size="sm"
                              variant="outline"
                              onClick={() => handleRevokeExitNode(node)}
                              disabled={revokeExitNode.isPending}
                            >
                              <ShieldOff size={14} />
                              Revoke exit
                            </Button>
                          )}
                          <Button
                            size="sm"
                            variant="ghost"
                            onClick={() => node.id && handleExpireNode(node.id, node.name || 'Unknown')}
                            disabled={expireNode.isPending}
                          >
                            <Clock size={14} />
                            Expire
                          </Button>
                          <Button
                            size="sm"
                            variant="ghost"
                            onClick={() => node.id && handleDeleteNode(node.id, node.name || 'Unknown')}
                            disabled={deleteNode.isPending}
                            aria-label="Delete client"
                          >
                            <Trash size={14} />
                          </Button>
                        </div>
                      </div>
                    )
                  })}
                </div>
              </>
            )}
          </div>
        )}
      </Card>

      {network.name && (
        <>
          <GenerateKeyDialog
            open={keyDialogOpen}
            onOpenChange={setKeyDialogOpen}
            networkName={network.name}
          />
          <RenameNetworkDialog
            open={renameDialogOpen}
            onOpenChange={setRenameDialogOpen}
            network={network}
          />
        </>
      )}
    </>
  )
}
