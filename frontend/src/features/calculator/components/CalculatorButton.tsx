import type { CalculatorKeyKind } from '../model/calculator.types'
import styles from '../styles/Calculator.module.css'

const kindClassNames: Record<CalculatorKeyKind, string | undefined> = {
  digit: styles.digit,
  decimal: styles.digit,
  operator: styles.operator,
  parenthesis: styles.functionKey,
  function: styles.functionKey,
  percentage: styles.operator,
  action: styles.action,
}

interface CalculatorButtonProps {
  label: string
  ariaLabel: string
  kind?: CalculatorKeyKind
  disabled?: boolean
  onClick?: () => void
  className?: string
}

export function CalculatorButton({
  label,
  ariaLabel,
  kind,
  disabled = false,
  onClick,
  className,
}: CalculatorButtonProps) {
  const kindClassName = kind ? kindClassNames[kind] : undefined
  const combinedClassName = [styles.button, kindClassName, className]
    .filter(Boolean)
    .join(' ')

  return (
    <button
      type="button"
      aria-label={ariaLabel}
      disabled={disabled}
      onClick={onClick}
      className={combinedClassName}
    >
      {label}
    </button>
  )
}
