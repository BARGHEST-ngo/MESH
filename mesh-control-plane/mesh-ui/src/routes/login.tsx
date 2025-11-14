import { createFileRoute, redirect } from '@tanstack/react-router'
import LoginForm from '../components/LoginForm'

export const Route = createFileRoute('/login')({
  beforeLoad: ({ context }) => {
    if (context.auth.isAuthenticated && !context.auth.isLoading) {
      throw redirect({ to: '/' })
    }
  },
  component: LoginRoute,
})

function LoginRoute() {
  return <LoginForm />
}

