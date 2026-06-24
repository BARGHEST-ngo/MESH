import * as React from "react"

interface BadgeProps {
  children: React.ReactNode
  variant?: "default" | "success" | "warning" | "destructive" | "accent" | "secondary"
  className?: string
}

export function Badge({ children, variant = "secondary", className = "" }: BadgeProps) {
  const variantClasses = {
    default: "bg-primary text-primary-foreground border-transparent",
    success: "bg-success-dim text-success border-success/25",
    warning: "bg-warning-dim text-warning border-warning/25",
    destructive: "bg-destructive-dim text-destructive border-destructive/25",
    accent: "bg-primary/15 text-primary border-primary/25",
    secondary: "bg-inset text-text2 border-border",
  }

  return (
    <span
      className={`inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full border text-[11.5px] font-semibold tracking-[0.2px] ${variantClasses[variant]} ${className}`}
    >
      {children}
    </span>
  )
}
