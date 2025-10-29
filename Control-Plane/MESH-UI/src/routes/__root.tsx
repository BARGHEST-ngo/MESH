import { Outlet, createRootRouteWithContext, useRouterState } from '@tanstack/react-router'

import Header from '../components/Header'

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

  return (
    <>
      {!isLoginPage && <Header />}
      <Outlet />
    </>
  )
}

export const Route = createRootRouteWithContext<MyRouterContext>()({
  component: RootComponent,
})
