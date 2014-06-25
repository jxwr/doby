package parser

import (
	"fmt"
	"regexp"
	"strings"
)

type Lexer struct {
	Src string
	Pos int
}

var (
	floatRe  = regexp.MustCompile("^[0-9]+\\.[0-9]+")
	intRe    = regexp.MustCompile("^[0-9]+")
	stringRe = regexp.MustCompile("^\".*\"")
	charRe   = regexp.MustCompile("^'.*'")
	identRe  = regexp.MustCompile("^[a-zA-Z_][a-zA-Z0-9_]*")
)

func (l *Lexer) Lex(lval *CalcSymType) int {
	cur := strings.TrimSpace(l.Src[l.Pos:])

	m := floatRe.FindString(cur)
	if m != "" {
		lval.lit = m
		l.Pos += len(m)
		return FLOAT
	}

	m = intRe.FindString(cur)
	if m != "" {
		lval.lit = m
		l.Pos += len(m)
		return INT
	}

	m = stringRe.FindString(cur)
	if m != "" {
		lval.lit = m
		l.Pos += len(m)
		return STRING
	}

	m = charRe.FindString(cur)
	if m != "" {
		lval.lit = m
		l.Pos += len(m)
		return CHAR
	}

	m = identRe.FindString(cur)
	if m != "" {
		lval.lit = m
		l.Pos += len(m)
		return IDENT
	}

	switch {
	case strings.HasPrefix(cur, "+"):
		l.Pos++
		return '+'
	case strings.HasPrefix(cur, "-"):
		l.Pos++
		return '-'
	case strings.HasPrefix(cur, "/"):
		l.Pos++
		return '/'
	case strings.HasPrefix(cur, "%"):
		l.Pos++
		return '%'
	case strings.HasPrefix(cur, "&"):
		l.Pos++
		return '&'
	case strings.HasPrefix(cur, "|"):
		l.Pos++
		return '|'
	}

	return 0
}

func (l *Lexer) Error(s string) {
	fmt.Printf("syntax error: %s\n", s)
}
