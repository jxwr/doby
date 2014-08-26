package parser

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/jxwr/doby/token"
)

type Lexer struct {
	FileName string
	Src      string
	Pos      int
	Line     int
	Col      int
	LastTok  *DobySymType

	SavedToks []*Tok
	Lines     []string
}

func NewLexer(filename, src string) *Lexer {
	lex := &Lexer{FileName: filename, Src: src + "\n", Pos: 0, Line: 1, Col: 0}
	lex.Lines = strings.Split(lex.Src, "\n")
	return lex
}

var (
	floatRe       = regexp.MustCompile("^[0-9]+\\.[0-9]+")
	intRe         = regexp.MustCompile("^[0-9]+")
	stringRe      = regexp.MustCompile("^\"[^\"]*\"")
	charRe        = regexp.MustCompile("^'.'")
	identRe       = regexp.MustCompile("^[a-zA-Z_][a-zA-Z0-9_]*")
	escapeRe      = regexp.MustCompile("\\\\.")
	lineCommentRe = regexp.MustCompile("^//.*")

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

func (l *Lexer) MkTok(lit string) Tok {
	t := Tok{lit, l.Line, l.Col, token.Pos(l.Pos)}
	l.SavedToks = append(l.SavedToks, &t)
	if len(l.SavedToks) > 16 {
		l.SavedToks = l.SavedToks[len(l.SavedToks)-16:]
	}
	return t
}

func (l *Lexer) Lex(lval *DobySymType) int {
	if l.Pos >= len(l.Src) {
		return 0
	}

	src := l.Src[l.Pos:]
	cur := strings.TrimLeft(src, " \t\r")
	l.Pos += len(src) - len(cur)

	l.LastTok = lval

	if cur[0] == '\n' {
		lval.tok = l.MkTok("\n")
		l.Pos++
		l.Line++
		l.Col = 0
		return EOL
	}

	m := lineCommentRe.FindString(cur)
	if m != "" {
		lval.tok = l.MkTok(m)
		l.Col += len(m)
		l.Pos += len(m)
		return EOL
	}

	m = floatRe.FindString(cur)
	if m != "" {
		lval.tok = l.MkTok(m)
		l.Col += len(m)
		l.Pos += len(m)
		return FLOAT
	}

	m = intRe.FindString(cur)
	if m != "" {
		lval.tok = l.MkTok(m)
		l.Col += len(m)
		l.Pos += len(m)
		return INT
	}

	for _, tok := range OpTokens {
		op := OpTokenMap[tok]

		if strings.HasPrefix(cur, op) {
			lval.tok = l.MkTok(op)
			l.Col += len(op)
			l.Pos += len(op)
			return tok
		}
	}

	for tok, kw := range KeywordTokenMap {
		if strings.HasPrefix(cur, kw) {
			lval.tok = l.MkTok(kw)
			l.Col += len(kw)
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
		lval.tok = l.MkTok(n)
		l.Col += len(m)
		l.Pos += len(m)
		return STRING
	}

	m = charRe.FindString(cur)
	if m != "" {
		lval.tok = l.MkTok(m)
		l.Col += len(m)
		l.Pos += len(m)
		return CHAR
	}

	m = identRe.FindString(cur)
	if m != "" {
		lval.tok = l.MkTok(m)
		l.Col += len(m)
		l.Pos += len(m)
		return IDENT
	}

	// otherwise
	l.Col++
	l.Pos++
	if len(cur) > 0 {
		return int(cur[0])
	}

	return 0
}

func (l *Lexer) Error(s string) {
	fmt.Printf("Syntax Error: Line %d, Col %d:\n", l.Line, l.Col)

	line := l.Line - 5
	if line < 0 {
		line = 0
	}

	for line < l.Line+5 && line < len(l.Lines) {
		if line == l.Line-1 {
			fmt.Printf("*%3d) %s\n", line+1, l.Lines[line])
		} else {
			fmt.Printf(" %3d) %s\n", line+1, l.Lines[line])
		}
		line++
	}
}

func (l *Lexer) PrintPosInfo(pos int) {
	lineNum := 1
	col := 0

	for _, line := range l.Lines {
		if pos < len(line) {
			col = pos + 1
			break
		}
		pos -= len(line) + 1
		lineNum++
	}

	if l.FileName != "" {
		fmt.Printf("\nFile: \"%s\", Line %d, Col %d\n", l.FileName, lineNum, col)
	} else {
		fmt.Printf("\nLine %d, Col %d\n", lineNum, col)
	}

	ln := lineNum - 5
	if ln < 0 {
		ln = 0
	}

	for ln < lineNum+4 && ln < len(l.Lines) {
		if ln == lineNum-1 {
			fmt.Printf("*%3d) %s\n", ln+1, l.Lines[ln])
			fmt.Print("      ")
			for i := 0; i < col-1; i++ {
				if l.Lines[ln][i] == '\t' {
					fmt.Print("\t")
				} else {
					fmt.Print(" ")
				}
			}
			fmt.Println("^")
		} else {
			fmt.Printf(" %3d) %s\n", ln+1, l.Lines[ln])
		}
		ln++
	}
}
