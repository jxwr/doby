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
	intRe    = regexp.MustCompile("^[-0-9]+")
	stringRe = regexp.MustCompile("^\"[^\"]*\"")
	charRe   = regexp.MustCompile("^'.*'")
	identRe  = regexp.MustCompile("^[a-zA-Z_][a-zA-Z0-9_]*")
	escapeRe = regexp.MustCompile("\\\\.")

	SpecTokens = map[int]string{
		EOF:     "EOF",
		EOL:     "EOL",
		COMMENT: "COMMENT",
	}

	AtomTokenMap = map[int]string{
		IDENT:  "IDENT",
		INT:    "INT",
		FLOAT:  "FLOAT",
		CHAR:   "CHAR",
		STRING: "STRING",
	}

	OpTokens = [...]int{
		ADD_ASSIGN, // "+=",
		SUB_ASSIGN, // "-=",
		MUL_ASSIGN, // "*=",
		QUO_ASSIGN, // "/=",
		REM_ASSIGN, // "%=",

		AND_ASSIGN,     // "&=",
		OR_ASSIGN,      // "|=",
		XOR_ASSIGN,     // "^=",
		SHL_ASSIGN,     // "<<=",
		SHR_ASSIGN,     // ">>=",
		AND_NOT_ASSIGN, // "&^=",

		SHL,     // "<<",
		SHR,     // ">>",
		AND_NOT, // "&^",

		LAND,  // "&&",
		LOR,   // "||",
		ARROW, // "<-",
		INC,   // "++",
		DEC,   // "--",
		EQL,   // "==",

		NEQ,      // "!=",
		LEQ,      // "<=",
		GEQ,      // ">=",
		DEFINE,   // ":=",
		ELLIPSIS, // "...",

		ADD, // "+",
		SUB, // "-",
		MUL, // "*",
		QUO, // "/",
		REM, // "%",

		AND, // "&",
		OR,  // "|",
		XOR, // "^",

		LSS,    // "<",
		GTR,    // ">",
		ASSIGN, // "=",
		NOT,    // "!",

		LPAREN, // "(",
		LBRACK, // "[",
		LBRACE, // "{",
		COMMA,  // ",",
		PERIOD, // ".",

		RPAREN,    // ")",
		RBRACK,    // "]",
		RBRACE,    // "}",
		SEMICOLON, // ";",
		COLON,     // ":",
	}

	OpTokenMap = map[int]string{
		ADD_ASSIGN: "+=",
		SUB_ASSIGN: "-=",
		MUL_ASSIGN: "*=",
		QUO_ASSIGN: "/=",
		REM_ASSIGN: "%=",

		AND_ASSIGN:     "&=",
		OR_ASSIGN:      "|=",
		XOR_ASSIGN:     "^=",
		SHL_ASSIGN:     "<<=",
		SHR_ASSIGN:     ">>=",
		AND_NOT_ASSIGN: "&^=",

		SHL:     "<<",
		SHR:     ">>",
		AND_NOT: "&^",

		LAND:  "&&",
		LOR:   "||",
		ARROW: "<-",
		INC:   "++",
		DEC:   "--",
		EQL:   "==",

		NEQ:      "!=",
		LEQ:      "<=",
		GEQ:      ">=",
		DEFINE:   ":=",
		ELLIPSIS: "...",

		ADD: "+",
		SUB: "-",
		MUL: "*",
		QUO: "/",
		REM: "%",

		AND: "&",
		OR:  "|",
		XOR: "^",

		LSS:    "<",
		GTR:    ">",
		ASSIGN: "=",
		NOT:    "!",

		LPAREN: "(",
		LBRACK: "[",
		LBRACE: "{",
		COMMA:  ",",
		PERIOD: ".",

		RPAREN:    ")",
		RBRACK:    "]",
		RBRACE:    "}",
		SEMICOLON: ";",
		COLON:     ":",
	}

	KeywordTokenMap = map[int]string{
		BREAK:    "break",
		CASE:     "case",
		CHAN:     "chan",
		CONST:    "const",
		CONTINUE: "continue",

		DEFAULT:     "default",
		DEFER:       "defer",
		ELSE:        "else",
		FALLTHROUGH: "fallthrough",
		FOR:         "for",

		FUNC:   "func",
		GO:     "go",
		GOTO:   "goto",
		IF:     "if",
		IMPORT: "import",

		INTERFACE: "interface",
		MAP:       "map",
		PACKAGE:   "package",
		RANGE:     "range",
		RETURN:    "return",

		SELECT: "select",
		STRUCT: "struct",
		SWITCH: "switch",
		TYPE:   "type",
		VAR:    "var",
	}
)

func (l *Lexer) Lex(lval *DoubiSymType) int {
	if l.Pos >= len(l.Src) {
		return 0
	}

	src := l.Src[l.Pos:]
	cur := strings.TrimLeft(src, " \t\r")
	l.Pos += len(src) - len(cur)

	if cur[0] == '\n' {
		lval.lit = "\n"
		l.Pos++
		return EOL
	}

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

	for _, tok := range OpTokens {
		op := OpTokenMap[tok]

		if strings.HasPrefix(cur, op) {
			lval.lit = op
			l.Pos += len(op)
			return tok
		}
	}

	for tok, kw := range KeywordTokenMap {
		if strings.HasPrefix(cur, kw) {
			lval.lit = kw
			l.Pos += len(kw)
			return tok
		}
	}

	m = stringRe.FindString(cur)
	if m != "" {
		n := escapeRe.ReplaceAllStringFunc(m, func(s string) string {
			esc := strings.TrimPrefix(s, "\\")
			switch esc {
			case "n":
				return "\n"
			case "t":
				return "\t"
			case "r":
				return "\r"
			case "\"":
				return "\""
			default:
				return esc
			}
		})
		lval.lit = n

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

	// otherwise
	l.Pos++
	if len(cur) > 0 {
		return int(cur[0])
	}

	return 0
}

func (l *Lexer) Error(s string) {
	fmt.Printf("syntax error: %s\n", s)
}
