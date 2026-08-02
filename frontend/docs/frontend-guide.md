# Sezzle Calculator — Frontend Implementation Guide

> **Stack:** React + Vite + TypeScript  
> **Backend:** Go + Fiber expression evaluation API  
> **Goal:** Build a clean, responsive, testable calculator UI that sends expressions to the backend.  
> **Important rule:** The frontend must **never evaluate expressions itself** and must never use JavaScript `eval()`.

---

## 1. Backend contract assumed by the frontend

This guide assumes that the completed backend exposes the following endpoint.

### Evaluate expression

```http
POST /api/v1/evaluate
Content-Type: application/json
```

Request:

```json
{
  "expression": "sqrt(16) + 2^-3"
}
```

Success response:

```json
{
  "result": 4.125
}
```

Error response:

```json
{
  "error": {
    "code": "INVALID_EXPRESSION",
    "message": "The expression is invalid."
  }
}
```

The exact error codes may differ in the backend. The TypeScript types should match the real backend response before implementation begins.

### Supported expression features

The UI should allow the user to build expressions containing:

- Numbers and decimals
- `+`, `-`, `*`, `/`
- Exponentiation: `^`
- Parentheses: `(` and `)`
- Square root: `sqrt(...)`
- Percentage: `%`
- Unary minus, including expressions such as `-5`, `3*-2`, and `2^-3`

---

## 2. Frontend architecture decisions

Use the following principles:

1. Use React local state with `useReducer`; do not install Redux or Zustand.
2. Use the native `fetch` API; do not install Axios for one endpoint.
3. Keep API communication outside React components.
4. Keep expression-building behavior outside visual components.
5. Treat the backend as the source of truth for expression validity and calculation results.
6. Perform only lightweight UX checks in the frontend.
7. Use CSS Modules for feature-level styling and one global stylesheet.
8. Use Vitest and React Testing Library for tests.
9. Use a Vite development proxy so the browser calls `/api/...` without CORS configuration.
10. Do not introduce React Router because the application has only one page.

---

## 3. Implementation order

Follow these phases in order:

```text
Phase 1  — Verify and clean the Vite project
Phase 2  — Install test dependencies and configure scripts
Phase 3  — Create the frontend folder structure
Phase 4  — Define API and calculator TypeScript types
Phase 5  — Implement the calculator API client
Phase 6  — Define calculator keys and expression helpers
Phase 7  — Implement calculator state with useReducer
Phase 8  — Implement the useCalculator hook
Phase 9  — Build presentational components
Phase 10 — Connect UI interactions to backend evaluation
Phase 11 — Add keyboard support and accessibility
Phase 12 — Add responsive styling
Phase 13 — Add loading and error behavior
Phase 14 — Write unit and component tests
Phase 15 — Connect frontend and backend Dev Containers
Phase 16 — Add the production Docker setup
Phase 17 — Run final quality checks and update README
```

Do not begin with styling. First make the complete data flow work with a simple UI.

---

# Phase 1 — Verify and clean the Vite project

The React + TypeScript Vite application is already scaffolded. From the frontend Dev Container:

```bash
cd /workspace/frontend
node --version
npm --version
npm install
npm run dev -- --host 0.0.0.0
```

Open:

```text
http://localhost:5173
```

Remove the default Vite demo files and assets that will not be used:

```text
src/assets/react.svg
public/vite.svg
```

Replace the demo `App.tsx` and demo CSS with a minimal application shell.

### Phase acceptance criteria

- The frontend Dev Container opens successfully.
- `npm install` succeeds.
- Vite is accessible at `http://localhost:5173`.
- The default counter demo is removed.

---

# Phase 2 — Install test dependencies

Install Vitest, jsdom, React Testing Library, jest-dom matchers, user-event, and coverage support:

```bash
npm install --save-dev \
  vitest \
  jsdom \
  @vitest/coverage-v8 \
  @testing-library/react \
  @testing-library/jest-dom \
  @testing-library/user-event
```

Add or verify these scripts in `package.json`:

```json
{
  "scripts": {
    "dev": "vite",
    "build": "tsc -b && vite build",
    "lint": "eslint .",
    "preview": "vite preview",
    "test": "vitest",
    "test:run": "vitest run",
    "test:coverage": "vitest run --coverage"
  }
}
```

Create `vitest.config.ts`:

```ts
import react from '@vitejs/plugin-react'
import { defineConfig } from 'vitest/config'

export default defineConfig({
  plugins: [react()],
  test: {
    environment: 'jsdom',
    setupFiles: ['./src/test/setup.ts'],
    css: true,
    clearMocks: true,
  },
})
```

Create `src/test/setup.ts`:

```ts
import '@testing-library/jest-dom/vitest'
```

### Phase acceptance criteria

```bash
npm run lint
npm run test:run
npm run build
```

All three commands must complete successfully, even if there are no meaningful tests yet.

---

# Phase 3 — Create the frontend folder structure

Use a feature-based structure:

```text
frontend/
├── src/
│   ├── app/
│   │   └── App.tsx
│   │
│   ├── features/
│   │   └── calculator/
│   │       ├── api/
│   │       │   ├── calculatorApi.ts
│   │       │   └── calculatorApi.test.ts
│   │       │
│   │       ├── components/
│   │       │   ├── Calculator.tsx
│   │       │   ├── Calculator.test.tsx
│   │       │   ├── CalculatorDisplay.tsx
│   │       │   ├── CalculatorKeypad.tsx
│   │       │   └── CalculatorButton.tsx
│   │       │
│   │       ├── hooks/
│   │       │   └── useCalculator.ts
│   │       │
│   │       ├── model/
│   │       │   ├── calculator.types.ts
│   │       │   ├── calculator.reducer.ts
│   │       │   ├── calculator.reducer.test.ts
│   │       │   ├── calculatorKeys.ts
│   │       │   └── expressionBuilder.ts
│   │       │
│   │       └── styles/
│   │           └── Calculator.module.css
│   │
│   ├── styles/
│   │   └── global.css
│   │
│   ├── test/
│   │   └── setup.ts
│   │
│   ├── main.tsx
│   └── vite-env.d.ts
│
├── .env.example
├── vite.config.ts
├── vitest.config.ts
├── package.json
├── tsconfig.json
└── Dockerfile
```

### Why this structure?

```text
api/
→ HTTP request and response handling

components/
→ Visual React components

hooks/
→ Connects state, API calls, and UI actions

model/
→ Pure TypeScript state, key definitions, and expression rules

styles/
→ Calculator-specific CSS
```

Keep `App.tsx` small. It should compose the page, not implement calculator behavior.

---

# Phase 4 — Define TypeScript types

Create `src/features/calculator/model/calculator.types.ts`.

## API types

```ts
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
```

## Calculator state

```ts
export type CalculatorStatus =
  | 'idle'
  | 'loading'
  | 'success'
  | 'error'

export interface CalculatorState {
  expression: string
  result: number | null
  status: CalculatorStatus
  error: string | null
  justEvaluated: boolean
}
```

Do not store duplicate derived values such as both `displayExpression` and `expression`. Derive the displayed value from the existing state.

## Reducer actions

```ts
export type CalculatorAction =
  | { type: 'APPEND'; value: string }
  | { type: 'DELETE_LAST' }
  | { type: 'CLEAR' }
  | { type: 'EVALUATE_START' }
  | { type: 'EVALUATE_SUCCESS'; result: number }
  | { type: 'EVALUATE_FAILURE'; message: string }
```

## Calculator key types

```ts
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
```

### Phase acceptance criteria

- All shared calculator types are centralized.
- Components do not redefine API response types.
- TypeScript compilation succeeds.

---

# Phase 5 — Implement the API client

Create `src/features/calculator/api/calculatorApi.ts`.

The browser should call a relative URL:

```text
/api/v1/evaluate
```

Do not hardcode `http://localhost:8080` inside React code.

## Custom API error

```ts
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
```

## Evaluate request

```ts
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
      'The calculator service is unavailable.',
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
```

### Responsibilities of the API file

- Serialize the request.
- Send the HTTP request.
- Parse success responses.
- Parse backend error responses.
- Convert network and malformed-response failures into a predictable error type.

### What must not be in this file?

- React state
- Button behavior
- Expression parsing
- UI rendering

---

# Phase 6 — Define calculator keys and expression helpers

## Calculator key definitions

Create `src/features/calculator/model/calculatorKeys.ts`.

Suggested layout:

```text
C    ⌫    (    )
√    ^    %    ÷
7    8    9    ×
4    5    6    −
1    2    3    +
0    .         =
```

Use display labels that are easy to read but store backend-compatible values:

```ts
import type { CalculatorKey } from './calculator.types'

export const calculatorKeys: readonly CalculatorKey[] = [
  {
    id: 'clear',
    label: 'C',
    kind: 'action',
    action: 'clear',
    ariaLabel: 'Clear expression',
  },
  {
    id: 'delete',
    label: '⌫',
    kind: 'action',
    action: 'delete',
    ariaLabel: 'Delete last character',
  },
  {
    id: 'left-parenthesis',
    label: '(',
    value: '(',
    kind: 'parenthesis',
    ariaLabel: 'Left parenthesis',
  },
  {
    id: 'right-parenthesis',
    label: ')',
    value: ')',
    kind: 'parenthesis',
    ariaLabel: 'Right parenthesis',
  },
  {
    id: 'square-root',
    label: '√',
    value: 'sqrt(',
    kind: 'function',
    ariaLabel: 'Square root',
  },
  {
    id: 'power',
    label: '^',
    value: '^',
    kind: 'operator',
    ariaLabel: 'Exponentiation',
  },
  {
    id: 'percentage',
    label: '%',
    value: '%',
    kind: 'percentage',
    ariaLabel: 'Percentage',
  },
  {
    id: 'divide',
    label: '÷',
    value: '/',
    kind: 'operator',
    ariaLabel: 'Divide',
  },
  // Add digits, multiplication, subtraction, addition, decimal, and evaluate.
]
```

The frontend state should store canonical operators:

```text
UI label × → expression value *
UI label ÷ → expression value /
UI label − → expression value -
UI label √ → expression value sqrt(
```

This avoids Unicode-normalization problems at the API boundary.

## Expression helper

Create `expressionBuilder.ts` for small, pure UI rules.

Suggested responsibilities:

- Append a selected key value.
- Delete the final character.
- Decide what happens after a successful evaluation.
- Optionally prevent duplicate decimal points in the current number.
- Optionally prevent an expression from exceeding the same maximum length used by the backend.

Do **not** reimplement the backend parser here.

### Post-evaluation behavior

After evaluating `3+4` and receiving `7`:

- Pressing a digit should start a new expression.
- Pressing `sqrt(` should start a new expression.
- Pressing an operator should continue from the result: `7+`.
- Pressing `%` may continue from the result: `7%`.
- Pressing clear should reset everything.

Implement this in one pure helper or in the reducer, not separately in several components.

---

# Phase 7 — Implement calculator state with `useReducer`

Create `calculator.reducer.ts`.

## Initial state

```ts
import type {
  CalculatorAction,
  CalculatorState,
} from './calculator.types'

export const initialCalculatorState: CalculatorState = {
  expression: '',
  result: null,
  status: 'idle',
  error: null,
  justEvaluated: false,
}
```

## Reducer responsibilities

```text
APPEND
→ Adds a value to the expression.
→ Clears an old error.
→ Applies the post-evaluation behavior.

DELETE_LAST
→ Removes the final expression character.
→ Does nothing while loading.

CLEAR
→ Returns the complete initial state.

EVALUATE_START
→ Sets loading state and clears an old error.

EVALUATE_SUCCESS
→ Stores result and marks the expression as evaluated.

EVALUATE_FAILURE
→ Stores a user-visible error without deleting the expression.
```

Keep the reducer pure:

```text
No fetch
No timers
No DOM access
No localStorage
No mutation
```

The reducer must return a new state object for each transition.

### Why `useReducer`?

The calculator has related state transitions involving expression, result, loading, errors, and post-evaluation behavior. A reducer keeps these transitions centralized and testable.

---

# Phase 8 — Implement `useCalculator`

Create `src/features/calculator/hooks/useCalculator.ts`.

This hook is the main frontend orchestration layer.

It should:

1. Hold reducer state.
2. Expose display values.
3. Handle calculator key actions.
4. Call `evaluateExpression`.
5. Convert API errors into reducer actions.
6. Prevent duplicate requests while loading.
7. Optionally cancel an in-flight request during unmount.

Suggested public interface:

```ts
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
```

Suggested evaluation flow:

```text
User presses =
      ↓
Reject empty/whitespace-only expression in UI
      ↓
Dispatch EVALUATE_START
      ↓
POST /api/v1/evaluate
      ↓
Success → EVALUATE_SUCCESS
Failure → EVALUATE_FAILURE
```

Do not trim or rewrite meaningful expression characters. Only use `trim()` to check whether the expression is empty.

Example error handling:

```ts
try {
  const result = await evaluateExpression(state.expression)
  dispatch({ type: 'EVALUATE_SUCCESS', result })
} catch (error) {
  const message =
    error instanceof CalculatorApiError
      ? error.message
      : 'An unexpected error occurred.'

  dispatch({ type: 'EVALUATE_FAILURE', message })
}
```

Be careful with stale state in async functions. The expression sent to the backend should be captured before awaiting:

```ts
const expressionToEvaluate = state.expression
```

---

# Phase 9 — Build presentational components

## `CalculatorButton.tsx`

A reusable button component.

Responsibilities:

- Render a semantic `<button>`.
- Receive label, `aria-label`, disabled state, and click handler.
- Apply CSS classes based on key type.
- Never know about API calls or reducer actions.

## `CalculatorDisplay.tsx`

Responsibilities:

- Show the current expression.
- Show the result after successful evaluation.
- Show a loading indicator without removing the expression.
- Show an accessible error message.

Recommended accessibility attributes:

```tsx
<div aria-live="polite" aria-atomic="true">
  {/* display value */}
</div>

<p role="alert">
  {/* error message */}
</p>
```

Use a readable fallback such as `0` when the expression is empty.

## `CalculatorKeypad.tsx`

Responsibilities:

- Render keys from `calculatorKeys`.
- Translate each key into `append`, `clear`, `deleteLast`, or `evaluate` calls.
- Disable interactive keys when needed.

It must not contain calculation logic.

## `Calculator.tsx`

This is the feature container.

Responsibilities:

- Call `useCalculator()`.
- Render `CalculatorDisplay`.
- Render `CalculatorKeypad`.
- Wire hook actions into presentational components.
- Register keyboard behavior through a dedicated hook or effect.

## `App.tsx`

Keep it small:

```tsx
import { Calculator } from '../features/calculator/components/Calculator'

export function App() {
  return (
    <main>
      <Calculator />
    </main>
  )
}
```

---

# Phase 10 — Complete frontend application flow

The final runtime flow must be:

```text
User presses calculator key
        ↓
CalculatorKeypad
        ↓
useCalculator action
        ↓
calculatorReducer
        ↓
React renders updated expression
        ↓
User presses =
        ↓
useCalculator.evaluate()
        ↓
calculatorApi.evaluateExpression()
        ↓
POST /api/v1/evaluate
        ↓
Go/Fiber backend
        ↓
Tokenizer → parser → evaluator
        ↓
JSON result or JSON error
        ↓
calculatorApi
        ↓
calculatorReducer
        ↓
CalculatorDisplay
```

Example:

```text
Button sequence: 2, ^, -, 3, =
Frontend expression: "2^-3"
Backend result: 0.125
Display: 0.125
```

The frontend must not know how unary minus or operator precedence works.

---

# Phase 11 — Keyboard support and accessibility

Support these keyboard inputs:

```text
0-9        → digits
.          → decimal
+ - * / ^  → operators
( )        → parentheses
%          → percentage
Enter or = → evaluate
Backspace  → delete last
Escape     → clear
```

Optional aliases:

```text
x or X → multiplication
```

Do not map arbitrary letters to functions. `sqrt(` can be inserted through the square-root button initially.

Keyboard handling rules:

- Ignore unsupported keys.
- Call `preventDefault()` only for keys handled by the calculator.
- Ignore actions while a request is loading where appropriate.
- Remove window listeners during cleanup.
- Button controls must remain usable without a keyboard.

Accessibility checklist:

- Use real `<button>` elements.
- Every symbolic button has an explicit `aria-label`.
- Loading state is announced with `aria-live` or `role="status"`.
- Errors use `role="alert"`.
- Visible focus states are not removed.
- Text contrast is sufficient.
- The UI remains usable at 200% zoom.

---

# Phase 12 — Responsive styling

Use CSS Modules in:

```text
src/features/calculator/styles/Calculator.module.css
```

Use global reset and page layout in:

```text
src/styles/global.css
```

## Layout goals

- Center the calculator on desktop.
- Fit comfortably on mobile screens.
- Use a four-column keypad grid.
- Give the zero or equals button a wider grid span if desired.
- Keep touch targets at least approximately 44px high and wide.
- Prevent long expressions from breaking the layout.

Suggested calculator width:

```css
.calculator {
  width: min(100%, 24rem);
}
```

Suggested keypad:

```css
.keypad {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 0.75rem;
}
```

Display overflow behavior:

```css
.display {
  overflow-x: auto;
  white-space: nowrap;
  text-align: right;
}
```

Use CSS variables for repeated design values, but avoid building a large design system for this assignment.

Do not add Tailwind, Material UI, or another component framework unless the assignment explicitly requires it.

---

# Phase 13 — Loading and error behavior

## Loading

When evaluation starts:

- Keep the current expression visible.
- Disable the equals button.
- Prevent duplicate evaluate requests.
- Show a subtle progress state such as `Calculating…`.
- Do not clear the result or expression prematurely.

## Backend validation errors

Examples:

```text
3 + * 4
10 / 0
sqrt(-4)
(3 + 4
```

Show the backend message near the display.

Do not convert every backend error into the same generic message. Preserve safe, user-facing messages returned by the API.

## Network errors

Show a stable message such as:

```text
The calculator service is unavailable. Please try again.
```

Do not display raw JavaScript stack traces or fetch error details.

## Error reset behavior

Clear the old error when:

- The user modifies the expression.
- The user clears the calculator.
- A new evaluation starts.

Keep the invalid expression visible so the user can correct it.

---

# Phase 14 — Testing plan

Use a test pyramid focused on pure logic and user-visible behavior.

## A. Reducer tests

Test `calculator.reducer.ts` with table-driven or individual tests:

```text
APPEND adds a value.
CLEAR resets all state.
DELETE_LAST removes one character.
EVALUATE_START sets loading state.
EVALUATE_SUCCESS stores the result.
EVALUATE_FAILURE stores the error.
Digit after evaluation starts a new expression.
Operator after evaluation continues from the result.
```

These tests should require no React rendering.

## B. API client tests

Mock `globalThis.fetch` and test:

```text
200 response returns result.
400/422 response throws CalculatorApiError with backend message.
500 non-JSON response creates a generic API error.
Rejected fetch creates NETWORK_ERROR.
Non-finite or malformed result creates INVALID_API_RESPONSE.
```

## C. Component tests

Use React Testing Library and `userEvent`.

Important tests:

```text
Calculator initially displays 0.
Clicking 3, +, 4 builds "3+4".
Clicking clear resets the expression.
Clicking delete removes the final character.
Clicking square root appends "sqrt(".
Clicking multiplication stores "*" even though the label is "×".
Clicking equals sends the current expression.
Successful response displays the result.
Backend error is shown with role="alert".
Equals is disabled while loading.
```

Prefer user-facing queries:

```ts
screen.getByRole('button', { name: /multiply/i })
screen.getByRole('alert')
screen.getByText('0.125')
```

Use `data-testid` only when semantic queries are impractical.

## D. Keyboard tests

Test:

```text
Typing 2^-3 and Enter evaluates the expression.
Backspace deletes the last character.
Escape clears the calculator.
Unsupported keys are ignored.
```

## E. Optional end-to-end test

Only add Playwright if time remains. It is optional for this assignment because backend unit tests plus frontend component tests provide better value first.

### Test commands

```bash
npm run test
npm run test:run
npm run test:coverage
```

Do not optimize for a vanity coverage number. Cover state transitions, API failures, and critical user behavior.

---

# Phase 15 — Connect Vite to the backend Dev Container

Use Vite’s development proxy. This avoids CORS configuration and keeps frontend API calls relative.

Update `vite.config.ts`:

```ts
import react from '@vitejs/plugin-react'
import { defineConfig, loadEnv } from 'vite'

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '')

  return {
    plugins: [react()],
    server: {
      host: '0.0.0.0',
      port: 5173,
      proxy: {
        '/api': {
          target: env.BACKEND_PROXY_TARGET || 'http://backend:8080',
          changeOrigin: true,
        },
      },
    },
  }
})
```

Create `.env.example`:

```dotenv
BACKEND_PROXY_TARGET=http://backend:8080
```

In the development Compose file, the backend service is named `backend`, so the Vite process can reach it through:

```text
http://backend:8080
```

The browser still calls:

```text
http://localhost:5173/api/v1/evaluate
```

Vite proxies the request internally:

```text
Browser
  ↓ localhost:5173/api/v1/evaluate
Vite dev server
  ↓ backend:8080/api/v1/evaluate
Go/Fiber backend
```

Run both services and test:

```bash
curl http://localhost:8080/health
npm run dev -- --host 0.0.0.0
```

Then evaluate an expression through the UI.

### Integration acceptance cases

```text
3 + 4                  → 7
3 + 4 - 7 * 2          → -7
(3 + 4) * 2            → 14
2 ^ -3                 → 0.125
sqrt(81)               → 9
100 * 15%              → 15
10 / 0                 → visible API error
3 + * 4                → visible API error
```

---

# Phase 16 — Production Docker setup

Do this only after the frontend works locally.

Use a multi-stage Dockerfile:

```text
Stage 1: Node image
→ npm ci
→ npm run build
→ creates /app/dist

Stage 2: nginx or another static web server
→ copy dist files
→ serve the React application
```

The production web server should:

1. Serve the Vite `dist` directory.
2. Fall back to `index.html` for SPA paths.
3. Proxy `/api/` requests to the backend service.

The browser should continue to call the same relative endpoint:

```text
/api/v1/evaluate
```

Do not bake a localhost backend URL into the production JavaScript bundle.

Example production request flow:

```text
Browser
  ↓ /api/v1/evaluate
Frontend web server / reverse proxy
  ↓ http://backend:8080/api/v1/evaluate
Go backend
```

---

# Phase 17 — Final quality checks

Run:

```bash
npm ci
npm run lint
npm run test:run
npm run test:coverage
npm run build
```

Then verify the generated production build:

```bash
npm run preview -- --host 0.0.0.0
```

## Manual acceptance checklist

- [ ] Calculator loads without console errors.
- [ ] All numeric keys work.
- [ ] Decimal input works.
- [ ] Operators display correctly and send canonical values.
- [ ] Parentheses work.
- [ ] `sqrt(` can be inserted.
- [ ] Percentage can be inserted.
- [ ] Unary-minus expressions can be sent.
- [ ] Loading state prevents duplicate requests.
- [ ] Backend validation errors are visible.
- [ ] Network failures show a stable message.
- [ ] Keyboard interaction works.
- [ ] Buttons are accessible by keyboard.
- [ ] Mobile layout is usable.
- [ ] `npm run lint` passes.
- [ ] `npm run test:run` passes.
- [ ] `npm run build` passes.
- [ ] Dockerized frontend can reach the backend.

---

# Recommended implementation milestones

## Milestone 1 — Static skeleton

Create the folders and render a calculator shell with a display and keypad.

**Do not call the backend yet.**

Result:

```text
App → Calculator → Display + Keypad
```

## Milestone 2 — Local expression building

Implement the reducer and key clicks so the display can build:

```text
sqrt(16)+2^-3
```

Do not calculate the result locally.

## Milestone 3 — Backend evaluation

Connect the equals action to:

```text
POST /api/v1/evaluate
```

Show the returned result.

## Milestone 4 — Error and loading states

Handle:

```text
Invalid syntax
Division by zero
Negative square root
Backend unavailable
Duplicate evaluate clicks
```

## Milestone 5 — Tests

Test reducer behavior, API behavior, and user interactions.

## Milestone 6 — Design and production

Finish responsive CSS, accessibility, Docker, and README documentation.

---

# Rules for Cursor while implementing

Give Cursor these constraints together with this document:

1. Implement one phase at a time.
2. Do not change the backend API contract without explicit approval.
3. Do not use `eval()`, `Function()`, or any frontend expression-evaluation library.
4. Do not calculate expressions in the frontend.
5. Do not add Redux, Zustand, Axios, React Router, Tailwind, or a UI framework.
6. Keep React components small and focused on rendering.
7. Keep state transitions in the reducer.
8. Keep HTTP behavior in `calculatorApi.ts`.
9. Keep expression-building UX rules in pure TypeScript helpers.
10. Use strict TypeScript types; avoid `any`.
11. Use semantic HTML and accessible labels.
12. Run lint, tests, and build after every milestone.
13. Do not create abstractions that are used only once unless they improve testability or readability.
14. Preserve backend error messages that are safe for users.
15. Ask before introducing a new production dependency.

---

# Final architecture summary

```text
App
 ↓
Calculator
 ↓
useCalculator
 ├── calculatorReducer
 └── calculatorApi
          ↓
   POST /api/v1/evaluate
          ↓
      Go backend
```

Responsibility boundaries:

```text
Components
→ Render UI and forward user actions

useCalculator
→ Coordinates state and asynchronous API calls

Reducer
→ Controls deterministic calculator state transitions

Expression helpers
→ Handle small frontend input-building rules

API client
→ Communicates with Go/Fiber backend

Backend
→ Validates, parses, and evaluates expressions
```

This separation keeps the frontend easy to test, prevents duplicated calculator logic, and clearly demonstrates a maintainable full-stack architecture.

---

# Official references

- Vite Getting Started: https://vite.dev/guide/
- React Managing State: https://react.dev/learn/managing-state
- React Sharing State: https://react.dev/learn/sharing-state-between-components
- Vitest Guide: https://vitest.dev/guide/
- React Testing Library: https://testing-library.com/docs/react-testing-library/intro/
- Testing Library Queries: https://testing-library.com/docs/queries/about/
