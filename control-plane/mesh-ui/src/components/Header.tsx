import { Link, useNavigate } from '@tanstack/react-router'

import { useState } from 'react'
import { Home, Menu, X, LogOut, Sun, Moon, Cog } from 'lucide-react'
import { useAuth } from '../lib/auth'
import { useTheme } from '../lib/theme'
import { Button } from './ui/button'

export default function Header() {
  const [isOpen, setIsOpen] = useState(false)
  const { isAuthenticated, logout } = useAuth()
  const { theme, toggleTheme } = useTheme()
  const navigate = useNavigate()

  const handleLogout = () => {
    logout()
    setIsOpen(false)
    navigate({ to: '/login' } as any)
  }

  return (
    <>
      <header className="p-4 flex items-center justify-between bg-card border-b border-border text-foreground shadow-lg">
        <div className="flex items-center">
          <button
            onClick={() => setIsOpen(true)}
            className="p-2 hover:bg-secondary rounded transition-colors"
            aria-label="Open menu"
          >
            <Menu size={24} />
          </button>
          <h1 className="ml-4 text-xl font-bold font-mono">
            <Link to="/">&gt; MESH</Link>
          </h1>
        </div>
        <button
          onClick={toggleTheme}
          className="p-2 hover:bg-secondary rounded transition-colors"
          aria-label="Toggle theme"
        >
          {theme === 'dark' ? <Sun size={20} /> : <Moon size={20} />}
        </button>
      </header>

      <aside
        className={`fixed top-0 left-0 h-full w-80 bg-sidebar border-r border-sidebar-border text-sidebar-foreground shadow-2xl z-50 transform transition-transform duration-300 ease-in-out flex flex-col ${
          isOpen ? 'translate-x-0' : '-translate-x-full'
        }`}
      >
        <div className="flex items-center justify-between p-4 border-b border-sidebar-border">
          <h2 className="text-xl font-bold font-mono">[ NAVIGATION ]</h2>
          <button
            onClick={() => setIsOpen(false)}
            className="p-2 hover:bg-sidebar-accent rounded transition-colors"
            aria-label="Close menu"
          >
            <X size={24} />
          </button>
        </div>

        <nav className="flex-1 p-4 overflow-y-auto">
          <Link
            to="/"
            onClick={() => setIsOpen(false)}
            className="flex items-center gap-3 p-3 rounded hover:bg-sidebar-accent transition-colors mb-2 font-mono"
            activeProps={{
              className:
                'flex items-center gap-3 p-3 rounded bg-sidebar-accent border border-sidebar-ring transition-colors mb-2 font-mono',
            }}
          >
            <Home size={20} />
            <span className="font-medium">&gt; HOME</span>
          </Link>

          <Link
            to="/acl"
            onClick={() => setIsOpen(false)}
            className="flex items-center gap-3 p-3 rounded hover:bg-sidebar-accent transition-colors mb-2 font-mono"
            activeProps={{
              className:
                'flex items-center gap-3 p-3 rounded bg-sidebar-accent border border-sidebar-ring transition-colors mb-2 font-mono',
            }}
          >
            <Cog size={20} />
            <span className="font-medium">&gt; CUSTOMIZE ACL</span>
          </Link>
        </nav>


        {isAuthenticated && (
          <div className="p-4 border-t border-sidebar-border">
            <Button
              onClick={handleLogout}
              variant="destructive"
              className="w-full flex items-center justify-center gap-2 font-mono"
            >
              <LogOut size={18} />
              <span>[ LOGOUT ]</span>
            </Button>
          </div>
        )}
      </aside>
    </>
  )
}
