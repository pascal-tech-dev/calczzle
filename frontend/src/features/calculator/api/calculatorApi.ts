import type {
  ApiErrorResponse,
  EvaluateExpressionRequest,
  EvaluateExpressionResponse,
} from '../model/calculator.types'

export class CalculatorApiError extends Error {
  readonly code: string
  readonly status: number

  constructor(message: string, code: string, status: number) {
    super(message)
    this.name = 'CalculatorApiError'
    this.code = code
    this.status = status
  }
}

const EVALUATE_ENDPOINT = '/api/v1/evaluate'

export async function evaluateExpression(
  expression: string,
  signal?: AbortSignal,
): Promise<number> {
  const request: EvaluateExpressionRequest = { expression }

  let response: Response

  try {
    response = await fetch(EVALUATE_ENDPOINT, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(request),
      signal,
    })
  } catch (error) {
    if (error instanceof DOMException && error.name === 'AbortError') {
      throw error
    }

    throw new CalculatorApiError(
      'The calculator service is unavailable. Please try again.',
      'NETWORK_ERROR',
      0,
    )
  }

  if (!response.ok) {
    let errorResponse: ApiErrorResponse | null = null

    try {
      errorResponse = (await response.json()) as ApiErrorResponse
    } catch {
      // The backend returned a non-JSON error response.
    }

    throw new CalculatorApiError(
      errorResponse?.error.message ?? 'The expression could not be evaluated.',
      errorResponse?.error.code ?? 'UNKNOWN_API_ERROR',
      response.status,
    )
  }

  const data = (await response.json()) as EvaluateExpressionResponse

  if (typeof data.result !== 'number' || !Number.isFinite(data.result)) {
    throw new CalculatorApiError(
      'The server returned an invalid result.',
      'INVALID_API_RESPONSE',
      response.status,
    )
  }

  return data.result
}
