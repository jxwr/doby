
//line grammar.y:2

package parser
import __yyfmt__ "fmt"
//line grammar.y:3
		
import (
       "github.com/jxwr/doby/ast"
       "github.com/jxwr/doby/token"
)

var ProgramAst []ast.Stmt

type Tok struct {
    Lit string
    Line int
    Col int
    Pos token.Pos
}

func (t Tok) String() string {
    return t.Lit
}


//line grammar.y:27
type DobySymType struct {
	yys int
    node ast.Node
    expr ast.Expr
    expr_list []ast.Expr
    stmt ast.Stmt
    stmt_list []ast.Stmt
    field *ast.Field
    field_list []*ast.Field
    ident_list []*ast.Ident
    tok Tok
}

const EOF = 57346
const EOL = 57347
const COMMENT = 57348
const IDENT = 57349
const INT = 57350
const FLOAT = 57351
const STRING = 57352
const CHAR = 57353
const SHL = 57354
const SHR = 57355
const AND_NOT = 57356
const ADD_ASSIGN = 57357
const SUB_ASSIGN = 57358
const MUL_ASSIGN = 57359
const QUO_ASSIGN = 57360
const REM_ASSIGN = 57361
const AND_ASSIGN = 57362
const OR_ASSIGN = 57363
const XOR_ASSIGN = 57364
const SHL_ASSIGN = 57365
const SHR_ASSIGN = 57366
const AND_NOT_ASSIGN = 57367
const LAND = 57368
const LOR = 57369
const ARROW = 57370
const INC = 57371
const DEC = 57372
const EQL = 57373
const NEQ = 57374
const LEQ = 57375
const GEQ = 57376
const DEFINE = 57377
const ELLIPSIS = 57378
const ADD = 57379
const SUB = 57380
const MUL = 57381
const QUO = 57382
const REM = 57383
const AND = 57384
const OR = 57385
const XOR = 57386
const LSS = 57387
const GTR = 57388
const ASSIGN = 57389
const NOT = 57390
const LPAREN = 57391
const LBRACK = 57392
const LBRACE = 57393
const COMMA = 57394
const PERIOD = 57395
const RPAREN = 57396
const RBRACK = 57397
const RBRACE = 57398
const SEMICOLON = 57399
const COLON = 57400
const BREAK = 57401
const CASE = 57402
const CHAN = 57403
const CONTINUE = 57404
const CONST = 57405
const DEFAULT = 57406
const DEFER = 57407
const ELSE = 57408
const FALLTHROUGH = 57409
const FOR = 57410
const FUNC = 57411
const GO = 57412
const GOTO = 57413
const IF = 57414
const IMPORT = 57415
const INTERFACE = 57416
const MAP = 57417
const PACKAGE = 57418
const RANGE = 57419
const RETURN = 57420
const SELECT = 57421
const STRUCT = 57422
const SWITCH = 57423
const TYPE = 57424
const VAR = 57425
const UMINUS = 57426

var DobyToknames = []string{
	"EOF",
	"EOL",
	"COMMENT",
	"IDENT",
	"INT",
	"FLOAT",
	"STRING",
	"CHAR",
	"SHL",
	"SHR",
	"AND_NOT",
	"ADD_ASSIGN",
	"SUB_ASSIGN",
	"MUL_ASSIGN",
	"QUO_ASSIGN",
	"REM_ASSIGN",
	"AND_ASSIGN",
	"OR_ASSIGN",
	"XOR_ASSIGN",
	"SHL_ASSIGN",
	"SHR_ASSIGN",
	"AND_NOT_ASSIGN",
	"LAND",
	"LOR",
	"ARROW",
	"INC",
	"DEC",
	"EQL",
	"NEQ",
	"LEQ",
	"GEQ",
	"DEFINE",
	"ELLIPSIS",
	"ADD",
	"SUB",
	"MUL",
	"QUO",
	"REM",
	"AND",
	"OR",
	"XOR",
	"LSS",
	"GTR",
	"ASSIGN",
	"NOT",
	"LPAREN",
	"LBRACK",
	"LBRACE",
	"COMMA",
	"PERIOD",
	"RPAREN",
	"RBRACK",
	"RBRACE",
	"SEMICOLON",
	"COLON",
	"BREAK",
	"CASE",
	"CHAN",
	"CONTINUE",
	"CONST",
	"DEFAULT",
	"DEFER",
	"ELSE",
	"FALLTHROUGH",
	"FOR",
	"FUNC",
	"GO",
	"GOTO",
	"IF",
	"IMPORT",
	"INTERFACE",
	"MAP",
	"PACKAGE",
	"RANGE",
	"RETURN",
	"SELECT",
	"STRUCT",
	"SWITCH",
	"TYPE",
	"VAR",
	"UMINUS",
}
var DobyStatenames = []string{}

const DobyEofCode = 1
const DobyErrCode = 2
const DobyMaxDepth = 200

//line yacctab:1
var DobyExca = []int{
	-1, 0,
	5, 122,
	57, 122,
	-2, 12,
	-1, 1,
	1, -1,
	-2, 0,
	-1, 18,
	5, 72,
	51, 72,
	56, 72,
	57, 72,
	60, 72,
	64, 72,
	-2, 13,
	-1, 24,
	5, 122,
	56, 122,
	57, 122,
	-2, 12,
	-1, 54,
	1, 127,
	5, 126,
	57, 126,
	-2, 12,
	-1, 94,
	5, 88,
	51, 88,
	56, 88,
	57, 88,
	60, 88,
	64, 88,
	-2, 65,
	-1, 105,
	57, 72,
	-2, 13,
	-1, 158,
	5, 126,
	56, 126,
	57, 126,
	60, 126,
	64, 126,
	-2, 12,
	-1, 190,
	5, 122,
	56, 122,
	57, 122,
	60, 122,
	64, 122,
	-2, 12,
	-1, 211,
	5, 122,
	56, 122,
	57, 122,
	60, 122,
	64, 122,
	-2, 12,
}

const DobyNprod = 128
const DobyPrivate = 57344

var DobyTokenNames []string
var DobyStates []string

const DobyLast = 1131

var DobyAct = []int{

	97, 18, 11, 178, 176, 101, 188, 2, 163, 186,
	164, 81, 190, 166, 165, 81, 81, 211, 196, 171,
	202, 95, 239, 225, 30, 18, 99, 18, 158, 105,
	232, 54, 98, 202, 216, 226, 202, 81, 203, 58,
	57, 113, 114, 56, 162, 81, 24, 184, 102, 108,
	109, 110, 237, 198, 180, 18, 18, 3, 120, 158,
	123, 124, 125, 126, 127, 128, 129, 130, 131, 132,
	133, 134, 135, 136, 137, 138, 139, 140, 141, 142,
	55, 119, 143, 55, 100, 221, 103, 43, 44, 45,
	46, 47, 217, 43, 44, 45, 46, 47, 116, 164,
	199, 206, 159, 165, 200, 167, 160, 107, 168, 234,
	157, 55, 117, 118, 194, 177, 36, 223, 49, 204,
	1, 179, 183, 43, 49, 161, 17, 16, 50, 48,
	51, 15, 14, 13, 50, 48, 51, 94, 12, 81,
	115, 208, 218, 10, 9, 185, 61, 62, 63, 53,
	59, 60, 61, 62, 63, 53, 58, 57, 8, 18,
	56, 81, 58, 57, 195, 52, 56, 191, 19, 7,
	187, 52, 6, 59, 60, 61, 62, 63, 64, 65,
	66, 5, 207, 4, 205, 58, 57, 18, 175, 56,
	96, 18, 41, 18, 215, 40, 39, 106, 212, 177,
	177, 38, 222, 219, 220, 37, 224, 42, 35, 34,
	33, 32, 18, 31, 18, 0, 117, 230, 231, 228,
	111, 0, 177, 0, 0, 0, 233, 122, 0, 235,
	0, 0, 236, 43, 44, 45, 46, 47, 0, 0,
	0, 238, 240, 0, 210, 0, 0, 0, 0, 0,
	214, 145, 146, 147, 148, 149, 150, 151, 152, 153,
	154, 155, 156, 0, 49, 0, 0, 0, 0, 0,
	0, 229, 0, 0, 50, 48, 51, 24, 0, 0,
	0, 172, 173, 104, 0, 22, 0, 0, 23, 43,
	44, 45, 46, 47, 28, 53, 20, 0, 25, 29,
	0, 0, 0, 0, 21, 27, 0, 26, 0, 0,
	0, 52, 0, 43, 44, 45, 46, 47, 0, 174,
	49, 43, 44, 45, 46, 47, 0, 0, 0, 0,
	50, 48, 51, 189, 0, 0, 0, 0, 145, 0,
	0, 0, 0, 197, 49, 0, 0, 0, 0, 0,
	0, 53, 49, 0, 50, 48, 51, 24, 0, 193,
	0, 0, 50, 48, 51, 22, 0, 52, 23, 43,
	44, 45, 46, 47, 28, 53, 20, 0, 25, 29,
	0, 0, 0, 53, 21, 27, 0, 26, 0, 0,
	144, 52, 43, 44, 45, 46, 47, 0, 0, 52,
	49, 43, 44, 45, 46, 47, 0, 0, 0, 0,
	50, 48, 51, 112, 0, 43, 44, 45, 46, 47,
	121, 0, 0, 49, 0, 0, 0, 0, 0, 0,
	0, 53, 49, 50, 48, 51, 0, 0, 0, 0,
	0, 0, 50, 48, 51, 0, 49, 52, 59, 60,
	61, 62, 63, 64, 53, 66, 50, 48, 51, 0,
	58, 57, 0, 53, 56, 0, 0, 0, 0, 0,
	52, 0, 0, 0, 0, 0, 0, 53, 0, 52,
	67, 68, 69, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 52, 76, 77, 0, 0, 0, 75,
	72, 73, 74, 0, 0, 59, 60, 61, 62, 63,
	64, 65, 66, 70, 71, 0, 0, 58, 57, 0,
	0, 56, 0, 182, 0, 0, 181, 67, 68, 69,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 76, 77, 0, 0, 0, 75, 72, 73, 74,
	0, 0, 59, 60, 61, 62, 63, 64, 65, 66,
	70, 71, 0, 0, 58, 57, 0, 0, 56, 67,
	68, 69, 0, 201, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 76, 77, 0, 0, 0, 75, 72,
	73, 74, 0, 0, 59, 60, 61, 62, 63, 64,
	65, 66, 70, 71, 0, 0, 58, 57, 0, 0,
	56, 67, 68, 69, 213, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 76, 77, 0, 0, 0,
	75, 72, 73, 74, 0, 0, 59, 60, 61, 62,
	63, 64, 65, 66, 70, 71, 0, 0, 58, 57,
	0, 0, 56, 67, 68, 69, 192, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 76, 77, 78,
	79, 80, 75, 72, 73, 74, 0, 0, 59, 60,
	61, 62, 63, 64, 65, 66, 70, 71, 0, 0,
	58, 57, 24, 0, 56, 67, 68, 69, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 76,
	77, 0, 0, 0, 75, 72, 73, 74, 0, 0,
	59, 60, 61, 62, 63, 64, 65, 66, 70, 71,
	0, 0, 58, 57, 0, 0, 56, 0, 227, 67,
	68, 69, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 76, 77, 0, 0, 0, 75, 72,
	73, 74, 0, 0, 59, 60, 61, 62, 63, 64,
	65, 66, 70, 71, 0, 0, 58, 57, 0, 0,
	56, 0, 209, 67, 68, 69, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 76, 77, 78,
	79, 80, 75, 72, 73, 74, 0, 0, 59, 60,
	61, 62, 63, 64, 65, 66, 70, 71, 0, 0,
	58, 57, 0, 0, 56, 67, 68, 69, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 76,
	77, 0, 0, 0, 75, 72, 73, 74, 0, 0,
	59, 60, 61, 62, 63, 64, 65, 66, 70, 71,
	0, 0, 58, 57, 0, 0, 56, 170, 67, 68,
	69, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 76, 77, 0, 0, 0, 75, 72, 73,
	74, 0, 0, 59, 60, 61, 62, 63, 64, 65,
	66, 70, 71, 0, 0, 58, 57, 24, 0, 56,
	67, 68, 69, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 76, 77, 0, 0, 0, 75,
	72, 73, 74, 0, 0, 59, 60, 61, 62, 63,
	64, 65, 66, 70, 71, 0, 0, 58, 57, 0,
	0, 56, 67, 68, 69, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 76, 0, 0, 0,
	0, 75, 72, 73, 74, 0, 0, 59, 60, 61,
	62, 63, 64, 65, 66, 70, 71, 0, 0, 58,
	57, 0, 0, 56, 67, 68, 69, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 75, 72, 73, 74, 0, 0, 59,
	60, 61, 62, 63, 64, 65, 66, 70, 71, 0,
	0, 58, 57, 0, 0, 56, 75, 72, 73, 74,
	0, 0, 59, 60, 61, 62, 63, 64, 65, 66,
	70, 71, 0, 0, 58, 57, 0, 0, 56, 75,
	72, 73, 74, 0, 0, 59, 60, 61, 62, 63,
	64, 65, 66, 0, 0, 0, 0, 58, 57, 0,
	0, 56, 83, 84, 85, 86, 87, 88, 89, 90,
	91, 92, 93, 83, 84, 85, 86, 87, 88, 89,
	90, 91, 92, 93, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 169, 0, 0, 0, 0, 81,
	0, 0, 0, 0, 0, 82, 0, 0, 0, 0,
	81,
}
var DobyPact = []int{

	306, -1000, 26, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, 771, 1078,
	394, 394, -1000, -1000, 306, 394, 306, -3, 226, 97,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, 394, 394,
	394, 408, -9, 91, 306, 306, 116, 362, 394, 394,
	394, 394, 394, 394, 394, 394, 394, 394, 394, 394,
	394, 394, 394, 394, 394, 394, 394, 394, 394, -1000,
	-1000, 385, 394, 394, 394, 394, 394, 394, 394, 394,
	394, 394, 394, 394, -1000, 898, -15, 898, 54, 856,
	-3, -1000, 39, -44, 394, 641, 1067, -1000, 813, -10,
	982, -36, 394, 314, 394, 114, 5, -1000, -1000, -1000,
	468, 394, -7, 107, 107, -10, -10, -10, 113, 411,
	113, 1005, 1005, 1005, 1028, 1028, 136, 136, 136, 136,
	982, 940, 898, 898, 394, -15, -15, -15, -15, -15,
	-15, -15, -15, -15, -15, -15, -15, -1000, 306, -57,
	-1000, -50, -1000, -1000, 394, -46, 394, 599, -1000, 282,
	-1000, -1000, 109, -37, 394, 48, -1000, 515, -16, 112,
	94, 86, -1000, 727, -1000, 898, 306, -1000, -1000, -41,
	306, 557, 306, 394, -21, -1000, -1000, 87, 394, 80,
	-1000, 394, 110, -5, -31, -19, -1000, 683, -1000, -1000,
	-1000, 306, 23, 306, -5, 856, -1000, -25, -1000, -1000,
	-1000, 394, 898, -1000, -1000, 102, -5, -1000, 23, -5,
	-1000, -1000, -1000, -1000, 3, -1000, -1000, 94, -32, -5,
	-1000,
}
var DobyPgo = []int{

	0, 0, 24, 213, 211, 210, 209, 208, 207, 116,
	205, 201, 196, 195, 192, 168, 4, 188, 3, 57,
	183, 181, 172, 169, 158, 144, 143, 2, 138, 8,
	5, 133, 132, 131, 127, 126, 7, 125, 120,
}
var DobyR1 = []int{

	0, 2, 3, 3, 3, 3, 4, 5, 7, 7,
	7, 6, 15, 15, 15, 15, 9, 10, 10, 11,
	11, 11, 11, 11, 11, 11, 11, 11, 11, 11,
	11, 11, 11, 11, 11, 11, 11, 11, 12, 12,
	12, 14, 14, 14, 16, 17, 17, 17, 17, 17,
	17, 17, 13, 18, 18, 18, 8, 8, 8, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 20, 21, 22, 22, 23, 23, 23, 23,
	23, 23, 23, 23, 23, 23, 23, 23, 24, 25,
	26, 26, 27, 28, 28, 29, 29, 37, 37, 37,
	30, 31, 32, 33, 33, 33, 34, 35, 19, 19,
	19, 19, 19, 19, 19, 19, 19, 19, 19, 19,
	19, 19, 36, 36, 36, 36, 36, 38,
}
var DobyR2 = []int{

	0, 1, 1, 1, 1, 1, 3, 3, 6, 5,
	5, 4, 0, 1, 3, 4, 4, 2, 2, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 5,
	4, 4, 6, 5, 3, 0, 1, 3, 3, 4,
	2, 3, 4, 0, 1, 3, 5, 6, 10, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 3, 2, 2, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 2, 2,
	1, 1, 3, 3, 5, 4, 3, 1, 1, 2,
	3, 3, 2, 7, 6, 3, 6, 2, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 0, 1, 3, 3, 2, 2,
}
var DobyChk = []int{

	-1000, -38, -36, -19, -20, -21, -22, -23, -24, -25,
	-26, -27, -28, -31, -32, -33, -34, -35, -1, -15,
	70, 78, 59, 62, 51, 72, 81, 79, 68, 73,
	-2, -3, -4, -5, -6, -7, -9, -10, -11, -12,
	-13, -14, -8, 7, 8, 9, 10, 11, 49, 38,
	48, 50, 85, 69, 5, 57, 53, 50, 49, 37,
	38, 39, 40, 41, 42, 43, 44, 12, 13, 14,
	45, 46, 32, 33, 34, 31, 26, 27, 28, 29,
	30, 52, 47, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, -9, -1, -15, -1, -36, -1,
	-19, -30, 51, -19, 57, -1, -15, 10, -1, -1,
	-1, -15, 5, 50, 51, 49, 7, -19, -19, -2,
	-1, 58, -15, -1, -1, -1, -1, -1, -1, -1,
	-1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
	-1, -1, -1, -1, 5, -15, -15, -15, -15, -15,
	-15, -15, -15, -15, -15, -15, -15, 56, 5, -27,
	-30, -37, 5, -29, 60, 64, 57, -1, -27, 47,
	54, 55, -15, -15, 5, -17, -16, -1, -18, 7,
	49, 58, 55, -1, 54, -1, 66, -29, 56, -15,
	58, -1, 57, 77, 5, 55, 55, -15, 5, 52,
	56, 58, 52, 54, 7, -18, 7, -1, 55, 55,
	-19, 58, -36, 57, -19, -1, 55, 5, 55, -16,
	-16, 5, -1, 7, -27, 54, 54, 55, -36, -19,
	-27, -27, 55, -16, 7, -27, -27, 49, -18, 54,
	-27,
}
var DobyDef = []int{

	-2, -2, 0, 123, 108, 109, 110, 111, 112, 113,
	114, 115, 116, 117, 118, 119, 120, 121, -2, 0,
	0, 12, 90, 91, -2, 0, 12, 0, 12, 0,
	59, 60, 61, 62, 63, 64, 65, 66, 67, 68,
	69, 70, 71, 1, 2, 3, 4, 5, 0, 0,
	0, 12, 0, 0, -2, 12, 0, 0, 12, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 74,
	75, 0, 12, 12, 12, 12, 12, 12, 12, 12,
	12, 12, 12, 12, -2, 0, 89, 13, 0, 0,
	0, 102, 0, 0, 0, -2, 0, 107, 0, 17,
	18, 0, 12, 12, 45, 53, 0, 124, 125, 7,
	0, 0, 0, 19, 20, 21, 22, 23, 24, 25,
	26, 27, 28, 29, 30, 31, 32, 33, 34, 35,
	36, 37, 73, 14, 0, 76, 77, 78, 79, 80,
	81, 82, 83, 84, 85, 86, 87, 92, -2, 93,
	101, 0, 97, 98, 12, 0, 0, 0, 105, 12,
	6, 38, 0, 0, 12, 0, 46, 0, 0, 54,
	53, 0, 11, 0, 16, 15, 12, 99, 100, 0,
	-2, 0, 12, 0, 0, 40, 41, 0, 50, 0,
	52, 0, 0, 0, 0, 0, 54, 0, 10, 9,
	94, -2, 96, 12, 0, 0, 39, 0, 43, 47,
	48, 51, 44, 55, 56, 0, 0, 8, 95, 0,
	104, 106, 42, 49, 0, 57, 103, 53, 0, 0,
	58,
}
var DobyTok1 = []int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 85,
}
var DobyTok2 = []int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, 34, 35, 36, 37, 38, 39, 40, 41,
	42, 43, 44, 45, 46, 47, 48, 49, 50, 51,
	52, 53, 54, 55, 56, 57, 58, 59, 60, 61,
	62, 63, 64, 65, 66, 67, 68, 69, 70, 71,
	72, 73, 74, 75, 76, 77, 78, 79, 80, 81,
	82, 83, 84,
}
var DobyTok3 = []int{
	0,
}

//line yaccpar:1

/*	parser for yacc output	*/

var DobyDebug = 0

type DobyLexer interface {
	Lex(lval *DobySymType) int
	Error(s string)
}

const DobyFlag = -1000

func DobyTokname(c int) string {
	// 4 is TOKSTART above
	if c >= 4 && c-4 < len(DobyToknames) {
		if DobyToknames[c-4] != "" {
			return DobyToknames[c-4]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func DobyStatname(s int) string {
	if s >= 0 && s < len(DobyStatenames) {
		if DobyStatenames[s] != "" {
			return DobyStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func Dobylex1(lex DobyLexer, lval *DobySymType) int {
	c := 0
	char := lex.Lex(lval)
	if char <= 0 {
		c = DobyTok1[0]
		goto out
	}
	if char < len(DobyTok1) {
		c = DobyTok1[char]
		goto out
	}
	if char >= DobyPrivate {
		if char < DobyPrivate+len(DobyTok2) {
			c = DobyTok2[char-DobyPrivate]
			goto out
		}
	}
	for i := 0; i < len(DobyTok3); i += 2 {
		c = DobyTok3[i+0]
		if c == char {
			c = DobyTok3[i+1]
			goto out
		}
	}

out:
	if c == 0 {
		c = DobyTok2[1] /* unknown char */
	}
	if DobyDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", DobyTokname(c), uint(char))
	}
	return c
}

func DobyParse(Dobylex DobyLexer) int {
	var Dobyn int
	var Dobylval DobySymType
	var DobyVAL DobySymType
	DobyS := make([]DobySymType, DobyMaxDepth)

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	Dobystate := 0
	Dobychar := -1
	Dobyp := -1
	goto Dobystack

ret0:
	return 0

ret1:
	return 1

Dobystack:
	/* put a state and value onto the stack */
	if DobyDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", DobyTokname(Dobychar), DobyStatname(Dobystate))
	}

	Dobyp++
	if Dobyp >= len(DobyS) {
		nyys := make([]DobySymType, len(DobyS)*2)
		copy(nyys, DobyS)
		DobyS = nyys
	}
	DobyS[Dobyp] = DobyVAL
	DobyS[Dobyp].yys = Dobystate

Dobynewstate:
	Dobyn = DobyPact[Dobystate]
	if Dobyn <= DobyFlag {
		goto Dobydefault /* simple state */
	}
	if Dobychar < 0 {
		Dobychar = Dobylex1(Dobylex, &Dobylval)
	}
	Dobyn += Dobychar
	if Dobyn < 0 || Dobyn >= DobyLast {
		goto Dobydefault
	}
	Dobyn = DobyAct[Dobyn]
	if DobyChk[Dobyn] == Dobychar { /* valid shift */
		Dobychar = -1
		DobyVAL = Dobylval
		Dobystate = Dobyn
		if Errflag > 0 {
			Errflag--
		}
		goto Dobystack
	}

Dobydefault:
	/* default state action */
	Dobyn = DobyDef[Dobystate]
	if Dobyn == -2 {
		if Dobychar < 0 {
			Dobychar = Dobylex1(Dobylex, &Dobylval)
		}

		/* look through exception table */
		xi := 0
		for {
			if DobyExca[xi+0] == -1 && DobyExca[xi+1] == Dobystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			Dobyn = DobyExca[xi+0]
			if Dobyn < 0 || Dobyn == Dobychar {
				break
			}
		}
		Dobyn = DobyExca[xi+1]
		if Dobyn < 0 {
			goto ret0
		}
	}
	if Dobyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			Dobylex.Error("syntax error")
			Nerrs++
			if DobyDebug >= 1 {
				__yyfmt__.Printf("%s", DobyStatname(Dobystate))
				__yyfmt__.Printf(" saw %s\n", DobyTokname(Dobychar))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for Dobyp >= 0 {
				Dobyn = DobyPact[DobyS[Dobyp].yys] + DobyErrCode
				if Dobyn >= 0 && Dobyn < DobyLast {
					Dobystate = DobyAct[Dobyn] /* simulate a shift of "error" */
					if DobyChk[Dobystate] == DobyErrCode {
						goto Dobystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if DobyDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", DobyS[Dobyp].yys)
				}
				Dobyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if DobyDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", DobyTokname(Dobychar))
			}
			if Dobychar == DobyEofCode {
				goto ret1
			}
			Dobychar = -1
			goto Dobynewstate /* try again in the same state */
		}
	}

	/* reduction by production Dobyn */
	if DobyDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", Dobyn, DobyStatname(Dobystate))
	}

	Dobynt := Dobyn
	Dobypt := Dobyp
	_ = Dobypt // guard against "declared and not used"

	Dobyp -= DobyR2[Dobyn]
	DobyVAL = DobyS[Dobyp+1]

	/* consult goto table to find next state */
	Dobyn = DobyR1[Dobyn]
	Dobyg := DobyPgo[Dobyn]
	Dobyj := Dobyg + DobyS[Dobyp].yys + 1

	if Dobyj >= DobyLast {
		Dobystate = DobyAct[Dobyg]
	} else {
		Dobystate = DobyAct[Dobyj]
		if DobyChk[Dobystate] != -Dobyn {
			Dobystate = DobyAct[Dobyg]
		}
	}
	// dummy call; replaced with literal code
	switch Dobynt {

	case 1:
		//line grammar.y:90
		{ DobyVAL.expr = &ast.Ident{DobyS[Dobypt-0].tok.Pos, DobyS[Dobypt-0].tok.Lit} }
	case 2:
		//line grammar.y:92
		{ DobyVAL.expr = &ast.BasicLit{DobyS[Dobypt-0].tok.Pos, token.INT, DobyS[Dobypt-0].tok.Lit} }
	case 3:
		//line grammar.y:93
		{ DobyVAL.expr = &ast.BasicLit{DobyS[Dobypt-0].tok.Pos, token.FLOAT, DobyS[Dobypt-0].tok.Lit} }
	case 4:
		//line grammar.y:94
		{ DobyVAL.expr = &ast.BasicLit{DobyS[Dobypt-0].tok.Pos, token.STRING, DobyS[Dobypt-0].tok.Lit} }
	case 5:
		//line grammar.y:95
		{ DobyVAL.expr = &ast.BasicLit{DobyS[Dobypt-0].tok.Pos, token.CHAR, DobyS[Dobypt-0].tok.Lit} }
	case 6:
		//line grammar.y:97
		{ DobyVAL.expr = &ast.ParenExpr{DobyS[Dobypt-2].tok.Pos, DobyS[Dobypt-1].expr, DobyS[Dobypt-0].tok.Pos} }
	case 7:
		//line grammar.y:99
		{ DobyVAL.expr = &ast.SelectorExpr{DobyS[Dobypt-2].expr, DobyS[Dobypt-0].expr.(*ast.Ident)} }
	case 8:
		//line grammar.y:102
		{ DobyVAL.expr = &ast.SliceExpr{DobyS[Dobypt-5].expr, DobyS[Dobypt-4].tok.Pos, DobyS[Dobypt-3].expr, DobyS[Dobypt-1].expr, DobyS[Dobypt-0].tok.Pos} }
	case 9:
		//line grammar.y:104
		{ DobyVAL.expr = &ast.SliceExpr{DobyS[Dobypt-4].expr, DobyS[Dobypt-3].tok.Pos, nil, DobyS[Dobypt-1].expr, DobyS[Dobypt-0].tok.Pos} }
	case 10:
		//line grammar.y:106
		{ DobyVAL.expr = &ast.SliceExpr{DobyS[Dobypt-4].expr, DobyS[Dobypt-3].tok.Pos, DobyS[Dobypt-2].expr, nil, DobyS[Dobypt-0].tok.Pos} }
	case 11:
		//line grammar.y:109
		{ DobyVAL.expr = &ast.IndexExpr{DobyS[Dobypt-3].expr, DobyS[Dobypt-2].tok.Pos, DobyS[Dobypt-1].expr, DobyS[Dobypt-2].tok.Pos} }
	case 12:
		//line grammar.y:111
		{ DobyVAL.expr_list = []ast.Expr{} }
	case 13:
		//line grammar.y:112
		{ DobyVAL.expr_list = []ast.Expr{DobyS[Dobypt-0].expr} }
	case 14:
		//line grammar.y:113
		{ DobyVAL.expr_list = append(DobyS[Dobypt-2].expr_list, DobyS[Dobypt-0].expr) }
	case 15:
		//line grammar.y:114
		{ DobyVAL.expr_list = append(DobyS[Dobypt-3].expr_list, DobyS[Dobypt-0].expr) }
	case 16:
		//line grammar.y:116
		{ DobyVAL.expr = &ast.CallExpr{DobyS[Dobypt-3].expr, DobyS[Dobypt-2].tok.Pos, DobyS[Dobypt-1].expr_list, DobyS[Dobypt-0].tok.Pos} }
	case 17:
		//line grammar.y:118
		{ DobyVAL.expr = &ast.UnaryExpr{DobyS[Dobypt-1].tok.Pos, token.SUB, DobyS[Dobypt-0].expr } }
	case 18:
		//line grammar.y:119
		{ DobyVAL.expr = &ast.UnaryExpr{DobyS[Dobypt-1].tok.Pos, token.NOT, DobyS[Dobypt-0].expr } }
	case 19:
		//line grammar.y:121
		{ DobyVAL.expr = &ast.BinaryExpr{DobyS[Dobypt-2].expr, DobyS[Dobypt-1].tok.Pos, token.ADD, DobyS[Dobypt-0].expr } }
	case 20:
		//line grammar.y:122
		{ DobyVAL.expr = &ast.BinaryExpr{DobyS[Dobypt-2].expr, DobyS[Dobypt-1].tok.Pos, token.SUB, DobyS[Dobypt-0].expr } }
	case 21:
		//line grammar.y:123
		{ DobyVAL.expr = &ast.BinaryExpr{DobyS[Dobypt-2].expr, DobyS[Dobypt-1].tok.Pos, token.MUL, DobyS[Dobypt-0].expr } }
	case 22:
		//line grammar.y:124
		{ DobyVAL.expr = &ast.BinaryExpr{DobyS[Dobypt-2].expr, DobyS[Dobypt-1].tok.Pos, token.QUO, DobyS[Dobypt-0].expr } }
	case 23:
		//line grammar.y:125
		{ DobyVAL.expr = &ast.BinaryExpr{DobyS[Dobypt-2].expr, DobyS[Dobypt-1].tok.Pos, token.REM, DobyS[Dobypt-0].expr } }
	case 24:
		//line grammar.y:126
		{ DobyVAL.expr = &ast.BinaryExpr{DobyS[Dobypt-2].expr, DobyS[Dobypt-1].tok.Pos, token.AND, DobyS[Dobypt-0].expr } }
	case 25:
		//line grammar.y:127
		{ DobyVAL.expr = &ast.BinaryExpr{DobyS[Dobypt-2].expr, DobyS[Dobypt-1].tok.Pos, token.OR, DobyS[Dobypt-0].expr } }
	case 26:
		//line grammar.y:128
		{ DobyVAL.expr = &ast.BinaryExpr{DobyS[Dobypt-2].expr, DobyS[Dobypt-1].tok.Pos, token.XOR, DobyS[Dobypt-0].expr } }
	case 27:
		//line grammar.y:129
		{ DobyVAL.expr = &ast.BinaryExpr{DobyS[Dobypt-2].expr, DobyS[Dobypt-1].tok.Pos, token.SHL, DobyS[Dobypt-0].expr } }
	case 28:
		//line grammar.y:130
		{ DobyVAL.expr = &ast.BinaryExpr{DobyS[Dobypt-2].expr, DobyS[Dobypt-1].tok.Pos, token.SHR, DobyS[Dobypt-0].expr } }
	case 29:
		//line grammar.y:131
		{ DobyVAL.expr = &ast.BinaryExpr{DobyS[Dobypt-2].expr, DobyS[Dobypt-1].tok.Pos, token.AND_NOT, DobyS[Dobypt-0].expr } }
	case 30:
		//line grammar.y:132
		{ DobyVAL.expr = &ast.BinaryExpr{DobyS[Dobypt-2].expr, DobyS[Dobypt-1].tok.Pos, token.LSS, DobyS[Dobypt-0].expr } }
	case 31:
		//line grammar.y:133
		{ DobyVAL.expr = &ast.BinaryExpr{DobyS[Dobypt-2].expr, DobyS[Dobypt-1].tok.Pos, token.GTR, DobyS[Dobypt-0].expr } }
	case 32:
		//line grammar.y:134
		{ DobyVAL.expr = &ast.BinaryExpr{DobyS[Dobypt-2].expr, DobyS[Dobypt-1].tok.Pos, token.NEQ, DobyS[Dobypt-0].expr } }
	case 33:
		//line grammar.y:135
		{ DobyVAL.expr = &ast.BinaryExpr{DobyS[Dobypt-2].expr, DobyS[Dobypt-1].tok.Pos, token.LEQ, DobyS[Dobypt-0].expr } }
	case 34:
		//line grammar.y:136
		{ DobyVAL.expr = &ast.BinaryExpr{DobyS[Dobypt-2].expr, DobyS[Dobypt-1].tok.Pos, token.GEQ, DobyS[Dobypt-0].expr } }
	case 35:
		//line grammar.y:137
		{ DobyVAL.expr = &ast.BinaryExpr{DobyS[Dobypt-2].expr, DobyS[Dobypt-1].tok.Pos, token.EQL, DobyS[Dobypt-0].expr } }
	case 36:
		//line grammar.y:139
		{ DobyVAL.expr = &ast.BinaryExpr{DobyS[Dobypt-2].expr, DobyS[Dobypt-1].tok.Pos, token.LAND, DobyS[Dobypt-0].expr } }
	case 37:
		//line grammar.y:140
		{ DobyVAL.expr = &ast.BinaryExpr{DobyS[Dobypt-2].expr, DobyS[Dobypt-1].tok.Pos, token.LOR, DobyS[Dobypt-0].expr } }
	case 38:
		//line grammar.y:143
		{ DobyVAL.expr = &ast.ArrayExpr{DobyS[Dobypt-2].tok.Pos, DobyS[Dobypt-1].expr_list, DobyS[Dobypt-0].tok.Pos} }
	case 39:
		//line grammar.y:145
		{ DobyVAL.expr = &ast.ArrayExpr{DobyS[Dobypt-4].tok.Pos, DobyS[Dobypt-2].expr_list, DobyS[Dobypt-1].tok.Pos} }
	case 40:
		//line grammar.y:147
		{ DobyVAL.expr = &ast.ArrayExpr{DobyS[Dobypt-3].tok.Pos, DobyS[Dobypt-1].expr_list, DobyS[Dobypt-0].tok.Pos} }
	case 41:
		//line grammar.y:150
		{ DobyVAL.expr = &ast.SetExpr{DobyS[Dobypt-2].tok.Pos, DobyS[Dobypt-1].expr_list, DobyS[Dobypt-0].tok.Pos} }
	case 42:
		//line grammar.y:152
		{ DobyVAL.expr = &ast.SetExpr{DobyS[Dobypt-4].tok.Pos, DobyS[Dobypt-2].expr_list, DobyS[Dobypt-0].tok.Pos} }
	case 43:
		//line grammar.y:154
		{ DobyVAL.expr = &ast.SetExpr{DobyS[Dobypt-3].tok.Pos, DobyS[Dobypt-1].expr_list, DobyS[Dobypt-0].tok.Pos} }
	case 44:
		//line grammar.y:157
		{ DobyVAL.field = &ast.Field{DobyS[Dobypt-2].expr, DobyS[Dobypt-1].tok.Pos, DobyS[Dobypt-0].expr} }
	case 45:
		//line grammar.y:159
		{ DobyVAL.field_list = []*ast.Field{} }
	case 46:
		//line grammar.y:160
		{ DobyVAL.field_list = []*ast.Field{DobyS[Dobypt-0].field} }
	case 47:
		//line grammar.y:161
		{ DobyVAL.field_list = append(DobyS[Dobypt-2].field_list, DobyS[Dobypt-0].field) }
	case 48:
		//line grammar.y:162
		{ DobyVAL.field_list = append(DobyS[Dobypt-2].field_list, DobyS[Dobypt-0].field) }
	case 49:
		//line grammar.y:163
		{ DobyVAL.field_list = append(DobyS[Dobypt-3].field_list, DobyS[Dobypt-0].field) }
	case 50:
		//line grammar.y:164
		{ DobyVAL.field_list = DobyS[Dobypt-1].field_list }
	case 51:
		//line grammar.y:165
		{ DobyVAL.field_list = DobyS[Dobypt-2].field_list }
	case 52:
		//line grammar.y:168
		{ DobyVAL.expr = &ast.DictExpr{DobyS[Dobypt-2].tok.Pos, DobyS[Dobypt-1].field_list, DobyS[Dobypt-0].tok.Pos} }
	case 53:
		//line grammar.y:171
		{ DobyVAL.ident_list = []*ast.Ident{} }
	case 54:
		//line grammar.y:173
		{ DobyVAL.ident_list = []*ast.Ident{&ast.Ident{DobyS[Dobypt-0].tok.Pos, DobyS[Dobypt-0].tok.Lit}} }
	case 55:
		//line grammar.y:175
		{ DobyVAL.ident_list = append(DobyS[Dobypt-2].ident_list, &ast.Ident{DobyS[Dobypt-0].tok.Pos, DobyS[Dobypt-0].tok.Lit}) }
	case 56:
		//line grammar.y:178
		{ DobyVAL.expr = &ast.FuncDeclExpr{DobyS[Dobypt-4].tok.Pos, nil, nil, nil, DobyS[Dobypt-2].ident_list, DobyS[Dobypt-0].stmt.(*ast.BlockStmt), []string{}} }
	case 57:
		//line grammar.y:180
		{ DobyVAL.expr = &ast.FuncDeclExpr{DobyS[Dobypt-5].tok.Pos, nil, nil, &ast.Ident{DobyS[Dobypt-4].tok.Pos, DobyS[Dobypt-4].tok.Lit}, DobyS[Dobypt-2].ident_list, DobyS[Dobypt-0].stmt.(*ast.BlockStmt), []string{}} }
	case 58:
		//line grammar.y:182
		{ DobyVAL.expr = &ast.FuncDeclExpr{DobyS[Dobypt-9].tok.Pos, &ast.Ident{DobyS[Dobypt-7].tok.Pos, DobyS[Dobypt-7].tok.Lit}, &ast.Ident{DobyS[Dobypt-6].tok.Pos, DobyS[Dobypt-6].tok.Lit},
	                                          &ast.Ident{DobyS[Dobypt-4].tok.Pos, DobyS[Dobypt-4].tok.Lit}, DobyS[Dobypt-2].ident_list, DobyS[Dobypt-0].stmt.(*ast.BlockStmt), []string{}} }
	case 59:
		DobyVAL.expr = DobyS[Dobypt-0].expr
	case 60:
		DobyVAL.expr = DobyS[Dobypt-0].expr
	case 61:
		DobyVAL.expr = DobyS[Dobypt-0].expr
	case 62:
		DobyVAL.expr = DobyS[Dobypt-0].expr
	case 63:
		DobyVAL.expr = DobyS[Dobypt-0].expr
	case 64:
		DobyVAL.expr = DobyS[Dobypt-0].expr
	case 65:
		DobyVAL.expr = DobyS[Dobypt-0].expr
	case 66:
		DobyVAL.expr = DobyS[Dobypt-0].expr
	case 67:
		DobyVAL.expr = DobyS[Dobypt-0].expr
	case 68:
		DobyVAL.expr = DobyS[Dobypt-0].expr
	case 69:
		DobyVAL.expr = DobyS[Dobypt-0].expr
	case 70:
		DobyVAL.expr = DobyS[Dobypt-0].expr
	case 71:
		DobyVAL.expr = DobyS[Dobypt-0].expr
	case 72:
		//line grammar.y:201
		{ DobyVAL.stmt = &ast.ExprStmt{DobyS[Dobypt-0].expr} }
	case 73:
		//line grammar.y:203
		{ DobyVAL.stmt = &ast.SendStmt{DobyS[Dobypt-2].expr, DobyS[Dobypt-1].tok.Pos, DobyS[Dobypt-0].expr} }
	case 74:
		//line grammar.y:205
		{ DobyVAL.stmt = &ast.IncDecStmt{DobyS[Dobypt-1].expr, DobyS[Dobypt-0].tok.Pos, token.INC} }
	case 75:
		//line grammar.y:206
		{ DobyVAL.stmt = &ast.IncDecStmt{DobyS[Dobypt-1].expr, DobyS[Dobypt-0].tok.Pos, token.DEC} }
	case 76:
		//line grammar.y:208
		{ DobyVAL.stmt = &ast.AssignStmt{DobyS[Dobypt-2].expr_list, DobyS[Dobypt-1].tok.Pos, token.ASSIGN, DobyS[Dobypt-0].expr_list} }
	case 77:
		//line grammar.y:209
		{ DobyVAL.stmt = &ast.AssignStmt{DobyS[Dobypt-2].expr_list, DobyS[Dobypt-1].tok.Pos, token.ADD_ASSIGN, DobyS[Dobypt-0].expr_list} }
	case 78:
		//line grammar.y:210
		{ DobyVAL.stmt = &ast.AssignStmt{DobyS[Dobypt-2].expr_list, DobyS[Dobypt-1].tok.Pos, token.SUB_ASSIGN, DobyS[Dobypt-0].expr_list} }
	case 79:
		//line grammar.y:211
		{ DobyVAL.stmt = &ast.AssignStmt{DobyS[Dobypt-2].expr_list, DobyS[Dobypt-1].tok.Pos, token.MUL_ASSIGN, DobyS[Dobypt-0].expr_list} }
	case 80:
		//line grammar.y:212
		{ DobyVAL.stmt = &ast.AssignStmt{DobyS[Dobypt-2].expr_list, DobyS[Dobypt-1].tok.Pos, token.QUO_ASSIGN, DobyS[Dobypt-0].expr_list} }
	case 81:
		//line grammar.y:213
		{ DobyVAL.stmt = &ast.AssignStmt{DobyS[Dobypt-2].expr_list, DobyS[Dobypt-1].tok.Pos, token.REM_ASSIGN, DobyS[Dobypt-0].expr_list} }
	case 82:
		//line grammar.y:214
		{ DobyVAL.stmt = &ast.AssignStmt{DobyS[Dobypt-2].expr_list, DobyS[Dobypt-1].tok.Pos, token.AND_ASSIGN, DobyS[Dobypt-0].expr_list} }
	case 83:
		//line grammar.y:215
		{ DobyVAL.stmt = &ast.AssignStmt{DobyS[Dobypt-2].expr_list, DobyS[Dobypt-1].tok.Pos, token.OR_ASSIGN, DobyS[Dobypt-0].expr_list} }
	case 84:
		//line grammar.y:216
		{ DobyVAL.stmt = &ast.AssignStmt{DobyS[Dobypt-2].expr_list, DobyS[Dobypt-1].tok.Pos, token.XOR_ASSIGN, DobyS[Dobypt-0].expr_list} }
	case 85:
		//line grammar.y:217
		{ DobyVAL.stmt = &ast.AssignStmt{DobyS[Dobypt-2].expr_list, DobyS[Dobypt-1].tok.Pos, token.SHL_ASSIGN, DobyS[Dobypt-0].expr_list} }
	case 86:
		//line grammar.y:218
		{ DobyVAL.stmt = &ast.AssignStmt{DobyS[Dobypt-2].expr_list, DobyS[Dobypt-1].tok.Pos, token.SHR_ASSIGN, DobyS[Dobypt-0].expr_list} }
	case 87:
		//line grammar.y:219
		{ DobyVAL.stmt = &ast.AssignStmt{DobyS[Dobypt-2].expr_list, DobyS[Dobypt-1].tok.Pos, token.AND_NOT_ASSIGN, DobyS[Dobypt-0].expr_list} }
	case 88:
		//line grammar.y:222
		{ DobyVAL.stmt = &ast.GoStmt{DobyS[Dobypt-1].tok.Pos, DobyS[Dobypt-0].expr.(*ast.CallExpr)} }
	case 89:
		//line grammar.y:225
		{ DobyVAL.stmt = &ast.ReturnStmt{DobyS[Dobypt-1].tok.Pos, DobyS[Dobypt-0].expr_list} }
	case 90:
		//line grammar.y:227
		{ DobyVAL.stmt = &ast.BranchStmt{DobyS[Dobypt-0].tok.Pos, token.BREAK} }
	case 91:
		//line grammar.y:228
		{ DobyVAL.stmt = &ast.BranchStmt{DobyS[Dobypt-0].tok.Pos, token.CONTINUE } }
	case 92:
		//line grammar.y:230
		{ DobyVAL.stmt = &ast.BlockStmt{DobyS[Dobypt-2].tok.Pos, DobyS[Dobypt-1].stmt_list ,DobyS[Dobypt-0].tok.Pos} }
	case 93:
		//line grammar.y:232
		{ DobyVAL.stmt = &ast.IfStmt{DobyS[Dobypt-2].tok.Pos, DobyS[Dobypt-1].expr, DobyS[Dobypt-0].stmt.(*ast.BlockStmt), nil} }
	case 94:
		//line grammar.y:233
		{ DobyVAL.stmt = &ast.IfStmt{DobyS[Dobypt-4].tok.Pos, DobyS[Dobypt-3].expr, DobyS[Dobypt-2].stmt.(*ast.BlockStmt), DobyS[Dobypt-0].stmt} }
	case 95:
		//line grammar.y:235
		{ DobyVAL.stmt = &ast.CaseClause{DobyS[Dobypt-3].tok.Pos, DobyS[Dobypt-2].expr_list, DobyS[Dobypt-1].tok.Pos, DobyS[Dobypt-0].stmt_list} }
	case 96:
		//line grammar.y:236
		{ DobyVAL.stmt = &ast.CaseClause{DobyS[Dobypt-2].tok.Pos, nil, DobyS[Dobypt-1].tok.Pos, DobyS[Dobypt-0].stmt_list} }
	case 97:
		//line grammar.y:238
		{ DobyVAL.stmt_list = []ast.Stmt{} }
	case 98:
		//line grammar.y:239
		{ DobyVAL.stmt_list = []ast.Stmt{DobyS[Dobypt-0].stmt} }
	case 99:
		//line grammar.y:240
		{ DobyVAL.stmt_list = append(DobyS[Dobypt-1].stmt_list, DobyS[Dobypt-0].stmt) }
	case 100:
		//line grammar.y:242
		{ DobyVAL.stmt = &ast.BlockStmt{DobyS[Dobypt-2].tok.Pos, DobyS[Dobypt-1].stmt_list, DobyS[Dobypt-0].tok.Pos} }
	case 101:
		//line grammar.y:244
		{ DobyVAL.stmt = &ast.SwitchStmt{DobyS[Dobypt-2].tok.Pos, DobyS[Dobypt-1].stmt, DobyS[Dobypt-0].stmt.(*ast.BlockStmt)} }
	case 102:
		//line grammar.y:246
		{ DobyVAL.stmt = &ast.SelectStmt{DobyS[Dobypt-1].tok.Pos, DobyS[Dobypt-0].stmt.(*ast.BlockStmt)} }
	case 103:
		//line grammar.y:249
		{ DobyVAL.stmt = &ast.ForStmt{DobyS[Dobypt-6].tok.Pos, DobyS[Dobypt-5].stmt, DobyS[Dobypt-3].expr, DobyS[Dobypt-1].stmt, DobyS[Dobypt-0].stmt.(*ast.BlockStmt)} }
	case 104:
		//line grammar.y:251
		{ DobyVAL.stmt = &ast.ForStmt{DobyS[Dobypt-5].tok.Pos, nil, DobyS[Dobypt-3].expr, DobyS[Dobypt-1].stmt, DobyS[Dobypt-0].stmt.(*ast.BlockStmt)} }
	case 105:
		//line grammar.y:253
		{ DobyVAL.stmt = &ast.ForStmt{DobyS[Dobypt-2].tok.Pos, nil, DobyS[Dobypt-1].expr, nil, DobyS[Dobypt-0].stmt.(*ast.BlockStmt)} }
	case 106:
		//line grammar.y:256
		{ DobyVAL.stmt = &ast.RangeStmt{DobyS[Dobypt-5].tok.Pos, DobyS[Dobypt-4].expr_list, DobyS[Dobypt-1].expr, DobyS[Dobypt-0].stmt.(*ast.BlockStmt)} }
	case 107:
		//line grammar.y:259
		{ DobyVAL.stmt = &ast.ImportStmt{DobyS[Dobypt-1].tok.Pos, []string{DobyS[Dobypt-0].tok.Lit}} }
	case 108:
		DobyVAL.stmt = DobyS[Dobypt-0].stmt
	case 109:
		DobyVAL.stmt = DobyS[Dobypt-0].stmt
	case 110:
		DobyVAL.stmt = DobyS[Dobypt-0].stmt
	case 111:
		DobyVAL.stmt = DobyS[Dobypt-0].stmt
	case 112:
		DobyVAL.stmt = DobyS[Dobypt-0].stmt
	case 113:
		DobyVAL.stmt = DobyS[Dobypt-0].stmt
	case 114:
		DobyVAL.stmt = DobyS[Dobypt-0].stmt
	case 115:
		DobyVAL.stmt = DobyS[Dobypt-0].stmt
	case 116:
		DobyVAL.stmt = DobyS[Dobypt-0].stmt
	case 117:
		DobyVAL.stmt = DobyS[Dobypt-0].stmt
	case 118:
		DobyVAL.stmt = DobyS[Dobypt-0].stmt
	case 119:
		DobyVAL.stmt = DobyS[Dobypt-0].stmt
	case 120:
		DobyVAL.stmt = DobyS[Dobypt-0].stmt
	case 121:
		DobyVAL.stmt = DobyS[Dobypt-0].stmt
	case 122:
		//line grammar.y:276
		{ DobyVAL.stmt_list = []ast.Stmt{} }
	case 123:
		//line grammar.y:277
		{ DobyVAL.stmt_list = []ast.Stmt{DobyS[Dobypt-0].stmt} }
	case 124:
		//line grammar.y:278
		{ DobyVAL.stmt_list = append(DobyS[Dobypt-2].stmt_list, DobyS[Dobypt-0].stmt) }
	case 125:
		//line grammar.y:279
		{ DobyVAL.stmt_list = append(DobyS[Dobypt-2].stmt_list, DobyS[Dobypt-0].stmt) }
	case 126:
		//line grammar.y:280
		{ DobyVAL.stmt_list = DobyS[Dobypt-1].stmt_list }
	case 127:
		//line grammar.y:285
		{ ProgramAst = DobyS[Dobypt-1].stmt_list }
	}
	goto Dobystack /* stack new state and value */
}
