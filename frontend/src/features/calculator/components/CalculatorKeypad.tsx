import { calculatorKeys } from '../model/calculatorKeys'
import type { CalculatorKey } from '../model/calculator.types'
import styles from '../styles/Calculator.module.css'
import { CalculatorButton } from './CalculatorButton'

interface CalculatorKeypadProps {
  isLoading: boolean
  onAppend: (value: string) => void
  onClear: () => void
  onDelete: () => void
  onEvaluate: () => void
}

function handleKeyPress(
  key: CalculatorKey,
  actions: Omit<CalculatorKeypadProps, 'isLoading'>,
) {
  if (key.action === 'clear') {
    actions.onClear()
    return
  }

  if (key.action === 'delete') {
    actions.onDelete()
    return
  }

  if (key.action === 'evaluate') {
    actions.onEvaluate()
    return
  }

  if (key.value !== undefined) {
    actions.onAppend(key.value)
  }
}

export function CalculatorKeypad({
  isLoading,
  onAppend,
  onClear,
  onDelete,
  onEvaluate,
}: CalculatorKeypadProps) {
  return (
    <div className={styles.keypad}>
      {calculatorKeys.map((key) => {
        const extraClassName =
          key.className === 'spanTwo' ? styles.spanTwo : undefined

        const isEvaluateKey = key.action === 'evaluate'

        return (
          <CalculatorButton
            key={key.id}
            label={key.label}
            ariaLabel={key.ariaLabel}
            kind={key.kind}
            className={extraClassName}
            disabled={isLoading && isEvaluateKey}
            onClick={() =>
              handleKeyPress(key, {
                onAppend,
                onClear,
                onDelete,
                onEvaluate,
              })
            }
          />
        )
      })}
    </div>
  )
}
