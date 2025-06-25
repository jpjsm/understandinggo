package booleanparser

import "errors"

type ExpressionValue struct {
	Value bool
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
