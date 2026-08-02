# Frontend Implementation Phases

Extracted from `docs/frontend-guide.md`. Implement **one phase at a time**. Do not start the next phase until the current phase acceptance criteria pass.

| Phase | Title | Status |
|------:|-------|--------|
| 1 | Verify and clean the Vite project | done |
| 2 | Install test dependencies and configure scripts | done |
| 3 | Create the frontend folder structure | done |
| 4 | Define API and calculator TypeScript types | pending |
| 5 | Implement the calculator API client | pending |
| 6 | Define calculator keys and expression helpers | pending |
| 7 | Implement calculator state with `useReducer` | pending |
| 8 | Implement the `useCalculator` hook | pending |
| 9 | Build presentational components | pending |
| 10 | Connect UI interactions to backend evaluation | pending |
| 11 | Add keyboard support and accessibility | pending |
| 12 | Add responsive styling | pending |
| 13 | Add loading and error behavior | pending |
| 14 | Write unit and component tests | pending |
| 15 | Connect frontend and backend Dev Containers | pending |
| 16 | Add the production Docker setup | pending |
| 17 | Run final quality checks and update README | pending |

## Milestones (optional grouping)

1. **Static skeleton** — folders + display/keypad shell (no backend)
2. **Local expression building** — reducer + key clicks
3. **Backend evaluation** — `POST /api/v1/evaluate`
4. **Error and loading states**
5. **Tests**
6. **Design and production** — CSS, a11y, Docker, README

## Hard rules (from the guide)

- Never evaluate expressions in the frontend (`eval`, `Function`, or libraries)
- Backend is the source of truth for validity and results
- No Redux, Zustand, Axios, React Router, Tailwind, or UI frameworks
- One phase at a time; verify before continuing
