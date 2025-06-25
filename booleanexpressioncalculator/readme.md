# Boolean Expression Calculator

## Introduction

There's the need for a Boolean expression calculator, capable of
checking the syntax of the expression and able to evaluate the expression
to obtain its Boolean value.

The expression to be evaluated can be something like any of the following
examples:

- `3461f6ad-6242-430f-954b-72b8380088f0`, a simple label.
- `READWRITE`, another simple label
- `3461f6ad-6242-430f-954b-72b8380088f0 | READWRITE`, a logical OR expression
and the resulting value will depend on the assigned values to the labels.
- `!(ba26bcea-634b-43fe-b38a-e01b98e1c84e | READONLY) &`
`(84feaab0-5a3a-4efa-a185-593ebe40aa04 | 7977e542-a838-417d-bc82-dbf740379522`
`| 3461f6ad-6242-430f-954b-72b8380088f0) ^ (READWRITE | ADMIN | BACKUP)`, a
rather complex expression.
- `84feaab0-5a3a-4efa-a185-593ebe40aa04, 7977e542-a838-417d-bc82-dbf740379522`
`, 3461f6ad-6242-430f-954b-72b8380088f0`, an OR expression using commas for
compatibility with list expressions.

In the expression, labels are used to identify constants associated to some
data set that we can name as the context in which the expression is to be
evaluated; if a label in the expression is part of the context the
label is associated with a `true` value, otherwise the label is associated
with a `false` value.

The desired operations for the calculator are:

Operation | Name | Symbol
--|--|--
Conjunction | _and_ | `&`
Disjunction | _or_ | `\|`, `,`
Unary Negation | _not_ | `!`
Exclusive OR | _xor_ | `^`
Group | _parentheses_ | `(` ... `)`

The precedence of the operations:

Precedence | Operation
--|--
1 | Group
2 | Negation
3 | Exclusive OR
4 | Conjunction, Disjunction (left to right evaluation)

To implement this calculator, a similar parser is found in chapter 6 of the
book "Programming, Principles and Practice Using C++", by Bjarne Stroustrup,
ISBN 978-0-321-54372; the book shows how to implement an arithmetic
calculator with a similar grammar; here will take the knowledge given in the
book and apply it to the Boolean calculator.

## Boolean Calculator Grammar

```regularGrammar
Expression:
    Term
    Expression "&" Term
    Expression "|" Term
    Expression "," Term

Term:
    Primary
    Term "^" Primary

Primary:
    Label
    "!" Primary
    "(" Expression ")"

Label:
    "[A-Z0-9_-]"
    "[A-Z0-9_-]" Label
```

## Funcionality

There are two main requirements and one desired requirement of the calculator:

1. [Must Have] The ability to validate the syntax of the expression and to
emit actionable messages when the syntax is invalid.
1. [Must Have] The ability to evaluate the expression, under a given context.
1. [Desired] The ability to evaluate one expression under several contexts and
compare each result against an expexted outcome (aka test validation of the
expression).
