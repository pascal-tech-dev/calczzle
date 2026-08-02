/** Align with backend limits when known; keeps the UI from sending oversized input. */
export const MAX_EXPRESSION_LENGTH = 256

export interface AppendExpressionOptions {
  justEvaluated: boolean
  result: number | null
}

function formatResult(result: number): string {
  return String(result)
}

function startsNewExpression(value: string): boolean {
  return /^\d$/.test(value) || value === '.' || value === '(' || value === 'sqrt('
}

function continuesFromResult(value: string): boolean {
  return /^[+\-*/^%)]$/.test(value)
}

function currentNumberFragment(expression: string): string {
  const match = expression.match(/(\d*\.?\d*)$/)
  return match?.[1] ?? ''
}

export function canAppendDecimal(expression: string): boolean {
  return !currentNumberFragment(expression).includes('.')
}

export function deleteLastCharacter(expression: string): string {
  if (expression.length === 0) {
    return expression
  }

  return expression.slice(0, -1)
}

export function appendToExpression(
  expression: string,
  value: string,
  options: AppendExpressionOptions,
): string {
  let base = expression

  if (options.justEvaluated && options.result !== null) {
    if (startsNewExpression(value)) {
      base = ''
    } else if (continuesFromResult(value)) {
      base = formatResult(options.result)
    }
  }

  if (value === '.' && !canAppendDecimal(base)) {
    return base
  }

  const next = `${base}${value}`

  if (next.length > MAX_EXPRESSION_LENGTH) {
    return base
  }

  return next
}
