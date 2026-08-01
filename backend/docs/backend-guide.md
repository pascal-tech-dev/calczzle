# Calczzle Backend Implementation Guide

## Purpose

This document is the implementation plan and architectural context for the Calczzle backend.

Use it as a reference when generating or modifying code in Cursor. Follow the phases in order and preserve the dependency boundaries described below.

## Project Goal

Build a Go backend microservice for a full-stack calculator application.

The backend must:

- Expose a REST API for evaluating arithmetic expressions.
- Support addition, subtraction, multiplication, and division.
- Optionally support exponentiation, square root, and percentage.
- Validate requests and expression syntax.
- Handle mathematical errors such as division by zero.
- Return consistent JSON responses.
- Use clean, maintainable, and testable architecture.

## Technology Decisions

- Language: Go
- HTTP framework: Fiber v3
- API style: REST
- Architecture: feature-first packages with controller, service, and expression-engine layers
- Testing: Go standard testing package
- Configuration: environment variables using the Go standard library
- Database: not required

Do not introduce a database, ORM, Redis, authentication, dependency-injection framework, or third-party validation library unless a real requirement appears.

---

## Implementation Order

Implement the backend in this order:

1. Verify the development environment.
2. Initialize the Go module.
3. Install Fiber.
4. Create the project structure.
5. Bootstrap the Fiber application.
6. Add a health endpoint.
7. Define the evaluation API contract.
8. Create controller and service boundaries.
9. Implement the expression engine.
10. Connect the complete request flow.
11. Add validation and error mapping.
12. Add tests and coverage.
13. Finalize Docker and documentation.

Do not begin with the tokenizer or parser. First establish the application skeleton and architectural boundaries.

---

## Phase 1: Verify the Development Environment

The backend runs inside a Dev Container, so Go does not need to be installed directly on the host machine.

Run:

```bash
go version
git --version
make --version
curl --version
pwd
```

Expected working directory:

```text
/workspace/backend
```

Useful Go commands:

```bash
go fmt ./...
go test ./...
go test ./... -cover
go vet ./...
go mod tidy
```

---

## Phase 2: Initialize the Go Module

Use the selected local module name:

```bash
go mod init pascal-tech/calczzle-backend
```

All internal imports must begin with the exact module path declared in `go.mod`.

Example:

```go
import "pascal-tech/calczzle-backend/internal/calculator/service"
```

Verify the module:

```bash
cat go.mod
```

Expected shape:

```go
module pascal-tech/calczzle-backend

go 1.26
```

Use the Go version actually installed in the Dev Container if it differs.

---

## Phase 3: Install Fiber

Install Fiber v3:

```bash
go get github.com/gofiber/fiber/v3
go mod tidy
```

Fiber import path:

```go
import "github.com/gofiber/fiber/v3"
```

Initially use only middleware with a clear purpose:

- Logger middleware
- Recover middleware
- CORS middleware only when the frontend proxy is insufficient

---

## Phase 4: Project Structure

Use this structure:

```text
backend/
├── cmd/
│   └── api/
│       └── main.go
│
├── internal/
│   ├── app/
│   │   └── app.go
│   │
│   ├── config/
│   │   └── config.go
│   │
│   ├── health/
│   │   ├── controller.go
│   │   └── controller_test.go
│   │
│   ├── calculator/
│   │   ├── controller/
│   │   │   ├── controller.go
│   │   │   ├── routes.go
│   │   │   ├── request.go
│   │   │   ├── response.go
│   │   │   └── controller_test.go
│   │   │
│   │   ├── service/
│   │   │   ├── interface.go
│   │   │   ├── service.go
│   │   │   └── service_test.go
│   │   │
│   │   └── expression/
│   │       ├── token.go
│   │       ├── tokenizer.go
│   │       ├── parser.go
│   │       ├── evaluator.go
│   │       ├── operator.go
│   │       ├── function.go
│   │       ├── errors.go
│   │       ├── tokenizer_test.go
│   │       ├── parser_test.go
│   │       └── evaluator_test.go
│   │
│   └── platform/
│       └── httpx/
│           ├── error.go
│           └── error_handler.go
│
├── .dockerignore
├── Dockerfile
├── go.mod
└── go.sum
```

Create the directories:

```bash
mkdir -p \
  cmd/api \
  internal/app \
  internal/config \
  internal/health \
  internal/calculator/controller \
  internal/calculator/service \
  internal/calculator/expression \
  internal/platform/httpx
```

---

## Package Responsibilities

### `cmd/api`

Executable entry point.

Responsibilities:

- Load configuration.
- Construct the application.
- Start the HTTP server.
- Handle startup failures.

It must not contain routes, request validation, expression parsing, or arithmetic logic.

### `internal/config`

Reads environment variables and provides application configuration.

Initial configuration:

```go
type Config struct {
    Port string
}
```

Default port:

```text
8080
```

Use `os.Getenv`; do not add an environment-library dependency for one variable.

### `internal/app`

Application composition root.

Responsibilities:

- Create the Fiber application.
- Register middleware.
- Construct services and controllers.
- Register routes.
- Configure the global error handler.

Dependency wiring belongs here, not in `main.go` or route files.

### `internal/health`

Infrastructure health endpoint.

Endpoint:

```http
GET /health
```

Response:

```json
{
  "status": "ok"
}
```

### `internal/calculator/controller`

Fiber HTTP layer.

Responsibilities:

- Receive HTTP requests.
- Bind request bodies.
- Perform transport-level validation.
- Call the calculator service.
- Return JSON responses.

It must not tokenize expressions, apply precedence rules, or perform arithmetic.

### `internal/calculator/service`

Use-case orchestration layer.

Public behavior:

```go
Evaluate(expression string) (float64, error)
```

Responsibilities:

- Coordinate tokenizer, parser, and evaluator.
- Return domain results and domain errors.

It must not know about Fiber, HTTP routes, JSON, or status codes.

### `internal/calculator/expression`

Framework-independent calculation engine.

Stages:

```text
Tokenizer -> Parser -> Evaluator
```

It must not import Fiber or HTTP packages.

### `internal/platform/httpx`

Shared HTTP concerns.

Responsibilities:

- API error response models.
- Global Fiber error handler.
- Mapping domain errors to HTTP status codes and public error codes.

---

## Dependency Direction

The dependency direction must remain:

```text
cmd/api
   ↓
internal/app
   ↓
controller
   ↓
service
   ↓
expression
```

Rules:

- `expression` must not import `service`.
- `service` must not import `controller`.
- `controller` may depend on a service interface.
- HTTP concepts must not leak into the expression engine.

---

## Phase 5: Bootstrap the Application

Implement only these components first:

- `config.Load()`
- `app.New()`
- `cmd/api/main.go`
- `GET /health`
- Health endpoint test

Conceptual `main.go`:

```go
func main() {
    cfg := config.Load()
    server := app.New(cfg)

    if err := server.Listen(":" + cfg.Port); err != nil {
        log.Fatal(err)
    }
}
```

Verify:

```bash
go fmt ./...
go vet ./...
go test ./...
go run ./cmd/api
```

Then test:

```bash
curl -i http://localhost:8080/health
```

Do not continue until the health endpoint works and all tests pass.

---

## Phase 6: API Contract

Use one evaluation endpoint:

```http
POST /api/v1/evaluate
Content-Type: application/json
```

Request:

```json
{
  "expression": "3 + 4 - 7 * 2"
}
```

Successful response:

```json
{
  "result": -7
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

Request model:

```go
type EvaluateRequest struct {
    Expression string `json:"expression"`
}
```

Response model:

```go
type EvaluateResponse struct {
    Result float64 `json:"result"`
}
```

Do not return internal tokens, postfix output, stack state, or parser implementation details.

---

## Phase 7: Service Boundary

Define the controller-to-service interface first.

```go
package service

type Evaluator interface {
    Evaluate(expression string) (float64, error)
}
```

Service skeleton:

```go
package service

type Service struct{}

func New() *Service {
    return &Service{}
}

func (s *Service) Evaluate(expression string) (float64, error) {
    return 0, nil
}
```

Controller dependency:

```go
type Controller struct {
    evaluator service.Evaluator
}

func New(evaluator service.Evaluator) *Controller {
    return &Controller{evaluator: evaluator}
}
```

Create an interface only at this meaningful boundary. Do not create interfaces for every type.

---

## Phase 8: Connect HTTP Before Parser Logic

Before implementing the expression engine, make this flow work:

```text
POST /api/v1/evaluate
        ↓
Controller
        ↓
Service interface
        ↓
Temporary result
```

Temporarily return a fixed result from the service:

```go
func (s *Service) Evaluate(expression string) (float64, error) {
    return 42, nil
}
```

Test:

```bash
curl -X POST http://localhost:8080/api/v1/evaluate \
  -H "Content-Type: application/json" \
  -d '{"expression":"3 + 4"}'
```

Temporary response:

```json
{
  "result": 42
}
```

This verifies the HTTP architecture independently from parser complexity.

---

## Phase 9: Expression Engine

Implement the expression engine in this order.

### 9.1 Token Model

Token types:

- Number
- Operator
- Left parenthesis
- Right parenthesis
- Function
- Percentage

Suggested model:

```go
type TokenType int

const (
    TokenNumber TokenType = iota
    TokenOperator
    TokenLeftParenthesis
    TokenRightParenthesis
    TokenFunction
    TokenPercentage
)

type Token struct {
    Type     TokenType
    Value    string
    Position int
}
```

Keep `Position` for useful syntax errors.

### 9.2 Tokenizer

Input:

```text
3 + 4 - 7 * 2
```

Output conceptually:

```text
[number 3]
[operator +]
[number 4]
[operator -]
[number 7]
[operator *]
[number 2]
```

Initial supported syntax:

- Digits: `0-9`
- Decimal separator: `.`
- Operators: `+ - * / ^`
- Parentheses: `( )`
- Function: `sqrt`
- Postfix percentage: `%`
- Whitespace

The tokenizer recognizes tokens only. It must not calculate results.

### 9.3 Operator Definitions

```go
type Associativity int

const (
    LeftAssociative Associativity = iota
    RightAssociative
)

type Operator struct {
    Symbol        string
    Precedence    int
    Associativity Associativity
}
```

Precedence:

```text
+ -    precedence 1, left-associative
* /    precedence 2, left-associative
^      precedence 3, right-associative
```

### 9.4 Parser

Convert infix tokens to postfix notation using the shunting-yard algorithm.

Example:

```text
Infix:   3 + 4 - 7 * 2
Postfix: 3 4 + 7 2 * -
```

Responsibilities:

- Operator precedence
- Associativity
- Parentheses
- Functions
- Percentage placement
- Invalid token ordering

The parser must not calculate the final result.

### 9.5 Evaluator

Evaluate postfix tokens using a stack.

Example:

```text
3 4 + 7 2 * -
```

Evaluation:

```text
3 4 +  -> 7
7 2 *  -> 14
7 14 - -> -7
```

Support:

- Addition
- Subtraction
- Multiplication
- Division
- Exponentiation
- Square root
- Percentage

### 9.6 Connect the Service

```go
func (s *Service) Evaluate(input string) (float64, error) {
    tokens, err := expression.Tokenize(input)
    if err != nil {
        return 0, err
    }

    postfix, err := expression.ToPostfix(tokens)
    if err != nil {
        return 0, err
    }

    return expression.EvaluatePostfix(postfix)
}
```

---

## Complete Request Flow

```text
React frontend
      ↓
POST /api/v1/evaluate
      ↓
Fiber router
      ↓
Calculator controller
      ↓
Calculator service
      ↓
Tokenizer
      ↓
Parser
      ↓
Evaluator
      ↓
Calculator service
      ↓
Controller response model
      ↓
JSON response
```

Concrete example:

```text
Request:
{"expression":"3 + 4 - 7 * 2"}

Tokenizer:
[3, +, 4, -, 7, *, 2]

Parser:
[3, 4, +, 7, 2, *, -]

Evaluator:
-7

Response:
{"result":-7}
```

---

## Validation Responsibilities

### Controller Validation

Validate transport-level concerns:

- Valid JSON body
- `expression` field exists
- `expression` is a string
- Expression is not empty
- Request body is not excessively large

### Parser Validation

Validate expression syntax:

- Unknown characters
- Invalid operator sequences
- Missing operand
- Missing operator
- Unmatched parentheses
- Invalid function syntax
- Multiple decimal points

### Evaluator Validation

Validate mathematical execution:

- Division by zero
- Square root of a negative number
- Non-finite power result
- Stack underflow
- Extra operands
- `NaN` or infinite result

---

## Error Mapping

Use domain errors in the expression and service layers. Convert them to HTTP responses centrally in `internal/platform/httpx`.

Recommended mapping:

```text
Malformed JSON            -> 400 INVALID_REQUEST
Empty expression          -> 400 EMPTY_EXPRESSION
Invalid syntax            -> 400 INVALID_EXPRESSION
Unsupported function      -> 400 UNSUPPORTED_FUNCTION
Division by zero          -> 422 DIVISION_BY_ZERO
Negative square root      -> 422 INVALID_SQUARE_ROOT
Unexpected internal error -> 500 INTERNAL_SERVER_ERROR
```

Do not expose internal stack traces or parser internals to API clients.

---

## Testing Order

Write tests in this order:

1. Health controller
2. Calculator controller with fake service
3. Tokenizer
4. Parser
5. Evaluator
6. Calculator service integration
7. Complete HTTP endpoint

Important table-driven cases:

```text
3 + 4             = 7
3 + 4 - 7 * 2     = -7
(3 + 4 - 7) * 2   = 0
2 ^ 3             = 8
2 ^ 3 ^ 2         = 512
sqrt(81)          = 9
100 * 15%         = 15
10 / 0            = error
sqrt(-4)          = error
3 + * 4           = error
(3 + 4            = error
```

Run:

```bash
go fmt ./...
go vet ./...
go test ./...
go test ./... -cover
```

---

## Cursor Implementation Rules

When generating code for this project:

1. Implement only the current phase unless explicitly asked to continue.
2. Keep functions small and responsibilities focused.
3. Use idiomatic Go and standard-library features when sufficient.
4. Avoid premature abstractions and unnecessary dependencies.
5. Keep Fiber-specific code inside the HTTP/controller and application layers.
6. Keep the expression engine independent from Fiber and HTTP.
7. Return typed or sentinel domain errors rather than comparing arbitrary error strings.
8. Write tests alongside each implemented component.
9. Run formatting, vetting, and tests after changes.
10. Do not silently change the API contract or package structure.

---

## Immediate Next Task

Implement only the initial application skeleton:

- `internal/config/config.go`
- `internal/app/app.go`
- `cmd/api/main.go`
- `internal/health/controller.go`
- `internal/health/controller_test.go`

Acceptance criteria:

```bash
go fmt ./...
go vet ./...
go test ./...
go run ./cmd/api
curl http://localhost:8080/health
```

Expected response:

```json
{
  "status": "ok"
}
```

Do not implement the calculator parser until the health endpoint works and all tests pass.
