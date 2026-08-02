import styles from '../styles/Calculator.module.css'

export function CalculatorDisplay() {
  return (
    <div className={styles.display} aria-live="polite" aria-atomic="true">
      0
    </div>
  )
}
