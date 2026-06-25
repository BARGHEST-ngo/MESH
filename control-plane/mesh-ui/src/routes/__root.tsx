import { Outlet, createRootRouteWithContext, useRouterState } from '@tanstack/react-router'

import Sidebar from '../components/Sidebar'
import { useAuth } from '../lib/auth'

import type { QueryClient } from '@tanstack/react-query'

interface MyRouterContext {
  queryClient: QueryClient
  auth: {
    isAuthenticated: boolean
    isLoading: boolean
  }
}

function RootComponent() {
  const routerState = useRouterState()
  const isLoginPage = routerState.location.pathname === '/login'
  const { isAuthenticated } = useAuth()
  const showSidebar = isAuthenticated && !isLoginPage

  if (!showSidebar) {
    return <Outlet />
  }

  return (
    <div className="flex h-screen overflow-hidden">
      <Sidebar />
      <main className="flex-1 min-w-0 overflow-auto">
        <Outlet />
      </main>
    </div>
  )
}

export const Route = createRootRouteWithContext<MyRouterContext>()({
  component: RootComponent,
})
