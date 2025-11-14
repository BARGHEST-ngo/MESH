import * as React from "react"

interface CardProps {
  children: React.ReactNode
  className?: string
}

export function Card({ children, className = "" }: CardProps) {
  return (
    <div className={`bg-card border border-border rounded-lg shadow-lg ${className}`}>
      {children}
    </div>
  )
}

interface CardHeaderProps {
  children: React.ReactNode
  className?: string
}

export function CardHeader({ children, className = "" }: CardHeaderProps) {
  return <div className={`p-6 pb-4 ${className}`}>{children}</div>
}

interface CardTitleProps {
  children: React.ReactNode
  className?: string
}

export function CardTitle({ children, className = "" }: CardTitleProps) {
  return <h3 className={`text-xl font-bold text-foreground font-mono ${className}`}>{children}</h3>
}

interface CardDescriptionProps {
  children: React.ReactNode
}

export function CardDescription({ children }: CardDescriptionProps) {
  return <p className="text-muted-foreground text-sm mt-1">{children}</p>
}

interface CardContentProps {
  children: React.ReactNode
  className?: string
}

export function CardContent({ children, className = "" }: CardContentProps) {
  return <div className={`p-6 pt-0 ${className}`}>{children}</div>
}

interface CardFooterProps {
  children: React.ReactNode
  className?: string
}

export function CardFooter({ children, className = "" }: CardFooterProps) {
  return <div className={`p-6 pt-0 flex items-center ${className}`}>{children}</div>
}

