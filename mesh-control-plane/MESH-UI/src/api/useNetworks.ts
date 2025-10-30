/**
 * Network (User) Management Hooks
 * 
 * In Headscale, "Users" act as network namespaces.
 */

import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import apiClient from './client'
import type { V1User, V1CreateUserRequest, V1ListUsersResponse, V1DeleteUserResponse } from './openapi/types.gen'

export function useNetworks() {
    return useQuery({
        queryKey: ['networks'],
        queryFn: async () => {
            const response = await apiClient.get<V1ListUsersResponse>('/user')
            return response.data.users || []
        },
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
            const response = await apiClient.post('/user', request)
            return response.data
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['networks'] })
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
            return response.data
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['networks'] })
            queryClient.invalidateQueries({ queryKey: ['nodes'] })
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

            const response = await apiClient.delete<V1DeleteUserResponse>(`/user/${user.id}`)
            return response.data
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['networks'] })
            queryClient.invalidateQueries({ queryKey: ['nodes'] })
        },
    })
}

