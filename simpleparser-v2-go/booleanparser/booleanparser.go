package booleanparser

import (
	"fmt"
)

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
