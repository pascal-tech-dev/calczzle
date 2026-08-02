export interface EvaluateExpressionRequest {
  expression: string
}

export interface EvaluateExpressionResponse {
  result: number
}

export interface ApiErrorBody {
  code: string
  message: string
}

export interface ApiErrorResponse {
  error: ApiErrorBody
}

export type CalculatorStatus = 'idle' | 'loading' | 'success' | 'error'

export interface CalculatorState {
  expression: string
  result: number | null
  status: CalculatorStatus
  error: string | null
  justEvaluated: boolean
}

export type CalculatorAction =
  | { type: 'APPEND'; value: string }
  | { type: 'DELETE_LAST' }
  | { type: 'CLEAR' }
  | { type: 'EVALUATE_START' }
  | { type: 'EVALUATE_SUCCESS'; result: number }
  | { type: 'EVALUATE_FAILURE'; message: string }

export type CalculatorKeyKind =
  | 'digit'
  | 'decimal'
  | 'operator'
  | 'parenthesis'
  | 'function'
  | 'percentage'
  | 'action'

export type CalculatorActionKey = 'clear' | 'delete' | 'evaluate'

export interface CalculatorKey {
  id: string
  label: string
  value?: string
  kind: CalculatorKeyKind
  action?: CalculatorActionKey
  ariaLabel: string
  className?: string
}
