package booleanparser

import (
	"regexp"
	"strings"
)

func IsProperLabel(label string) bool {
	ok, err := regexp.MatchString(`^[A-Za-z0-9_][A-Za-z0-9_-]*$`, label)
	if !ok || err != nil {
		return false
	}

	return true
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
