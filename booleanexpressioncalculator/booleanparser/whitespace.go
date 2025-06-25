package booleanparser

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
