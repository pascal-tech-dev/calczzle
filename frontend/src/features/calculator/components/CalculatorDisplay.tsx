import styles from '../styles/Calculator.module.css'

interface CalculatorDisplayProps {
  displayValue: string
  expression: string
  isLoading: boolean
  error: string | null
}

export function CalculatorDisplay({
  displayValue,
  expression,
  isLoading,
  error,
}: CalculatorDisplayProps) {
  const showExpressionDetail =
    expression !== '' && displayValue !== expression

  return (
    <div className={styles.displayPanel}>
      {showExpressionDetail ? (
        <div className={styles.expression}>{expression}</div>
      ) : null}
      <div
        className={styles.display}
        aria-live="polite"
        aria-atomic="true"
        aria-label="Display"
      >
        {displayValue}
      </div>
      {isLoading ? (
        <p className={styles.status} role="status">
          Calculating…
        </p>
      ) : null}
      {error ? (
        <p className={styles.error} role="alert">
          {error}
        </p>
      ) : null}
    </div>
  )
}
