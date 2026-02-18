import axios, { type AxiosInstance } from 'axios'

export const getBaseURL = () => '/api/v1'

const apiClient: AxiosInstance = axios.create({
    baseURL: getBaseURL(),
    timeout: 10000,
    headers: {
        'Content-Type': 'application/json',
    },
})

// Request interceptor to add auth key
apiClient.interceptors.request.use(
    (config) => {
        const authKey = localStorage.getItem('authKey')
        if (authKey) {
            config.headers.Authorization = `Bearer ${authKey}`
        }
        return config
    },
    (error) => {
        return Promise.reject(error)
    }
)

apiClient.interceptors.response.use(
    (response) => response,
    (error) => {
        if (error.response?.status === 401) {
            // Clear auth on 401
            localStorage.removeItem('authKey')
            window.location.href = '/login'
        }
        return Promise.reject(error)
    }
)

export default apiClient

