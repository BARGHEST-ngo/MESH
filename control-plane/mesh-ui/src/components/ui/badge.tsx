import * as React from "react"

interface BadgeProps {
  children: React.ReactNode
  variant?: "default" | "success" | "warning" | "destructive" | "secondary"
  className?: string
}

export function Badge({ children, variant = "default", className = "" }: BadgeProps) {
  const variantClasses = {
    default: "bg-primary text-primary-foreground",
    success: "bg-green-500/20 text-green-400 border border-green-500/50",
    warning: "bg-yellow-500/20 text-yellow-400 border border-yellow-500/50",
    destructive: "bg-red-500/20 text-red-400 border border-red-500/50",
    secondary: "bg-secondary text-secondary-foreground",
  }

  return (
    <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-semibold font-mono ${variantClasses[variant]} ${className}`}>
      {children}
    </span>
  )
}

