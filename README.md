<p align="center">
  <img src="assets/logo.png" alt="Cloudws E2e RPN Go logo" width="180">
</p>

# cloudws-e2e-rpn-go

A small REST API service implementing a scientific calculator that evaluates
[Reverse Polish Notation (RPN)](https://en.wikipedia.org/wiki/Reverse_Polish_notation)
expressions. Built with Go 1.23 and the `net/http` standard library — no
external dependencies.

## Why RPN?

Reverse Polish Notation places operators after their operands (e.g. `3 4 +`
instead of `3 + 4`), which removes the need for parentheses or operator
precedence rules. Expressions are evaluated left to right using a stack:
numbers are pushed, and operators pop their operands, compute a result, and
push it back.

## Supported operators

| Token  | Meaning                          | Arity  |
|--------|-----------------------------------|--------|
| `+`    | addition                          | binary |
| `-`    | subtraction                       | binary |
| `*`    | multiplication                    | binary |
| `/`    | division                          | binary |
| `^`    | power (`a ^ b` = a raised to b)   | binary |
| `sqrt` | square root                       | unary  |
| `sin`  | sine (radians)                    | unary  |
| `cos`  | cosine (radians)                  | unary  |
| `tan`  | tangent (radians)                 | unary  |
| `log`  | natural logarithm                 | unary  |

For binary operators, the expression `a b OP` computes `a OP b` (operand
order is preserved, so `5 2 -` is `5 - 2 = 3`).

## Error handling

The service returns HTTP 400 with a JSON error body for:
- **Stack underflow** — an operator doesn't have enough operands available.
- **Division by zero**.
- **Unknown operators** — any token that isn't a number or a supported
  operator.
- **Non-numeric tokens** — tokens that fail to parse as numbers or operators.
- **Malformed expressions** — e.g. leftover values on the stack after
  evaluation, or an empty expression.
- **Domain errors** — e.g. `sqrt` of a negative number, `log` of a
  non-positive number.
- **Invalid JSON** request bodies.

## Build

Requires Go 1.23+.

```bash
go build ./...
```

## Run

```bash
go run ./cmd/server
```

The server listens on port `8080` by default. Set the `PORT` environment
variable to override:

```bash
PORT=9090 go run ./cmd/server
```

## API

### `POST /calculate`

Request body:

```json
{ "expression": "<space-separated RPN tokens>" }
```

Success response (`200 OK`):

```json
{ "result": <number> }
```

Error response (`400 Bad Request`):

```json
{ "error": "<message>" }
```

## Example requests

```bash
# Basic arithmetic: (3 + 4) * 2 = 14
curl -X POST http://localhost:8080/calculate \
  -H 'Content-Type: application/json' \
  -d '{"expression": "3 4 + 2 *"}'
# {"result":14}

# Power and sqrt: 2^10, then sqrt(16)
curl -X POST http://localhost:8080/calculate \
  -H 'Content-Type: application/json' \
  -d '{"expression": "16 sqrt"}'
# {"result":4}

# Trig functions (radians)
curl -X POST http://localhost:8080/calculate \
  -H 'Content-Type: application/json' \
  -d '{"expression": "0 sin"}'
# {"result":0}

# Natural log
curl -X POST http://localhost:8080/calculate \
  -H 'Content-Type: application/json' \
  -d '{"expression": "1 log"}'
# {"result":0}

# Error: division by zero
curl -X POST http://localhost:8080/calculate \
  -H 'Content-Type: application/json' \
  -d '{"expression": "5 0 /"}'
# {"error":"division by zero"}

# Error: stack underflow
curl -X POST http://localhost:8080/calculate \
  -H 'Content-Type: application/json' \
  -d '{"expression": "1 +"}'
# {"error":"stack underflow: operator \"+\" requires 2 operands, got 1"}
```

## Testing

```bash
go test ./... -v
```

## Project layout

```
cmd/server        entry point (main.go)
internal/api       HTTP handler for /calculate
internal/rpn       RPN expression parser/evaluator (core logic)
```
