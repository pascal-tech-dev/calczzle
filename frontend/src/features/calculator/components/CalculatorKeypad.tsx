import { CalculatorButton } from './CalculatorButton'
import styles from '../styles/Calculator.module.css'

export function CalculatorKeypad() {
  return (
    <div className={styles.keypad}>
      <CalculatorButton label="…" ariaLabel="Keys coming soon" disabled />
    </div>
  )
}
