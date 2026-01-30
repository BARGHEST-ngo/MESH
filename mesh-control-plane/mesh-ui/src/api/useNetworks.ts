import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import apiClient from './client'
import type { V1User, V1CreateUserRequest, V1ListUsersResponse, V1DeleteUserResponse } from './openapi/types.gen'
import type { V1GetPolicyResponse, V1SetPolicyResponse } from './openapi/types.gen'
import { parseHuJSON, toHuJSON } from './usePolicy'
import { generateNetworkIsolationPolicy } from '../lib/aclPolicyGenerator'

export function useNetworks() {
    return useQuery({
        queryKey: ['networks'],
        queryFn: async () => {
            const response = await apiClient.get<V1ListUsersResponse>('/user')
            return response.data.users || []
        },
    })
}

// TODO impliment timestamping since multiple users can overlap
async function syncNetworkIsolationPolicy(): Promise<void> {
    // Fetch current networks
    const networksResponse = await apiClient.get<V1ListUsersResponse>('/user')
    const networks = networksResponse.data.users || []
    const currentPolicyResponse = await apiClient.get<V1GetPolicyResponse>('/policy')
    const currentRaw = currentPolicyResponse.data.policy || '{}'
    const currentPolicy = parseHuJSON(currentRaw) as Record<string, unknown>
    const generatedPolicy = generateNetworkIsolationPolicy(networks)
    const mergedPolicy = {
        ...currentPolicy,
        ...generatedPolicy,
        // Ensure we always overwrite ACL rules with the generated ones TODO: test usecases for why you wouldnt want this
        acls: generatedPolicy.acls,
    }

    //apply
    await apiClient.put<V1SetPolicyResponse>('/policy', {
        policy: toHuJSON(mergedPolicy),
    })
}

export function useCreateNetwork() {
    const queryClient = useQueryClient()

    return useMutation({
        mutationFn: async (data: { name: string; displayName?: string; email?: string }) => {
            const request: V1CreateUserRequest = {
                name: data.name,
                displayName: data.displayName,
                email: data.email,
            }

            // Create the network first
            const response = await apiClient.post('/user', request)

            try {
                await syncNetworkIsolationPolicy()
            } catch (aclError) {
                console.error('Failed to sync ACL policy after network creation:', aclError)
            }

            return response.data
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['networks'] })
            queryClient.invalidateQueries({ queryKey: ['policy'] })
        },
    })
}

export function useRenameNetwork() {
    const queryClient = useQueryClient()

    return useMutation({
        mutationFn: async (data: { user: V1User; newName: string }) => {
            if (!data.user.id) {
                throw new Error('User ID is required')
            }

            const response = await apiClient.post(`/user/${data.user.id}/rename/${data.newName}`)

           
            try {
                await syncNetworkIsolationPolicy()
            } catch (aclError) {
                console.error('Failed to sync ACL policy after network rename:', aclError)
            }

            return response.data
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['networks'] })
            queryClient.invalidateQueries({ queryKey: ['nodes'] })
            queryClient.invalidateQueries({ queryKey: ['policy'] })
        },
    })
}

export function useDeleteNetwork() {
    const queryClient = useQueryClient()

    return useMutation({
        mutationFn: async (user: V1User) => {
            if (!user.id) {
                throw new Error('User ID is required')
            }

            // Delete the network first
            const response = await apiClient.delete<V1DeleteUserResponse>(`/user/${user.id}`)

            // Sync ACL policy to remove the deleted network's rule
            // This regenerates the entire policy to ensure consistency
            try {
                await syncNetworkIsolationPolicy()
            } catch (aclError) {
                console.error('Failed to sync ACL policy after network deletion:', aclError)
            }

            return response.data
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['networks'] })
            queryClient.invalidateQueries({ queryKey: ['nodes'] })
            queryClient.invalidateQueries({ queryKey: ['policy'] })
        },
    })
}

export function useSyncNetworkIsolation() {
    const queryClient = useQueryClient()

    return useMutation({
        mutationFn: async () => {
            await syncNetworkIsolationPolicy()
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['policy'] })
        },
    })
}

