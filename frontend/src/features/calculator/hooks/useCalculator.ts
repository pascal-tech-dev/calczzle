import { useEffect, useReducer, useRef } from 'react'
import {
  CalculatorApiError,
  evaluateExpression,
} from '../api/calculatorApi'
import {
  calculatorReducer,
  initialCalculatorState,
} from '../model/calculator.reducer'
import type { CalculatorStatus } from '../model/calculator.types'

export interface UseCalculatorResult {
  expression: string
  displayValue: string
  status: CalculatorStatus
  error: string | null
  isLoading: boolean
  append: (value: string) => void
  clear: () => void
  deleteLast: () => void
  evaluate: () => Promise<void>
}

export function useCalculator(): UseCalculatorResult {
  const [state, dispatch] = useReducer(
    calculatorReducer,
    initialCalculatorState,
  )
  const abortControllerRef = useRef<AbortController | null>(null)

  useEffect(() => {
    return () => {
      abortControllerRef.current?.abort()
    }
  }, [])

  const isLoading = state.status === 'loading'

  const displayValue =
    isLoading || !(state.justEvaluated && state.result !== null)
      ? state.expression === ''
        ? '0'
        : state.expression
      : String(state.result)

  function abortInFlightRequest() {
    abortControllerRef.current?.abort()
    abortControllerRef.current = null
  }

  function append(value: string) {
    if (isLoading) {
      abortInFlightRequest()
    }
    dispatch({ type: 'APPEND', value })
  }

  function clear() {
    abortInFlightRequest()
    dispatch({ type: 'CLEAR' })
  }

  function deleteLast() {
    if (isLoading) {
      abortInFlightRequest()
    }
    dispatch({ type: 'DELETE_LAST' })
  }

  async function evaluate() {
    if (state.status === 'loading') {
      return
    }

    const expressionToEvaluate = state.expression

    if (expressionToEvaluate.trim() === '') {
      return
    }

    abortInFlightRequest()
    const controller = new AbortController()
    abortControllerRef.current = controller

    dispatch({ type: 'EVALUATE_START' })

    try {
      const result = await evaluateExpression(
        expressionToEvaluate,
        controller.signal,
      )
      dispatch({ type: 'EVALUATE_SUCCESS', result })
    } catch (error) {
      if (error instanceof DOMException && error.name === 'AbortError') {
        return
      }

      const message =
        error instanceof CalculatorApiError
          ? error.message
          : 'An unexpected error occurred.'

      dispatch({ type: 'EVALUATE_FAILURE', message })
    }
  }

  return {
    expression: state.expression,
    displayValue,
    status: state.status,
    error: state.error,
    isLoading,
    append,
    clear,
    deleteLast,
    evaluate,
  }
}
