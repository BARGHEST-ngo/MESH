import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import apiClient from './client'
import type { V1Node, V1ListNodesResponse, V1DeleteNodeResponse, V1ExpireNodeResponse, V1SetApprovedRoutesResponse } from './openapi/types.gen'
import { syncNetworkIsolationPolicy } from './useNetworks'

export const EXIT_ROUTES = ['0.0.0.0/0', '::/0']

// Node is advertising exit node capability (client-side)
export function isAdvertisingExitNode(node: V1Node): boolean {
    return EXIT_ROUTES.every(r => node.availableRoutes?.includes(r))
}

/// Admin has approved this node as an exit node
export function isApprovedExitNode(node: V1Node): boolean {
    return EXIT_ROUTES.every(r => node.approvedRoutes?.includes(r))
}

// List all nodes across all networks
export function useNodes() {
    return useQuery({
        queryKey: ['nodes'],
        queryFn: async () => {
            const response = await apiClient.get<V1ListNodesResponse>('/node')
            return response.data.nodes || []
        },
        refetchInterval: 10000,
    })
}

export function useNodesByNetwork(userName?: string) {
    return useQuery({
        queryKey: ['nodes', userName],
        queryFn: async () => {
            if (!userName) return []

            const response = await apiClient.get<V1ListNodesResponse>('/node')

            const filteredNodes = response.data.nodes?.filter(node =>
                node.preAuthKey?.user?.name === userName
            ) || []

            return filteredNodes
        },
        enabled: !!userName,
        refetchInterval: 10000,
    })
}

export function useDeleteNode() {
    const queryClient = useQueryClient()

    return useMutation({
        mutationFn: async (nodeId: string) => {
            const response = await apiClient.delete<V1DeleteNodeResponse>(`/node/${nodeId}`)
            return response.data
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['nodes'] })
        },
    })
}

export function useExpireNode() {
    const queryClient = useQueryClient()

    return useMutation({
        mutationFn: async (nodeId: string) => {
            const response = await apiClient.post<V1ExpireNodeResponse>(`/node/${nodeId}/expire`)
            return response.data
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['nodes'] })
        },
    })
}

export function useSetNodeTags() {
    const queryClient = useQueryClient()

    return useMutation({
        mutationFn: async (data: { nodeId: string; tags: string[] }) => {
            const response = await apiClient.post(`/node/${data.nodeId}/tags`, {
                tags: data.tags,
            })
            return response.data
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['nodes'] })
        },
    })
}

export function useSetApprovedRoutes() {
    const queryClient = useQueryClient()

    return useMutation({
        mutationFn: async (data: { nodeId: string; routes: string[] }) => {
            const response = await apiClient.post<V1SetApprovedRoutesResponse>(
                `/node/${data.nodeId}/approve_routes`,
                { routes: data.routes }
            )
            return response.data
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['nodes'] })
        },
    })
}

export function useApproveExitNode() {
    const setApprovedRoutes = useSetApprovedRoutes()
    const queryClient = useQueryClient()

    return useMutation({
        mutationFn: async (node: V1Node) => {
            if (!node.id) throw new Error('Node ID is required')
            const existingApproved = node.approvedRoutes || []
            const merged = Array.from(new Set([...existingApproved, ...EXIT_ROUTES]))
            const result = await setApprovedRoutes.mutateAsync({ nodeId: node.id, routes: merged })

            // Sync ACL policy to add autogroup:internet rule for this network
            try {
                await syncNetworkIsolationPolicy()
            } catch (aclError) {
                console.error('Failed to sync ACL policy after exit node approval:', aclError)
            }

            return result
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['policy'] })
        },
    })
}

export function useRevokeExitNode() {
    const setApprovedRoutes = useSetApprovedRoutes()
    const queryClient = useQueryClient()

    return useMutation({
        mutationFn: async (node: V1Node) => {
            if (!node.id) throw new Error('Node ID is required')
            const existingApproved = node.approvedRoutes || []
            const withoutExit = existingApproved.filter(r => !EXIT_ROUTES.includes(r))
            const result = await setApprovedRoutes.mutateAsync({ nodeId: node.id, routes: withoutExit })

            // Sync ACL policy to remove autogroup:internet rule if no exit nodes remain
            try {
                await syncNetworkIsolationPolicy()
            } catch (aclError) {
                console.error('Failed to sync ACL policy after exit node revocation:', aclError)
            }

            return result
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['policy'] })
        },
    })
}

