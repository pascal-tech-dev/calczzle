import { CalculatorDisplay } from './CalculatorDisplay'
import { CalculatorKeypad } from './CalculatorKeypad'
import styles from '../styles/Calculator.module.css'

export function Calculator() {
  return (
    <div className={styles.calculator}>
      <CalculatorDisplay />
      <CalculatorKeypad />
    </div>
  )
}
