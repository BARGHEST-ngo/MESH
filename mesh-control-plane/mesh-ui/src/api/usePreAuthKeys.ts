/**
 * Pre-Auth Key Management Hooks
 * 
 * Pre-auth keys are used to register nodes to a network.
 */

import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import apiClient from './client'
import type {
    V1CreatePreAuthKeyRequest,
    V1CreatePreAuthKeyResponse,
    V1ListPreAuthKeysResponse,
    V1ExpirePreAuthKeyRequest,
    V1ExpirePreAuthKeyResponse,
} from './openapi/types.gen'

import { networkTagForNetwork } from '../lib/aclPolicyGenerator'

export function usePreAuthKeys(userName?: string) {
    return useQuery({
        queryKey: ['preAuthKeys', userName],
        queryFn: async () => {
            if (!userName) return []

            const usersResponse = await apiClient.get('/user')
            const users = usersResponse.data.users || []
            const user = users.find((u: any) => u.name === userName)

            if (!user) {
                throw new Error(`User '${userName}' not found`)
            }

            const response = await apiClient.get<V1ListPreAuthKeysResponse>(`/preauthkey`, {
                params: { user: user.id }
            })
            return response.data.preAuthKeys || []
        },
        enabled: !!userName,
    })
}

export function useCreatePreAuthKey() {
    const queryClient = useQueryClient()

    return useMutation({
        mutationFn: async (data: {
            user: string
            reusable?: boolean
            ephemeral?: boolean
            expiration?: string
            aclTags?: string[]
        }) => {
            const usersResponse = await apiClient.get('/user')
            const users = usersResponse.data.users || []
            const user = users.find((u: any) => u.name === data.user)

            if (!user) {
                throw new Error(`User '${data.user}' not found`)
            }

			// Ensure every key includes the per-network tag (in addition to any role tags).
			const requestedTags = data.aclTags || []
			const networkTag = networkTagForNetwork(user)
			const aclTags = Array.from(new Set([...(requestedTags || []), ...(networkTag ? [networkTag] : [])]))
				.filter((t): t is string => typeof t === 'string' && t.length > 0)

            const request: V1CreatePreAuthKeyRequest = {
                user: user.id,
                reusable: data.reusable || false,
                ephemeral: data.ephemeral || false,
                expiration: data.expiration,
				aclTags,
            }

            const response = await apiClient.post<V1CreatePreAuthKeyResponse>('/preauthkey', request)
            return response.data
        },
        onSuccess: (_, variables) => {
            queryClient.invalidateQueries({ queryKey: ['preAuthKeys', variables.user] })
        },
    })
}

export function useExpirePreAuthKey() {
    const queryClient = useQueryClient()

    return useMutation({
        mutationFn: async (data: { user: string; key: string }) => {
            const usersResponse = await apiClient.get('/user')
            const users = usersResponse.data.users || []
            const user = users.find((u: any) => u.name === data.user)

            if (!user) {
                throw new Error(`User '${data.user}' not found`)
            }

            const request: V1ExpirePreAuthKeyRequest = {
                user: user.id, // Use the user ID instead of name
                key: data.key,
            }
            const response = await apiClient.post<V1ExpirePreAuthKeyResponse>('/preauthkey/expire', request)
            return response.data
        },
        onSuccess: (_, variables) => {
            queryClient.invalidateQueries({ queryKey: ['preAuthKeys', variables.user] })
        },
    })
}

