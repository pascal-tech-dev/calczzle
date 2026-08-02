import { Calculator } from '../features/calculator/components/Calculator'

export function App() {
  return (
    <main>
      <header className="appHeader">
        <h1 className="brand">Calczzle</h1>
        <p className="tagline">Expression calculator powered by the backend</p>
      </header>

      <Calculator />

      <section className="usageGuide" aria-label="Usage guide">
        <h2>How to use</h2>
        <ul>
          <li>
            Build an expression with the keypad, then press <kbd>=</kbd> or{' '}
            <kbd>Enter</kbd>.
          </li>
          <li>
            Supported: <code>+ - * / ^</code>, parentheses, <code>%</code>, and{' '}
            <code>sqrt(</code>.
          </li>
          <li>
            Example: <code>2^-3</code>, <code>sqrt(16)+2</code>,{' '}
            <code>100*15%</code>
          </li>
          <li>
            Keyboard: digits and operators, <kbd>Backspace</kbd> delete,{' '}
            <kbd>Esc</kbd> clear, <kbd>x</kbd> for multiply.
          </li>
          <li>Math is evaluated on the server — the UI never uses eval.</li>
        </ul>
      </section>
    </main>
  )
}
