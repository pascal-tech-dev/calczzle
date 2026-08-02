interface CalculatorButtonProps {
  label: string
  ariaLabel: string
  disabled?: boolean
  onClick?: () => void
  className?: string
}

export function CalculatorButton({
  label,
  ariaLabel,
  disabled = false,
  onClick,
  className,
}: CalculatorButtonProps) {
  return (
    <button
      type="button"
      aria-label={ariaLabel}
      disabled={disabled}
      onClick={onClick}
      className={className}
    >
      {label}
    </button>
  )
}
