import React, { createContext, useContext, useState, useEffect, type ReactNode } from 'react'

interface AuthContextType {
  isAuthenticated: boolean
  authKey: string
  login: (authKey: string) => Promise<boolean>
  logout: () => void
  isLoading: boolean
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export const useAuth = () => {
  const context = useContext(AuthContext)
  if (!context) throw new Error('useAuth must be used within an AuthProvider')
  return context
}

export const AuthProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [isAuthenticated, setIsAuthenticated] = useState(false)
  const [authKey, setAuthKey] = useState('')
  const [isLoading, setIsLoading] = useState(true)

  const clearAuth = () => {
    localStorage.removeItem('authKey')
    setAuthKey('')
    setIsAuthenticated(false)
  }

  useEffect(() => {
    const storedAuthKey = localStorage.getItem('authKey')
    
    if (storedAuthKey) {
      setAuthKey(storedAuthKey)
      setIsAuthenticated(true)
    }
    
    setIsLoading(false)
  }, [])

  const login = async (key: string) => {
    try {
      if (!key || key.trim().length === 0) {
        throw new Error('Invalid authentication key')
      }

      localStorage.setItem('authKey', key)
      setAuthKey(key)
      setIsAuthenticated(true)
      
      return true
    } catch (error) {
      clearAuth()
      return false
    }
  }

  return (
    <AuthContext.Provider value={{
      isAuthenticated,
      authKey,
      login,
      logout: clearAuth,
      isLoading
    }}>
      {children}
    </AuthContext.Provider>
  )
}
