import { useState } from 'react'
import { Card, CardHeader, CardTitle, CardContent } from './ui/card'
import { Button } from './ui/button'
import { Badge } from './ui/badge'
import { Trash2, Key, ChevronDown, ChevronRight, Wifi, WifiOff, Clock, Trash, Edit, Globe, ShieldOff } from 'lucide-react'
import { useNodesByNetwork } from '../api/useNodes'
import { useDeleteNetwork } from '../api/useNetworks'
import { useDeleteNode, useExpireNode, useApproveExitNode, useRevokeExitNode, isAdvertisingExitNode, isApprovedExitNode } from '../api/useNodes'
import { GenerateKeyDialog } from './GenerateKeyDialog'
import { RenameNetworkDialog } from './RenameNetworkDialog'
import type { V1User, V1Node } from '../api/openapi/types.gen'

interface NetworkCardProps {
  network: V1User
}

export function NetworkCard({ network }: NetworkCardProps) {
  const [expanded, setExpanded] = useState(false)
  const [keyDialogOpen, setKeyDialogOpen] = useState(false)
  const [renameDialogOpen, setRenameDialogOpen] = useState(false)
  
  const { data: nodes = [], isLoading: nodesLoading } = useNodesByNetwork(expanded ? network.name : undefined)
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
        await deleteNetwork.mutateAsync(network)
      } catch (error) {
        console.error('Failed to delete network:', error)
        alert('Failed to delete network. Make sure all clients are removed first.')
      }
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
      <Card className="hover:border-primary/50 transition-colors">
        <CardHeader>
          <div className="flex items-start justify-between">
            <div className="flex-1">
              <button
                onClick={() => setExpanded(!expanded)}
                className="flex items-center gap-2 hover:text-primary transition-colors group w-full text-left"
              >
                {expanded ? (
                  <ChevronDown size={20} className="text-primary" />
                ) : (
                  <ChevronRight size={20} className="group-hover:text-primary" />
                )}
                <div>
                  <CardTitle>
                    {network.displayName || network.name}
                  </CardTitle>
                  {network.displayName && (
                    <p className="text-xs text-muted-foreground font-mono mt-1">
                      {network.name}
                    </p>
                  )}
                </div>
              </button>
              
              <div className="flex items-center gap-2 mt-3 ml-7">
                <Badge variant={onlineCount > 0 ? "success" : "secondary"}>
                  {onlineCount} / {totalCount} ONLINE
                </Badge>
                {network.email && (
                  <span className="text-xs text-muted-foreground">{network.email}</span>
                )}
              </div>
            </div>

            <div className="flex gap-2">
              <Button
                size="sm"
                variant="outline"
                onClick={() => setKeyDialogOpen(true)}
                className="flex items-center gap-2"
              >
                <Key size={16} />
                [ KEY ]
              </Button>
              <Button
                size="sm"
                variant="outline"
                onClick={() => setRenameDialogOpen(true)}
                className="flex items-center gap-2"
              >
                <Edit size={16} />
                [ RENAME ]
              </Button>
              <Button
                size="sm"
                variant="destructive"
                onClick={handleDeleteNetwork}
                disabled={deleteNetwork.isPending}
                className="flex items-center gap-2"
              >
                <Trash2 size={16} />
                [ DELETE ]
              </Button>
            </div>
          </div>
        </CardHeader>

        {expanded && (
          <CardContent>
            {nodesLoading ? (
              <div className="text-center py-8 text-muted-foreground">
                Loading clients...
              </div>
            ) : nodes.length === 0 ? (
              <div className="text-center py-8 border-2 border-dashed border-border rounded-lg">
                <p className="text-muted-foreground mb-2">No clients registered</p>
                <Button size="sm" onClick={() => setKeyDialogOpen(true)}>
                  [ GENERATE KEY TO ADD CLIENTS ]
                </Button>
              </div>
            ) : (
              <div className="space-y-2">
                <h4 className="text-sm font-semibold text-foreground mb-3 font-mono">
                  &gt; CLIENTS
                </h4>
                {nodes.map((node) => (
                  <div
                    key={node.id}
                    className="flex items-center justify-between p-3 bg-secondary/50 border border-border rounded hover:border-primary/30 transition-colors"
                  >
                    <div className="flex items-center gap-3 flex-1">
                      {node.online ? (
                        <Wifi size={18} className="text-green-400" />
                      ) : (
                        <WifiOff size={18} className="text-muted-foreground" />
                      )}
                      <div className="flex-1">
                        <div className="flex items-center gap-2">
                          <span className="font-semibold text-foreground">
                            {node.givenName || node.name}
                          </span>
                          {node.online ? (
                            <Badge variant="success">ONLINE</Badge>
                          ) : (
                            <Badge variant="secondary">OFFLINE</Badge>
                          )}
                          {node.validTags && node.validTags.length > 0 && (
                            node.validTags.map((tag) => {
                              const tagName = tag.replace('tag:', '').toUpperCase()
                              const isAnalyst = tagName === 'ANALYST'
                              return (
                                <Badge
                                  key={tag}
                                  variant={isAnalyst ? "default" : "destructive"}
                                  className="text-xs"
                                >
                                  {tagName}
                                </Badge>
                              )
                            })
                          )}
                          {isApprovedExitNode(node) && (
                            <Badge variant="warning" className="text-xs">
                              <Globe size={12} className="mr-1" />
                              EXIT NODE
                            </Badge>
                          )}
                        </div>
                        <div className="text-xs text-muted-foreground mt-1 space-y-0.5">
                          <div className="font-mono">
                            {node.ipAddresses?.join(', ') || 'No IP'}
                          </div>
                          {node.lastSeen && (
                            <div>
                              Last seen: {new Date(node.lastSeen).toLocaleString()}
                            </div>
                          )}
                          {isAdvertisingExitNode(node) && !isApprovedExitNode(node) && (
                            <div className="text-yellow-400 font-semibold">
                              [ ADVERTISING EXIT NODE — NOT APPROVED ]
                            </div>
                          )}
                        </div>
                      </div>
                    </div>

                    <div className="flex gap-2">
                      {isAdvertisingExitNode(node) && !isApprovedExitNode(node) && (
                        <Button
                          size="sm"
                          variant="outline"
                          onClick={() => handleApproveExitNode(node)}
                          disabled={approveExitNode.isPending}
                          className="flex items-center gap-1 border-yellow-500/50 text-yellow-400 hover:bg-yellow-500/10"
                        >
                          <Globe size={14} />
                          APPROVE EXIT
                        </Button>
                      )}
                      {isApprovedExitNode(node) && (
                        <Button
                          size="sm"
                          variant="outline"
                          onClick={() => handleRevokeExitNode(node)}
                          disabled={revokeExitNode.isPending}
                          className="flex items-center gap-1"
                        >
                          <ShieldOff size={14} />
                          REVOKE EXIT
                        </Button>
                      )}
                      <Button
                        size="sm"
                        variant="outline"
                        onClick={() => node.id && handleExpireNode(node.id, node.name || 'Unknown')}
                        disabled={expireNode.isPending}
                        className="flex items-center gap-1"
                      >
                        <Clock size={14} />
                        EXPIRE
                      </Button>
                      <Button
                        size="sm"
                        variant="destructive"
                        onClick={() => node.id && handleDeleteNode(node.id, node.name || 'Unknown')}
                        disabled={deleteNode.isPending}
                        className="flex items-center gap-1"
                      >
                        <Trash size={14} />
                      </Button>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </CardContent>
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

