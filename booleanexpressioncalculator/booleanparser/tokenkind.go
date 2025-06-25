package booleanparser

import "fmt"

type TokenKind int

const (
	LABEL TokenKind = iota + 1
	OPENPARENTHESES
	CLOSEPARENTHESES
	NOT
	XOR
	AND
	OR
	INVALID
)

func GetTokenKindName(t TokenKind) string {
	names := []string{
		"Undefined",
		"LABEL",
		"OPENPARENTHESES",
		"CLOSEPARENTHESES",
		"NOT",
		"XOR",
		"AND",
		"OR",
		"INVALID",
	}

	if t >= LABEL && t <= INVALID {
		return names[t]
	}

	return fmt.Sprintf("Undefined: '%d'", t)
}
