import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { Calculator } from './Calculator'

function mockFetchResult(result: number) {
  vi.stubGlobal(
    'fetch',
    vi.fn().mockResolvedValue({
      ok: true,
      status: 200,
      json: async () => ({ result }),
    }),
  )
}

function mockFetchError(message: string, status = 400) {
  vi.stubGlobal(
    'fetch',
    vi.fn().mockResolvedValue({
      ok: false,
      status,
      json: async () => ({
        error: {
          code: 'INVALID_EXPRESSION',
          message,
        },
      }),
    }),
  )
}

function getDisplayValue() {
  return screen.getByLabelText('Display')
}

describe('Calculator', () => {
  afterEach(() => {
    vi.unstubAllGlobals()
    vi.restoreAllMocks()
  })

  it('initially displays 0', () => {
    render(<Calculator />)
    expect(getDisplayValue()).toHaveTextContent('0')
  })

  it('builds 3+4 when clicking 3, +, 4', async () => {
    const user = userEvent.setup()
    render(<Calculator />)

    await user.click(screen.getByRole('button', { name: '3' }))
    await user.click(screen.getByRole('button', { name: 'Add' }))
    await user.click(screen.getByRole('button', { name: '4' }))

    expect(getDisplayValue()).toHaveTextContent('3+4')
  })

  it('clears the expression when clicking clear', async () => {
    const user = userEvent.setup()
    render(<Calculator />)

    await user.click(screen.getByRole('button', { name: '3' }))
    await user.click(screen.getByRole('button', { name: 'Clear expression' }))

    expect(getDisplayValue()).toHaveTextContent('0')
  })

  it('removes the final character when clicking delete', async () => {
    const user = userEvent.setup()
    render(<Calculator />)

    await user.click(screen.getByRole('button', { name: '3' }))
    await user.click(screen.getByRole('button', { name: '4' }))
    await user.click(
      screen.getByRole('button', { name: 'Delete last character' }),
    )

    expect(getDisplayValue()).toHaveTextContent('3')
  })

  it('appends sqrt( when clicking square root', async () => {
    const user = userEvent.setup()
    render(<Calculator />)

    await user.click(screen.getByRole('button', { name: /square root/i }))

    expect(getDisplayValue()).toHaveTextContent('sqrt(')
  })

  it('stores * when clicking multiplication even though the label is ×', async () => {
    const user = userEvent.setup()
    render(<Calculator />)

    await user.click(screen.getByRole('button', { name: '3' }))
    await user.click(screen.getByRole('button', { name: /multiply/i }))
    await user.click(screen.getByRole('button', { name: '4' }))

    expect(getDisplayValue()).toHaveTextContent('3*4')
  })

  it('sends the current expression when clicking equals', async () => {
    const user = userEvent.setup()
    mockFetchResult(7)
    render(<Calculator />)

    await user.click(screen.getByRole('button', { name: '3' }))
    await user.click(screen.getByRole('button', { name: 'Add' }))
    await user.click(screen.getByRole('button', { name: '4' }))
    await user.click(
      screen.getByRole('button', { name: /evaluate expression/i }),
    )

    await waitFor(() => {
      expect(fetch).toHaveBeenCalledWith(
        '/api/v1/evaluate',
        expect.objectContaining({
          method: 'POST',
          body: JSON.stringify({ expression: '3+4' }),
        }),
      )
    })
  })

  it('displays the result after a successful response', async () => {
    const user = userEvent.setup()
    mockFetchResult(0.125)
    render(<Calculator />)

    await user.click(screen.getByRole('button', { name: '2' }))
    await user.click(screen.getByRole('button', { name: /exponentiation/i }))
    await user.click(screen.getByRole('button', { name: 'Subtract' }))
    await user.click(screen.getByRole('button', { name: '3' }))
    await user.click(
      screen.getByRole('button', { name: /evaluate expression/i }),
    )

    await waitFor(() => {
      expect(getDisplayValue()).toHaveTextContent('0.125')
    })
  })

  it('shows a backend error with role="alert"', async () => {
    const user = userEvent.setup()
    mockFetchError('The expression is invalid.')
    render(<Calculator />)

    await user.click(screen.getByRole('button', { name: '3' }))
    await user.click(
      screen.getByRole('button', { name: /evaluate expression/i }),
    )

    expect(await screen.findByRole('alert')).toHaveTextContent(
      'The expression is invalid.',
    )
  })

  it('disables equals while loading', async () => {
    const user = userEvent.setup()
    let resolveFetch: ((value: unknown) => void) | undefined

    vi.stubGlobal(
      'fetch',
      vi.fn().mockImplementation(
        () =>
          new Promise((resolve) => {
            resolveFetch = resolve
          }),
      ),
    )

    render(<Calculator />)

    await user.click(screen.getByRole('button', { name: '1' }))
    const equals = screen.getByRole('button', {
      name: /evaluate expression/i,
    })
    await user.click(equals)

    await waitFor(() => {
      expect(equals).toBeDisabled()
    })
    expect(screen.getByRole('status')).toHaveTextContent('Calculating…')

    resolveFetch?.({
      ok: true,
      status: 200,
      json: async () => ({ result: 1 }),
    })

    await waitFor(() => {
      expect(equals).toBeEnabled()
    })
  })

  it('evaluates 2^-3 when typed on the keyboard and Enter is pressed', async () => {
    const user = userEvent.setup()
    mockFetchResult(0.125)
    render(<Calculator />)

    await user.keyboard('2^-3{Enter}')

    await waitFor(() => {
      expect(fetch).toHaveBeenCalledWith(
        '/api/v1/evaluate',
        expect.objectContaining({
          body: JSON.stringify({ expression: '2^-3' }),
        }),
      )
    })
    await waitFor(() => {
      expect(getDisplayValue()).toHaveTextContent('0.125')
    })
  })

  it('deletes the last character on Backspace', async () => {
    const user = userEvent.setup()
    render(<Calculator />)

    await user.keyboard('12{Backspace}')

    expect(getDisplayValue()).toHaveTextContent('1')
  })

  it('clears the calculator on Escape', async () => {
    const user = userEvent.setup()
    render(<Calculator />)

    await user.keyboard('99{Escape}')

    expect(getDisplayValue()).toHaveTextContent('0')
  })

  it('ignores unsupported keys', async () => {
    const user = userEvent.setup()
    render(<Calculator />)

    await user.keyboard('3abc4')

    expect(getDisplayValue()).toHaveTextContent('34')
  })
})
