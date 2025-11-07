/**
 * API Key Management Hooks
 * 
 * TODO: Implement key management page
 */

import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import apiClient from './client'
import type {
    V1CreateApiKeyRequest,
    V1CreateApiKeyResponse,
    V1ListApiKeysResponse,
    V1ExpireApiKeyRequest,
    V1ExpireApiKeyResponse,
} from './openapi/types.gen'

export function useApiKeys() {
    return useQuery({
        queryKey: ['apiKeys'],
        queryFn: async () => {
            const response = await apiClient.get<V1ListApiKeysResponse>('/apikey')
            return response.data.apiKeys || []
        },
    })
}

export function useCreateApiKey() {
    const queryClient = useQueryClient()

    return useMutation({
        mutationFn: async (data: {
            expiration?: string
        }) => {
            const request: V1CreateApiKeyRequest = {
                expiration: data.expiration,
            }
            const response = await apiClient.post<V1CreateApiKeyResponse>('/apikey', request)
            return response.data
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['apiKeys'] })
        },
    })
}

export function useExpireApiKey() {
    const queryClient = useQueryClient()

    return useMutation({
        mutationFn: async (data: { prefix: string }) => {
            const request: V1ExpireApiKeyRequest = {
                prefix: data.prefix,
            }
            const response = await apiClient.post<V1ExpireApiKeyResponse>('/apikey/expire', request)
            return response.data
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['apiKeys'] })
        },
    })
}

export function useDeleteApiKey() {
    const queryClient = useQueryClient()

    return useMutation({
        mutationFn: async (prefix: string) => {
            const response = await apiClient.delete(`/apikey/${prefix}`)
            return response.data
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['apiKeys'] })
        },
    })
}
