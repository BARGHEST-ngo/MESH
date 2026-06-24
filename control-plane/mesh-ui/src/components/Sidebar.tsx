import { Link, useNavigate } from '@tanstack/react-router'
import { Home, Shield, LogOut, Sun, Moon } from 'lucide-react'
import { useAuth } from '../lib/auth'
import { useTheme } from '../lib/theme'
import { Button } from './ui/button'
import { Wordmark } from './Wordmark'

const navBase = 'flex items-center gap-3 px-3.5 py-3 rounded-[10px] text-[14.5px] font-semibold transition-colors border'
const navItem = `${navBase} border-transparent text-foreground hover:bg-card-hi`
const navItemActive = `${navBase} border-primary/40 bg-primary/15 text-primary`

export default function Sidebar() {
  const { logout } = useAuth()
  const { theme, toggleTheme } = useTheme()
  const navigate = useNavigate()

  const handleLogout = () => {
    logout()
    navigate({ to: '/login' })
  }

  return (
    <aside className="w-[252px] shrink-0 bg-sidebar border-r border-border flex flex-col">
      <div className="flex items-end gap-2.5 px-[18px] h-[60px] border-b border-border">
        <Wordmark height={17} className="mb-5 text-foreground" />
        <span className="text-[13px] text-muted-foreground mb-5 leading-none">Control Plane</span>
      </div>

      <nav className="flex-1 p-3 flex flex-col gap-1">
        <div className="text-[11px] font-semibold tracking-[0.3px] text-muted-foreground uppercase px-3.5 pt-2 pb-1.5">
          Menu
        </div>
        <Link to="/" className={navItem} activeOptions={{ exact: true }} activeProps={{ className: navItemActive }}>
          <Home size={19} />
          Networks
        </Link>
        <Link to="/acl" className={navItem} activeProps={{ className: navItemActive }}>
          <Shield size={19} />
          Access policy
        </Link>
      </nav>

      <div className="p-3 border-t border-border flex flex-col gap-2">
        <button
          onClick={toggleTheme}
          className="flex items-center gap-3 px-3.5 py-[11px] rounded-[10px] border border-border text-foreground text-[13.5px] font-semibold w-full text-left hover:bg-card-hi transition-colors"
        >
          {theme === 'dark' ? <Sun size={18} className="text-text2" /> : <Moon size={18} className="text-text2" />}
          {theme === 'dark' ? 'Light mode' : 'Dark mode'}
        </button>
        <Button variant="destructive" onClick={handleLogout} className="w-full">
          <LogOut size={18} />
          Sign out
        </Button>
      </div>
    </aside>
  )
}
