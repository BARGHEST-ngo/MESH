/**
 * Node (Client/Device) Management Hooks
 * 
 * Nodes are the actual devices/clients connected to the network.
 */

import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import apiClient from './client'
import type { V1ListNodesResponse, V1DeleteNodeResponse, V1ExpireNodeResponse } from './openapi/types.gen'

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
                node.user?.name === userName
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

