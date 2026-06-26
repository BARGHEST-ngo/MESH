import * as React from "react"

interface DialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  children: React.ReactNode
}

export function Dialog({ open, onOpenChange, children }: DialogProps) {
  if (!open) return null

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-6">
      <div
        className="fixed inset-0 bg-black/55 backdrop-blur-[2px] motion-safe:[animation:cpFade_.18s_ease]"
        onClick={() => onOpenChange(false)}
      />
      <div className="relative z-50 w-full max-w-lg motion-safe:[animation:cpPop_.22s_cubic-bezier(0.16,1,0.3,1)]">
        {children}
      </div>
    </div>
  )
}

interface DialogContentProps {
  children: React.ReactNode
  className?: string
}

export function DialogContent({ children, className = "" }: DialogContentProps) {
  return (
    <div className={`bg-bg-raised border border-border rounded-[18px] shadow-[0_24px_70px_rgba(0,0,0,0.5)] p-6 w-full max-h-[90vh] overflow-y-auto ${className}`}>
      {children}
    </div>
  )
}

interface DialogHeaderProps {
  children: React.ReactNode
}

export function DialogHeader({ children }: DialogHeaderProps) {
  return <div className="mb-4">{children}</div>
}

interface DialogTitleProps {
  children: React.ReactNode
}

export function DialogTitle({ children }: DialogTitleProps) {
  return <h2 className="text-xl font-bold text-foreground tracking-tight">{children}</h2>
}

interface DialogDescriptionProps {
  children: React.ReactNode
}

export function DialogDescription({ children }: DialogDescriptionProps) {
  return <p className="text-muted-foreground text-sm mt-2">{children}</p>
}

interface DialogFooterProps {
  children: React.ReactNode
}

export function DialogFooter({ children }: DialogFooterProps) {
  return <div className="mt-6 flex justify-end gap-3">{children}</div>
}

