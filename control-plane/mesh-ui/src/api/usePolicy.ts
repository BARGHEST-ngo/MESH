import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import apiClient from './client'
import type { V1GetPolicyResponse, V1SetPolicyResponse } from './openapi/types.gen'


export interface ACLRule {
    action: 'accept'
    src: string[]
    dst: string[]
    proto?: 'tcp' | 'udp' | 'icmp'
}

export interface ACLPolicy {
    groups?: Record<string, string[]>
    tagOwners?: Record<string, string[]>
    hosts?: Record<string, string>
    acls?: ACLRule[]
    tests?: ACLTest[]
}

export interface ACLTest {
    src: string
    accept?: string[]
    deny?: string[]
}

export function parseHuJSON(huJson: string): ACLPolicy {
    if (!huJson || huJson.trim() === '') {
        return {}
    }
    
    try {
        let cleaned = huJson.replace(/\/\/.*$/gm, '')
        cleaned = cleaned.replace(/\/\*[\s\S]*?\*\//g, '')
        cleaned = cleaned.replace(/,(\s*[}\]])/g, '$1')
        
        return JSON.parse(cleaned)
    } catch (error) {
        console.error('Failed to parse HuJSON policy:', error)
        return {}
    }
}

export function toHuJSON(policy: ACLPolicy): string {
    // Use standard JSON with 2-space indentation
    return JSON.stringify(policy, null, 2)
}

export function usePolicy() {
    return useQuery({
        queryKey: ['policy'],
        queryFn: async () => {
            const response = await apiClient.get<V1GetPolicyResponse>('/policy')
            const raw = response.data.policy || '{}'
            return {
                raw,
                parsed: parseHuJSON(raw),
                updatedAt: response.data.updatedAt,
            }
        },
        // Don't refetch too frequently as policy changes are infrequent
        staleTime: 30000,
    })
}

export function useSetPolicy() {
    const queryClient = useQueryClient()

    return useMutation({
        mutationFn: async (policy: string) => {
            const response = await apiClient.put<V1SetPolicyResponse>('/policy', {
                policy,
            })
            return response.data
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['policy'] })
        },
    })
}

export function useSetPolicyObject() {
    const setPolicy = useSetPolicy()

    return useMutation({
        mutationFn: async (policy: ACLPolicy) => {
            const huJson = toHuJSON(policy)
            return setPolicy.mutateAsync(huJson)
        },
    })
}

