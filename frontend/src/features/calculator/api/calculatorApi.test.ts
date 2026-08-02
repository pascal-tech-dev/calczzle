import { afterEach, describe, expect, it, vi } from 'vitest'
import { CalculatorApiError, evaluateExpression } from './calculatorApi'

describe('calculatorApi', () => {
  afterEach(() => {
    vi.unstubAllGlobals()
    vi.restoreAllMocks()
  })

  it('returns result on 200 response', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn().mockResolvedValue({
        ok: true,
        status: 200,
        json: async () => ({ result: 0.125 }),
      }),
    )

    await expect(evaluateExpression('2^-3')).resolves.toBe(0.125)
    expect(fetch).toHaveBeenCalledWith(
      '/api/v1/evaluate',
      expect.objectContaining({
        method: 'POST',
        body: JSON.stringify({ expression: '2^-3' }),
      }),
    )
  })

  it('throws CalculatorApiError with backend message on 400', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn().mockResolvedValue({
        ok: false,
        status: 400,
        json: async () => ({
          error: {
            code: 'INVALID_EXPRESSION',
            message: 'The expression is invalid.',
          },
        }),
      }),
    )

    await expect(evaluateExpression('3+*4')).rejects.toMatchObject({
      name: 'CalculatorApiError',
      message: 'The expression is invalid.',
      code: 'INVALID_EXPRESSION',
      status: 400,
    } satisfies Partial<CalculatorApiError>)
  })

  it('creates a generic API error for 500 non-JSON responses', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn().mockResolvedValue({
        ok: false,
        status: 500,
        json: async () => {
          throw new Error('not json')
        },
      }),
    )

    await expect(evaluateExpression('1+1')).rejects.toMatchObject({
      name: 'CalculatorApiError',
      message: 'The expression could not be evaluated.',
      code: 'UNKNOWN_API_ERROR',
      status: 500,
    })
  })

  it('creates NETWORK_ERROR when fetch is rejected', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn().mockRejectedValue(new TypeError('Failed to fetch')),
    )

    await expect(evaluateExpression('1+1')).rejects.toMatchObject({
      name: 'CalculatorApiError',
      code: 'NETWORK_ERROR',
      status: 0,
      message: 'The calculator service is unavailable. Please try again.',
    })
  })

  it('creates INVALID_API_RESPONSE for non-finite results', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn().mockResolvedValue({
        ok: true,
        status: 200,
        json: async () => ({ result: Number.NaN }),
      }),
    )

    await expect(evaluateExpression('1+1')).rejects.toMatchObject({
      name: 'CalculatorApiError',
      code: 'INVALID_API_RESPONSE',
      message: 'The server returned an invalid result.',
    })
  })

  it('creates INVALID_API_RESPONSE for malformed results', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn().mockResolvedValue({
        ok: true,
        status: 200,
        json: async () => ({ result: 'not-a-number' }),
      }),
    )

    await expect(evaluateExpression('1+1')).rejects.toMatchObject({
      code: 'INVALID_API_RESPONSE',
    })
  })
})
