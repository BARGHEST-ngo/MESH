import { useQuery } from "@tanstack/react-query";

export function useControlPlaneUrl() {
    return useQuery({
        queryKey: ['config'],
        queryFn: async () => {
            const response = await fetch('/config.json')
            if (!response.ok) return ''
            const config = await response.json()
            return config.controlPlaneUrl ?? ''
        }
    })
}