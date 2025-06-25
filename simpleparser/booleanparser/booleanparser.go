package booleanparser

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const _whitespace = " \b\f\n\r\t\v"

var _whitespacemap = map[rune]bool{
	' ':  true,
	'\b': true,
	'\f': true,
	'\n': true,
	'\r': true,
	'\t': true,
	'\v': true,
}

func IsWhiteSpace(r rune) bool {
	_, ok := _whitespacemap[r]
	return ok
}

const _invalidchartemplate = "Invalid or unexpected character in expression: %s"
const _unexpextedendtemplate = "Unexpected end of expression"

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

type ExpressionValue struct {
	Value bool
}

func IsProperLabel(label string) bool {
	ok, err := regexp.MatchString(`^[A-Za-z0-9_][A-Za-z0-9_-]*$`, label)
	if !ok || err != nil {
		return false
	}

	return true
}

type Context struct {
	c map[string]bool
}

// Adds a label to the context
// Returns:
// - TRUE, if label is a proper label
// - FALSE, if the label isn't a proper label
func (ctx *Context) Add(label string) bool {
	if IsProperLabel(label) {
		ctx.c[strings.ToUpper(label)] = true
		return true
	}

	return false
}

func (ctx *Context) Contains(label string) bool {
	return ctx.c[strings.ToUpper(label)]
}

type LabelIdPair struct {
	Label string
	Id    string
}

type Universe struct {
	u6e map[string]string
}

func (u *Universe) Add(p LabelIdPair) bool {
	if IsProperLabel(p.Label) && IsProperLabel(p.Id) {
		l := strings.ToUpper(p.Label)
		i := strings.ToUpper(p.Id)
		u.u6e[l] = l
		u.u6e[i] = l
		return true
	}

	return false
}

func (u *Universe) Contains(label string) bool {
	return u.u6e[strings.ToUpper(label)] != ""
}

func (u *Universe) GetLabel(label string) string {
	return u.u6e[strings.ToUpper(label)]
}

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
		errmsg := _unexpextedendtemplate
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
		errmsg := fmt.Sprintf(_invalidchartemplate, newtoken.Operator)
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

func Primary(ts *TokenStream) (*ExpressionValue, error) {
	t := ts.Get()
	if t == nil {
		return nil, errors.New("Unexpected end-of-stream, Primary not found.")
	}

	switch tokenkind := t.Kind; tokenkind {
	case NOT:
		ev, err := Primary(ts)
		if err != nil {
			return nil, err
		}

		return &ExpressionValue{Value: !ev.Value}, nil

	case OPENPARENTHESES:
		ev, err := Expression(ts)
		if err != nil {
			return nil, err
		}

		t = ts.Get()
		if t == nil {
			return nil, errors.New("Unexpected end-of-stream, closing parentheses not found.")
		}

		if t.Kind != CLOSEPARENTHESES {
			return nil, errors.New("Missing closing parentheses.")
		}

		return &ExpressionValue{Value: ev.Value}, nil

	case LABEL:
		return &ExpressionValue{Value: t.Value}, nil

	default:
		// Do nothing, the error will be thrown in the return after the
		// switch statement
	}

	return nil, errors.New("Primary expected.")
}

func Term(ts *TokenStream) (*ExpressionValue, error) {
	ev, err := Primary(ts)
	if err != nil {
		return nil, err
	}

	left := ev.Value

	for t := ts.Get(); t != nil; {
		switch tokenkind := t.Kind; tokenkind {
		case XOR:
			ev, err = Term(ts)
			if err != nil {
				return nil, err
			}

			right := ev.Value
			left = left != right
		default:
			ts.Push(t)
			return &ExpressionValue{Value: left}, nil
		}

		t = ts.Get()
	}

	return &ExpressionValue{Value: left}, nil
}

func Expression(ts *TokenStream) (*ExpressionValue, error) {
	ev, err := Term(ts)
	if err != nil {
		return nil, err
	}

	left := ev.Value

	for t := ts.Get(); t != nil; {
		switch tokenkind := t.Kind; tokenkind {
		case AND:
			ev, err = Term(ts)
			if err != nil {
				return nil, err
			}

			left = left && ev.Value
		case OR:
			ev, err = Term(ts)
			if err != nil {
				return nil, err
			}

			left = left || ev.Value
		default:
			ts.Push(t)
			return &ExpressionValue{Value: left}, nil
		}

		t = ts.Get()
	}

	return &ExpressionValue{Value: left}, nil
}

func EvaluateBooleanExpression(
	expression string,
	ctx []string,
	unvrs [][]string) (bool, error) {
	universe := &Universe{u6e: make(map[string]string)}

	for _, pair := range unvrs {
		label := pair[0]
		id := pair[1]
		if !universe.Add(LabelIdPair{Label: label, Id: id}) {
			fmt.Printf("[Universe] Failed insertion of '%s': '%s'\n", label, id)
		}
	}

	extendedctx := &Context{c: make(map[string]bool)}

	for _, c := range ctx {
		if !extendedctx.Add(c) {
			fmt.Printf("[Context] Failed insertion of '%s''\n", c)
		}

		if universelabel := universe.GetLabel(c); universelabel != "" {
			extendedctx.Add(universelabel)
		}
	}

	tokens, tokenizeerror := Tokenize(expression, extendedctx, universe)
	if tokenizeerror != nil {
		return false, tokenizeerror
	}

	ts := &TokenStream{tokens: make([]*Token, len(tokens))}

	for _, t := range tokens {
		tkn := t
		ts.Append(&tkn)
	}

	ev, evaluationerror := Expression(ts)

	if evaluationerror != nil {
		return false, evaluationerror
	}

	return ev.Value, nil
}
