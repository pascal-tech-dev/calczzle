import type { CalculatorAction, CalculatorState } from './calculator.types'
import { appendToExpression, deleteLastCharacter } from './expressionBuilder'

export const initialCalculatorState: CalculatorState = {
  expression: '',
  result: null,
  status: 'idle',
  error: null,
  justEvaluated: false,
}

export function calculatorReducer(
  state: CalculatorState,
  action: CalculatorAction,
): CalculatorState {
  switch (action.type) {
    case 'APPEND':
      return {
        ...state,
        expression: appendToExpression(state.expression, action.value, {
          justEvaluated: state.justEvaluated,
          result: state.result,
        }),
        result: null,
        status: 'idle',
        error: null,
        justEvaluated: false,
      }

    case 'DELETE_LAST':
      return {
        ...state,
        expression: deleteLastCharacter(state.expression),
        result: null,
        status: 'idle',
        error: null,
        justEvaluated: false,
      }

    case 'CLEAR':
      return { ...initialCalculatorState }

    case 'EVALUATE_START':
      return {
        ...state,
        status: 'loading',
        error: null,
      }

    case 'EVALUATE_SUCCESS':
      return {
        ...state,
        result: action.result,
        status: 'success',
        error: null,
        justEvaluated: true,
      }

    case 'EVALUATE_FAILURE':
      return {
        ...state,
        status: 'error',
        error: action.message,
        justEvaluated: false,
      }

    default: {
      const _exhaustive: never = action
      return _exhaustive
    }
  }
}
