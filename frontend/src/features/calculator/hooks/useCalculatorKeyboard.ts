import { useEffect, useRef } from 'react'

interface CalculatorKeyboardActions {
  isLoading: boolean
  append: (value: string) => void
  clear: () => void
  deleteLast: () => void
  evaluate: () => void
}

function isTypingInField(target: EventTarget | null): boolean {
  if (!(target instanceof HTMLElement)) {
    return false
  }

  const tagName = target.tagName
  return (
    tagName === 'INPUT' ||
    tagName === 'TEXTAREA' ||
    tagName === 'SELECT' ||
    target.isContentEditable
  )
}

export function useCalculatorKeyboard(actions: CalculatorKeyboardActions) {
  const actionsRef = useRef(actions)

  useEffect(() => {
    actionsRef.current = actions
  })

  useEffect(() => {
    function handleKeyDown(event: KeyboardEvent) {
      if (isTypingInField(event.target)) {
        return
      }

      const {
        isLoading,
        append,
        clear,
        deleteLast,
        evaluate,
      } = actionsRef.current

      const { key } = event

      if (isLoading && (key === 'Enter' || key === '=')) {
        event.preventDefault()
        return
      }

      if (/^[0-9]$/.test(key)) {
        event.preventDefault()
        append(key)
        return
      }

      switch (key) {
        case '.':
          event.preventDefault()
          append('.')
          return
        case '+':
        case '-':
        case '*':
        case '/':
        case '^':
        case '%':
        case '(':
        case ')':
          event.preventDefault()
          append(key)
          return
        case 'x':
        case 'X':
          event.preventDefault()
          append('*')
          return
        case 'Enter':
        case '=':
          event.preventDefault()
          evaluate()
          return
        case 'Backspace':
          event.preventDefault()
          deleteLast()
          return
        case 'Escape':
          event.preventDefault()
          clear()
          return
        default:
          break
      }
    }

    window.addEventListener('keydown', handleKeyDown)
    return () => {
      window.removeEventListener('keydown', handleKeyDown)
    }
  }, [])
}
