import { useMutation } from '@tanstack/react-query'
import axios from 'axios'
import type { V1ListApiKeysResponse } from './openapi/types.gen'
import { getBaseURL } from './client'

interface LoginRequest {
    authKey: string
}

interface LoginResponse {
    success: boolean
    message?: string
}

const validateAuthKey = async (authKey: string): Promise<LoginResponse> => {
    try {
        // Test request to validate auth
        const response = await axios.get<V1ListApiKeysResponse>(
            `${getBaseURL()}/apikey`,
            {
                headers: {
                    Authorization: `Bearer ${authKey}`,
                },
                timeout: 5000,
            }
        )

        if (response.status === 200) {
            return {
                success: true,
                message: 'Authentication successful'
            }
        }

        throw new Error('Invalid authentication key')
    } catch (error) {
        if (axios.isAxiosError(error)) {
            if (error.response?.status === 401) {
                return {
                    success: false,
                    message: 'Invalid authentication key'
                }
            }
            return {
                success: false,
                message: error.response?.data?.message || 'Failed to connect to Headscale API'
            }
        }
        return {
            success: false,
            message: error instanceof Error ? error.message : 'Authentication failed'
        }
    }
}

export const useLogin = () => {
    return useMutation({
        mutationFn: async ({ authKey }: LoginRequest) => {
            const result = await validateAuthKey(authKey)

            if (result.success) {
                // Store auth key in localStorage
                localStorage.setItem('authKey', authKey)
                return result
            }

            throw new Error(result.message || 'Authentication failed')
        },
    })
}
