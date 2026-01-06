# Typechecker REPL

A small Go REPL that feeds input into a typechecker for a Lisp-like expression language.

## Features

- Scalars: numbers, booleans, nil/null, single- and double-quoted strings
- Binary ops: `+ - * / %`
- Multiple input styles: prefix, parenthesized, or JSON arrays

## Examples

```
1
'hello'
"hello"
true
nil
+ 1 2
(+ 1 2)
["+", 1, 2]
```

## Run

```
go run .
```

## Tests

```
go test ./...
```
