import { describe, expect, it } from 'vitest'
import {
  calculatorReducer,
  initialCalculatorState,
} from './calculator.reducer'
import type { CalculatorState } from './calculator.types'

describe('calculatorReducer', () => {
  it('APPEND adds a value', () => {
    const next = calculatorReducer(initialCalculatorState, {
      type: 'APPEND',
      value: '3',
    })

    expect(next.expression).toBe('3')
    expect(next.error).toBeNull()
    expect(next.status).toBe('idle')
  })

  it('CLEAR resets all state', () => {
    const dirty: CalculatorState = {
      expression: '3+4',
      result: 7,
      status: 'success',
      error: 'old',
      justEvaluated: true,
    }

    expect(calculatorReducer(dirty, { type: 'CLEAR' })).toEqual(
      initialCalculatorState,
    )
  })

  it('DELETE_LAST removes one character', () => {
    const state: CalculatorState = {
      ...initialCalculatorState,
      expression: '12',
    }

    const next = calculatorReducer(state, { type: 'DELETE_LAST' })
    expect(next.expression).toBe('1')
  })

  it('EVALUATE_START sets loading state', () => {
    const state: CalculatorState = {
      ...initialCalculatorState,
      expression: '1+1',
      error: 'old error',
    }

    const next = calculatorReducer(state, { type: 'EVALUATE_START' })
    expect(next.status).toBe('loading')
    expect(next.error).toBeNull()
    expect(next.expression).toBe('1+1')
  })

  it('EVALUATE_SUCCESS stores the result', () => {
    const state: CalculatorState = {
      ...initialCalculatorState,
      expression: '3+4',
      status: 'loading',
    }

    const next = calculatorReducer(state, {
      type: 'EVALUATE_SUCCESS',
      result: 7,
    })

    expect(next.result).toBe(7)
    expect(next.status).toBe('success')
    expect(next.justEvaluated).toBe(true)
    expect(next.error).toBeNull()
  })

  it('EVALUATE_FAILURE stores the error', () => {
    const state: CalculatorState = {
      ...initialCalculatorState,
      expression: '10/0',
      status: 'loading',
    }

    const next = calculatorReducer(state, {
      type: 'EVALUATE_FAILURE',
      message: 'Division by zero',
    })

    expect(next.status).toBe('error')
    expect(next.error).toBe('Division by zero')
    expect(next.expression).toBe('10/0')
    expect(next.justEvaluated).toBe(false)
  })

  it('digit after evaluation starts a new expression', () => {
    const state: CalculatorState = {
      expression: '3+4',
      result: 7,
      status: 'success',
      error: null,
      justEvaluated: true,
    }

    const next = calculatorReducer(state, { type: 'APPEND', value: '9' })
    expect(next.expression).toBe('9')
    expect(next.justEvaluated).toBe(false)
    expect(next.result).toBeNull()
  })

  it('operator after evaluation continues from the result', () => {
    const state: CalculatorState = {
      expression: '3+4',
      result: 7,
      status: 'success',
      error: null,
      justEvaluated: true,
    }

    const next = calculatorReducer(state, { type: 'APPEND', value: '+' })
    expect(next.expression).toBe('7+')
    expect(next.justEvaluated).toBe(false)
  })
})
