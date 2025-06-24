package booleanparser

import (
	"errors"
	"fmt"
	"strings"
)

type Token struct {
	Kind       TokenKind
	Operator   string
	HasValue   bool
	Value      bool
	InUniverse bool
}

// Using a circular list that resizes as needed
type TokenStream struct {
	tokens []*Token
	head   int
	tail   int
	count  int
}

// Add to the end of the list
func (ts *TokenStream) Append(t *Token) {
	if ts.head == ts.tail && ts.count > 0 {
		tokens := make([]*Token, len(ts.tokens)*2)
		copy(tokens, ts.tokens[ts.head:])
		copy(tokens[len(ts.tokens)-ts.head:], ts.tokens[:ts.head])
		ts.head = 0
		ts.tail = len(ts.tokens)
	}

	ts.tokens[ts.tail] = t
	ts.tail = (ts.tail + 1) % len(ts.tokens)
	ts.count++
}

// Add to the beginning of the list, make this element the first of the list.
func (ts *TokenStream) Push(t *Token) {
	if ts.head == ts.tail && ts.count > 0 {
		tokens := make([]*Token, len(ts.tokens)*2)
		copy(tokens, ts.tokens[ts.head:])
		copy(tokens[len(ts.tokens)-ts.head:], ts.tokens[:ts.head])
		ts.head = 0
		ts.tail = len(ts.tokens)
	}

	if ts.head == 0 {
		ts.head = len(ts.tokens) - 1
	} else {
		ts.head--
	}

	ts.tokens[ts.head] = t
	ts.count++
}

// Remove the first element of the list.
func (ts *TokenStream) Get() *Token {
	if ts.count == 0 {
		return nil
	}

	t := ts.tokens[ts.head]
	ts.head = (ts.head + 1) % len(ts.tokens)
	ts.count--
	return t
}

func isvalidruneforlabel(_rune rune) bool {
	if _rune >= 'A' && _rune <= 'Z' {
		return true
	}
	if _rune >= 'a' && _rune <= 'z' {
		return true
	}
	if _rune >= '0' && _rune <= '9' {
		return true
	}
	if _rune == '-' || _rune == '_' {
		return true
	}
	return false
}

func isvalidfirstruneforlabel(_rune rune) bool {
	if _rune >= 'A' && _rune <= 'Z' {
		return true
	}
	if _rune >= 'a' && _rune <= 'z' {
		return true
	}
	if _rune >= '0' && _rune <= '9' {
		return true
	}
	if _rune == '_' {
		return true
	}
	return false
}

func Get_Token(expressionrunes []rune, index int, cp *Context, up *Universe) (Token, int, error) {
	var newtoken Token

	// skip whitespace in between tokens
	for ; index < len(expressionrunes) && IsWhiteSpace(expressionrunes[index]); index++ {
	}

	if index >= len(expressionrunes) {
		newtoken.Kind = INVALID
		newtoken.Operator = ""
		newtoken.HasValue = false
		errmsg := UNEXPECTED_END_OF_TEMPLATE
		return newtoken, index, errors.New(errmsg)
	}

	switch _rune := expressionrunes[index]; _rune {
	case '&':
		newtoken.Kind = AND
		newtoken.Operator = "&"
		newtoken.HasValue = false

		return newtoken, index, nil

	case '|', ',':
		newtoken.Kind = OR
		newtoken.Operator = "|"
		newtoken.HasValue = false

		return newtoken, index, nil

	case '!':
		newtoken.Kind = NOT
		newtoken.Operator = "!"
		newtoken.HasValue = false

		return newtoken, index, nil

	case '^':
		newtoken.Kind = XOR
		newtoken.Operator = "^"
		newtoken.HasValue = false

		return newtoken, index, nil

	case '(':
		newtoken.Kind = OPENPARENTHESES
		newtoken.Operator = "("
		newtoken.HasValue = false

		return newtoken, index, nil

	case ')':
		newtoken.Kind = CLOSEPARENTHESES
		newtoken.Operator = ")"
		newtoken.HasValue = false

		return newtoken, index, nil
	}

	if !isvalidfirstruneforlabel(expressionrunes[index]) {
		newtoken.Kind = INVALID
		newtoken.Operator = string(expressionrunes[index])
		newtoken.HasValue = false
		errmsg := fmt.Sprintf(INVALID_CHAR_ON_TEMPLATE, newtoken.Operator)
		return newtoken, index, errors.New(errmsg)
	}

	label := ""
	for ; index < len(expressionrunes) && isvalidruneforlabel(expressionrunes[index]); index++ {
		label += string(expressionrunes[index])
	}

	ulabel := strings.ToUpper(label)
	newtoken.Kind = LABEL
	newtoken.InUniverse = up.Contains(ulabel)
	newtoken.Operator = ulabel
	if newtoken.InUniverse {
		newtoken.Operator = up.GetLabel(ulabel)
	}
	newtoken.HasValue = true
	newtoken.Value = cp.Contains(ulabel)
	return newtoken, index - 1, nil
}

func Tokenize(expression string, cp *Context, up *Universe) ([]Token, error) {
	var _expression string
	var _expressionrunes []rune
	var _tokens []Token

	expression = strings.Trim(expression, _whitespace)
	if expression == "" {
		return nil, errors.New("Empty expression not allowed")
	}

	_expression = expression
	_expressionrunes = []rune(_expression)

	for runeindex := 0; runeindex < len(_expressionrunes); {
		token, lastindex, get_token_error := Get_Token(_expressionrunes, runeindex, cp, up)
		if get_token_error != nil {
			return nil, errors.Join(errors.New(fmt.Sprintf("Syntax error found near '%s' (position: %d).", string(_expressionrunes[lastindex]), lastindex)), get_token_error)
		}

		_tokens = append(_tokens, token)
		runeindex = lastindex + 1
	}

	return _tokens, nil
}
