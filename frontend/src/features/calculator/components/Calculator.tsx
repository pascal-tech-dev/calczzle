import { useCalculator } from '../hooks/useCalculator'
import { useCalculatorKeyboard } from '../hooks/useCalculatorKeyboard'
import styles from '../styles/Calculator.module.css'
import { CalculatorDisplay } from './CalculatorDisplay'
import { CalculatorKeypad } from './CalculatorKeypad'

export function Calculator() {
  const {
    expression,
    displayValue,
    error,
    isLoading,
    append,
    clear,
    deleteLast,
    evaluate,
  } = useCalculator()

  useCalculatorKeyboard({
    isLoading,
    append,
    clear,
    deleteLast,
    evaluate: () => {
      void evaluate()
    },
  })

  return (
    <div
      className={styles.calculator}
      role="region"
      aria-label="Calculator"
    >
      <CalculatorDisplay
        displayValue={displayValue}
        expression={expression}
        isLoading={isLoading}
        error={error}
      />
      <CalculatorKeypad
        isLoading={isLoading}
        onAppend={append}
        onClear={clear}
        onDelete={deleteLast}
        onEvaluate={() => {
          void evaluate()
        }}
      />
    </div>
  )
}
