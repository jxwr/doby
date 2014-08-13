
//line grammer.y:2

package parser
import __yyfmt__ "fmt"
//line grammer.y:3
		
import (
       "github.com/jxwr/doubi/ast"
       "github.com/jxwr/doubi/token"
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


//line grammer.y:27
type DoubiSymType struct {
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

var DoubiToknames = []string{
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
var DoubiStatenames = []string{}

const DoubiEofCode = 1
const DoubiErrCode = 2
const DoubiMaxDepth = 200

//line yacctab:1
var DoubiExca = []int{
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

const DoubiNprod = 128
const DoubiPrivate = 57344

var DoubiTokenNames []string
var DoubiStates []string

const DoubiLast = 1131

var DoubiAct = []int{

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
var DoubiPact = []int{

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
var DoubiPgo = []int{

	0, 0, 24, 213, 211, 210, 209, 208, 207, 116,
	205, 201, 196, 195, 192, 168, 4, 188, 3, 57,
	183, 181, 172, 169, 158, 144, 143, 2, 138, 8,
	5, 133, 132, 131, 127, 126, 7, 125, 120,
}
var DoubiR1 = []int{

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
var DoubiR2 = []int{

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
var DoubiChk = []int{

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
var DoubiDef = []int{

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
var DoubiTok1 = []int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 85,
}
var DoubiTok2 = []int{

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
var DoubiTok3 = []int{
	0,
}

//line yaccpar:1

/*	parser for yacc output	*/

var DoubiDebug = 0

type DoubiLexer interface {
	Lex(lval *DoubiSymType) int
	Error(s string)
}

const DoubiFlag = -1000

func DoubiTokname(c int) string {
	// 4 is TOKSTART above
	if c >= 4 && c-4 < len(DoubiToknames) {
		if DoubiToknames[c-4] != "" {
			return DoubiToknames[c-4]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func DoubiStatname(s int) string {
	if s >= 0 && s < len(DoubiStatenames) {
		if DoubiStatenames[s] != "" {
			return DoubiStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func Doubilex1(lex DoubiLexer, lval *DoubiSymType) int {
	c := 0
	char := lex.Lex(lval)
	if char <= 0 {
		c = DoubiTok1[0]
		goto out
	}
	if char < len(DoubiTok1) {
		c = DoubiTok1[char]
		goto out
	}
	if char >= DoubiPrivate {
		if char < DoubiPrivate+len(DoubiTok2) {
			c = DoubiTok2[char-DoubiPrivate]
			goto out
		}
	}
	for i := 0; i < len(DoubiTok3); i += 2 {
		c = DoubiTok3[i+0]
		if c == char {
			c = DoubiTok3[i+1]
			goto out
		}
	}

out:
	if c == 0 {
		c = DoubiTok2[1] /* unknown char */
	}
	if DoubiDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", DoubiTokname(c), uint(char))
	}
	return c
}

func DoubiParse(Doubilex DoubiLexer) int {
	var Doubin int
	var Doubilval DoubiSymType
	var DoubiVAL DoubiSymType
	DoubiS := make([]DoubiSymType, DoubiMaxDepth)

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	Doubistate := 0
	Doubichar := -1
	Doubip := -1
	goto Doubistack

ret0:
	return 0

ret1:
	return 1

Doubistack:
	/* put a state and value onto the stack */
	if DoubiDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", DoubiTokname(Doubichar), DoubiStatname(Doubistate))
	}

	Doubip++
	if Doubip >= len(DoubiS) {
		nyys := make([]DoubiSymType, len(DoubiS)*2)
		copy(nyys, DoubiS)
		DoubiS = nyys
	}
	DoubiS[Doubip] = DoubiVAL
	DoubiS[Doubip].yys = Doubistate

Doubinewstate:
	Doubin = DoubiPact[Doubistate]
	if Doubin <= DoubiFlag {
		goto Doubidefault /* simple state */
	}
	if Doubichar < 0 {
		Doubichar = Doubilex1(Doubilex, &Doubilval)
	}
	Doubin += Doubichar
	if Doubin < 0 || Doubin >= DoubiLast {
		goto Doubidefault
	}
	Doubin = DoubiAct[Doubin]
	if DoubiChk[Doubin] == Doubichar { /* valid shift */
		Doubichar = -1
		DoubiVAL = Doubilval
		Doubistate = Doubin
		if Errflag > 0 {
			Errflag--
		}
		goto Doubistack
	}

Doubidefault:
	/* default state action */
	Doubin = DoubiDef[Doubistate]
	if Doubin == -2 {
		if Doubichar < 0 {
			Doubichar = Doubilex1(Doubilex, &Doubilval)
		}

		/* look through exception table */
		xi := 0
		for {
			if DoubiExca[xi+0] == -1 && DoubiExca[xi+1] == Doubistate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			Doubin = DoubiExca[xi+0]
			if Doubin < 0 || Doubin == Doubichar {
				break
			}
		}
		Doubin = DoubiExca[xi+1]
		if Doubin < 0 {
			goto ret0
		}
	}
	if Doubin == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			Doubilex.Error("syntax error")
			Nerrs++
			if DoubiDebug >= 1 {
				__yyfmt__.Printf("%s", DoubiStatname(Doubistate))
				__yyfmt__.Printf(" saw %s\n", DoubiTokname(Doubichar))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for Doubip >= 0 {
				Doubin = DoubiPact[DoubiS[Doubip].yys] + DoubiErrCode
				if Doubin >= 0 && Doubin < DoubiLast {
					Doubistate = DoubiAct[Doubin] /* simulate a shift of "error" */
					if DoubiChk[Doubistate] == DoubiErrCode {
						goto Doubistack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if DoubiDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", DoubiS[Doubip].yys)
				}
				Doubip--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if DoubiDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", DoubiTokname(Doubichar))
			}
			if Doubichar == DoubiEofCode {
				goto ret1
			}
			Doubichar = -1
			goto Doubinewstate /* try again in the same state */
		}
	}

	/* reduction by production Doubin */
	if DoubiDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", Doubin, DoubiStatname(Doubistate))
	}

	Doubint := Doubin
	Doubipt := Doubip
	_ = Doubipt // guard against "declared and not used"

	Doubip -= DoubiR2[Doubin]
	DoubiVAL = DoubiS[Doubip+1]

	/* consult goto table to find next state */
	Doubin = DoubiR1[Doubin]
	Doubig := DoubiPgo[Doubin]
	Doubij := Doubig + DoubiS[Doubip].yys + 1

	if Doubij >= DoubiLast {
		Doubistate = DoubiAct[Doubig]
	} else {
		Doubistate = DoubiAct[Doubij]
		if DoubiChk[Doubistate] != -Doubin {
			Doubistate = DoubiAct[Doubig]
		}
	}
	// dummy call; replaced with literal code
	switch Doubint {

	case 1:
		//line grammer.y:90
		{ DoubiVAL.expr = &ast.Ident{DoubiS[Doubipt-0].tok.Pos, DoubiS[Doubipt-0].tok.Lit} }
	case 2:
		//line grammer.y:92
		{ DoubiVAL.expr = &ast.BasicLit{DoubiS[Doubipt-0].tok.Pos, token.INT, DoubiS[Doubipt-0].tok.Lit} }
	case 3:
		//line grammer.y:93
		{ DoubiVAL.expr = &ast.BasicLit{DoubiS[Doubipt-0].tok.Pos, token.FLOAT, DoubiS[Doubipt-0].tok.Lit} }
	case 4:
		//line grammer.y:94
		{ DoubiVAL.expr = &ast.BasicLit{DoubiS[Doubipt-0].tok.Pos, token.STRING, DoubiS[Doubipt-0].tok.Lit} }
	case 5:
		//line grammer.y:95
		{ DoubiVAL.expr = &ast.BasicLit{DoubiS[Doubipt-0].tok.Pos, token.CHAR, DoubiS[Doubipt-0].tok.Lit} }
	case 6:
		//line grammer.y:97
		{ DoubiVAL.expr = &ast.ParenExpr{DoubiS[Doubipt-2].tok.Pos, DoubiS[Doubipt-1].expr, DoubiS[Doubipt-0].tok.Pos} }
	case 7:
		//line grammer.y:99
		{ DoubiVAL.expr = &ast.SelectorExpr{DoubiS[Doubipt-2].expr, DoubiS[Doubipt-0].expr.(*ast.Ident)} }
	case 8:
		//line grammer.y:102
		{ DoubiVAL.expr = &ast.SliceExpr{DoubiS[Doubipt-5].expr, DoubiS[Doubipt-4].tok.Pos, DoubiS[Doubipt-3].expr, DoubiS[Doubipt-1].expr, DoubiS[Doubipt-0].tok.Pos} }
	case 9:
		//line grammer.y:104
		{ DoubiVAL.expr = &ast.SliceExpr{DoubiS[Doubipt-4].expr, DoubiS[Doubipt-3].tok.Pos, nil, DoubiS[Doubipt-1].expr, DoubiS[Doubipt-0].tok.Pos} }
	case 10:
		//line grammer.y:106
		{ DoubiVAL.expr = &ast.SliceExpr{DoubiS[Doubipt-4].expr, DoubiS[Doubipt-3].tok.Pos, DoubiS[Doubipt-2].expr, nil, DoubiS[Doubipt-0].tok.Pos} }
	case 11:
		//line grammer.y:109
		{ DoubiVAL.expr = &ast.IndexExpr{DoubiS[Doubipt-3].expr, DoubiS[Doubipt-2].tok.Pos, DoubiS[Doubipt-1].expr, DoubiS[Doubipt-2].tok.Pos} }
	case 12:
		//line grammer.y:111
		{ DoubiVAL.expr_list = []ast.Expr{} }
	case 13:
		//line grammer.y:112
		{ DoubiVAL.expr_list = []ast.Expr{DoubiS[Doubipt-0].expr} }
	case 14:
		//line grammer.y:113
		{ DoubiVAL.expr_list = append(DoubiS[Doubipt-2].expr_list, DoubiS[Doubipt-0].expr) }
	case 15:
		//line grammer.y:114
		{ DoubiVAL.expr_list = append(DoubiS[Doubipt-3].expr_list, DoubiS[Doubipt-0].expr) }
	case 16:
		//line grammer.y:116
		{ DoubiVAL.expr = &ast.CallExpr{DoubiS[Doubipt-3].expr, DoubiS[Doubipt-2].tok.Pos, DoubiS[Doubipt-1].expr_list, DoubiS[Doubipt-0].tok.Pos} }
	case 17:
		//line grammer.y:118
		{ DoubiVAL.expr = &ast.UnaryExpr{DoubiS[Doubipt-1].tok.Pos, token.SUB, DoubiS[Doubipt-0].expr } }
	case 18:
		//line grammer.y:119
		{ DoubiVAL.expr = &ast.UnaryExpr{DoubiS[Doubipt-1].tok.Pos, token.NOT, DoubiS[Doubipt-0].expr } }
	case 19:
		//line grammer.y:121
		{ DoubiVAL.expr = &ast.BinaryExpr{DoubiS[Doubipt-2].expr, DoubiS[Doubipt-1].tok.Pos, token.ADD, DoubiS[Doubipt-0].expr } }
	case 20:
		//line grammer.y:122
		{ DoubiVAL.expr = &ast.BinaryExpr{DoubiS[Doubipt-2].expr, DoubiS[Doubipt-1].tok.Pos, token.SUB, DoubiS[Doubipt-0].expr } }
	case 21:
		//line grammer.y:123
		{ DoubiVAL.expr = &ast.BinaryExpr{DoubiS[Doubipt-2].expr, DoubiS[Doubipt-1].tok.Pos, token.MUL, DoubiS[Doubipt-0].expr } }
	case 22:
		//line grammer.y:124
		{ DoubiVAL.expr = &ast.BinaryExpr{DoubiS[Doubipt-2].expr, DoubiS[Doubipt-1].tok.Pos, token.QUO, DoubiS[Doubipt-0].expr } }
	case 23:
		//line grammer.y:125
		{ DoubiVAL.expr = &ast.BinaryExpr{DoubiS[Doubipt-2].expr, DoubiS[Doubipt-1].tok.Pos, token.REM, DoubiS[Doubipt-0].expr } }
	case 24:
		//line grammer.y:126
		{ DoubiVAL.expr = &ast.BinaryExpr{DoubiS[Doubipt-2].expr, DoubiS[Doubipt-1].tok.Pos, token.AND, DoubiS[Doubipt-0].expr } }
	case 25:
		//line grammer.y:127
		{ DoubiVAL.expr = &ast.BinaryExpr{DoubiS[Doubipt-2].expr, DoubiS[Doubipt-1].tok.Pos, token.OR, DoubiS[Doubipt-0].expr } }
	case 26:
		//line grammer.y:128
		{ DoubiVAL.expr = &ast.BinaryExpr{DoubiS[Doubipt-2].expr, DoubiS[Doubipt-1].tok.Pos, token.XOR, DoubiS[Doubipt-0].expr } }
	case 27:
		//line grammer.y:129
		{ DoubiVAL.expr = &ast.BinaryExpr{DoubiS[Doubipt-2].expr, DoubiS[Doubipt-1].tok.Pos, token.SHL, DoubiS[Doubipt-0].expr } }
	case 28:
		//line grammer.y:130
		{ DoubiVAL.expr = &ast.BinaryExpr{DoubiS[Doubipt-2].expr, DoubiS[Doubipt-1].tok.Pos, token.SHR, DoubiS[Doubipt-0].expr } }
	case 29:
		//line grammer.y:131
		{ DoubiVAL.expr = &ast.BinaryExpr{DoubiS[Doubipt-2].expr, DoubiS[Doubipt-1].tok.Pos, token.AND_NOT, DoubiS[Doubipt-0].expr } }
	case 30:
		//line grammer.y:132
		{ DoubiVAL.expr = &ast.BinaryExpr{DoubiS[Doubipt-2].expr, DoubiS[Doubipt-1].tok.Pos, token.LSS, DoubiS[Doubipt-0].expr } }
	case 31:
		//line grammer.y:133
		{ DoubiVAL.expr = &ast.BinaryExpr{DoubiS[Doubipt-2].expr, DoubiS[Doubipt-1].tok.Pos, token.GTR, DoubiS[Doubipt-0].expr } }
	case 32:
		//line grammer.y:134
		{ DoubiVAL.expr = &ast.BinaryExpr{DoubiS[Doubipt-2].expr, DoubiS[Doubipt-1].tok.Pos, token.NEQ, DoubiS[Doubipt-0].expr } }
	case 33:
		//line grammer.y:135
		{ DoubiVAL.expr = &ast.BinaryExpr{DoubiS[Doubipt-2].expr, DoubiS[Doubipt-1].tok.Pos, token.LEQ, DoubiS[Doubipt-0].expr } }
	case 34:
		//line grammer.y:136
		{ DoubiVAL.expr = &ast.BinaryExpr{DoubiS[Doubipt-2].expr, DoubiS[Doubipt-1].tok.Pos, token.GEQ, DoubiS[Doubipt-0].expr } }
	case 35:
		//line grammer.y:137
		{ DoubiVAL.expr = &ast.BinaryExpr{DoubiS[Doubipt-2].expr, DoubiS[Doubipt-1].tok.Pos, token.EQL, DoubiS[Doubipt-0].expr } }
	case 36:
		//line grammer.y:139
		{ DoubiVAL.expr = &ast.BinaryExpr{DoubiS[Doubipt-2].expr, DoubiS[Doubipt-1].tok.Pos, token.LAND, DoubiS[Doubipt-0].expr } }
	case 37:
		//line grammer.y:140
		{ DoubiVAL.expr = &ast.BinaryExpr{DoubiS[Doubipt-2].expr, DoubiS[Doubipt-1].tok.Pos, token.LOR, DoubiS[Doubipt-0].expr } }
	case 38:
		//line grammer.y:143
		{ DoubiVAL.expr = &ast.ArrayExpr{DoubiS[Doubipt-2].tok.Pos, DoubiS[Doubipt-1].expr_list, DoubiS[Doubipt-0].tok.Pos} }
	case 39:
		//line grammer.y:145
		{ DoubiVAL.expr = &ast.ArrayExpr{DoubiS[Doubipt-4].tok.Pos, DoubiS[Doubipt-2].expr_list, DoubiS[Doubipt-1].tok.Pos} }
	case 40:
		//line grammer.y:147
		{ DoubiVAL.expr = &ast.ArrayExpr{DoubiS[Doubipt-3].tok.Pos, DoubiS[Doubipt-1].expr_list, DoubiS[Doubipt-0].tok.Pos} }
	case 41:
		//line grammer.y:150
		{ DoubiVAL.expr = &ast.SetExpr{DoubiS[Doubipt-2].tok.Pos, DoubiS[Doubipt-1].expr_list, DoubiS[Doubipt-0].tok.Pos} }
	case 42:
		//line grammer.y:152
		{ DoubiVAL.expr = &ast.SetExpr{DoubiS[Doubipt-4].tok.Pos, DoubiS[Doubipt-2].expr_list, DoubiS[Doubipt-0].tok.Pos} }
	case 43:
		//line grammer.y:154
		{ DoubiVAL.expr = &ast.SetExpr{DoubiS[Doubipt-3].tok.Pos, DoubiS[Doubipt-1].expr_list, DoubiS[Doubipt-0].tok.Pos} }
	case 44:
		//line grammer.y:157
		{ DoubiVAL.field = &ast.Field{DoubiS[Doubipt-2].expr, DoubiS[Doubipt-1].tok.Pos, DoubiS[Doubipt-0].expr} }
	case 45:
		//line grammer.y:159
		{ DoubiVAL.field_list = []*ast.Field{} }
	case 46:
		//line grammer.y:160
		{ DoubiVAL.field_list = []*ast.Field{DoubiS[Doubipt-0].field} }
	case 47:
		//line grammer.y:161
		{ DoubiVAL.field_list = append(DoubiS[Doubipt-2].field_list, DoubiS[Doubipt-0].field) }
	case 48:
		//line grammer.y:162
		{ DoubiVAL.field_list = append(DoubiS[Doubipt-2].field_list, DoubiS[Doubipt-0].field) }
	case 49:
		//line grammer.y:163
		{ DoubiVAL.field_list = append(DoubiS[Doubipt-3].field_list, DoubiS[Doubipt-0].field) }
	case 50:
		//line grammer.y:164
		{ DoubiVAL.field_list = DoubiS[Doubipt-1].field_list }
	case 51:
		//line grammer.y:165
		{ DoubiVAL.field_list = DoubiS[Doubipt-2].field_list }
	case 52:
		//line grammer.y:168
		{ DoubiVAL.expr = &ast.DictExpr{DoubiS[Doubipt-2].tok.Pos, DoubiS[Doubipt-1].field_list, DoubiS[Doubipt-0].tok.Pos} }
	case 53:
		//line grammer.y:171
		{ DoubiVAL.ident_list = []*ast.Ident{} }
	case 54:
		//line grammer.y:173
		{ DoubiVAL.ident_list = []*ast.Ident{&ast.Ident{DoubiS[Doubipt-0].tok.Pos, DoubiS[Doubipt-0].tok.Lit}} }
	case 55:
		//line grammer.y:175
		{ DoubiVAL.ident_list = append(DoubiS[Doubipt-2].ident_list, &ast.Ident{DoubiS[Doubipt-0].tok.Pos, DoubiS[Doubipt-0].tok.Lit}) }
	case 56:
		//line grammer.y:178
		{ DoubiVAL.expr = &ast.FuncDeclExpr{DoubiS[Doubipt-4].tok.Pos, nil, nil, nil, DoubiS[Doubipt-2].ident_list, DoubiS[Doubipt-0].stmt.(*ast.BlockStmt), []string{}} }
	case 57:
		//line grammer.y:180
		{ DoubiVAL.expr = &ast.FuncDeclExpr{DoubiS[Doubipt-5].tok.Pos, nil, nil, &ast.Ident{DoubiS[Doubipt-4].tok.Pos, DoubiS[Doubipt-4].tok.Lit}, DoubiS[Doubipt-2].ident_list, DoubiS[Doubipt-0].stmt.(*ast.BlockStmt), []string{}} }
	case 58:
		//line grammer.y:182
		{ DoubiVAL.expr = &ast.FuncDeclExpr{DoubiS[Doubipt-9].tok.Pos, &ast.Ident{DoubiS[Doubipt-7].tok.Pos, DoubiS[Doubipt-7].tok.Lit}, &ast.Ident{DoubiS[Doubipt-6].tok.Pos, DoubiS[Doubipt-6].tok.Lit},
	                                          &ast.Ident{DoubiS[Doubipt-4].tok.Pos, DoubiS[Doubipt-4].tok.Lit}, DoubiS[Doubipt-2].ident_list, DoubiS[Doubipt-0].stmt.(*ast.BlockStmt), []string{}} }
	case 59:
		DoubiVAL.expr = DoubiS[Doubipt-0].expr
	case 60:
		DoubiVAL.expr = DoubiS[Doubipt-0].expr
	case 61:
		DoubiVAL.expr = DoubiS[Doubipt-0].expr
	case 62:
		DoubiVAL.expr = DoubiS[Doubipt-0].expr
	case 63:
		DoubiVAL.expr = DoubiS[Doubipt-0].expr
	case 64:
		DoubiVAL.expr = DoubiS[Doubipt-0].expr
	case 65:
		DoubiVAL.expr = DoubiS[Doubipt-0].expr
	case 66:
		DoubiVAL.expr = DoubiS[Doubipt-0].expr
	case 67:
		DoubiVAL.expr = DoubiS[Doubipt-0].expr
	case 68:
		DoubiVAL.expr = DoubiS[Doubipt-0].expr
	case 69:
		DoubiVAL.expr = DoubiS[Doubipt-0].expr
	case 70:
		DoubiVAL.expr = DoubiS[Doubipt-0].expr
	case 71:
		DoubiVAL.expr = DoubiS[Doubipt-0].expr
	case 72:
		//line grammer.y:201
		{ DoubiVAL.stmt = &ast.ExprStmt{DoubiS[Doubipt-0].expr} }
	case 73:
		//line grammer.y:203
		{ DoubiVAL.stmt = &ast.SendStmt{DoubiS[Doubipt-2].expr, DoubiS[Doubipt-1].tok.Pos, DoubiS[Doubipt-0].expr} }
	case 74:
		//line grammer.y:205
		{ DoubiVAL.stmt = &ast.IncDecStmt{DoubiS[Doubipt-1].expr, DoubiS[Doubipt-0].tok.Pos, token.INC} }
	case 75:
		//line grammer.y:206
		{ DoubiVAL.stmt = &ast.IncDecStmt{DoubiS[Doubipt-1].expr, DoubiS[Doubipt-0].tok.Pos, token.DEC} }
	case 76:
		//line grammer.y:208
		{ DoubiVAL.stmt = &ast.AssignStmt{DoubiS[Doubipt-2].expr_list, DoubiS[Doubipt-1].tok.Pos, token.ASSIGN, DoubiS[Doubipt-0].expr_list} }
	case 77:
		//line grammer.y:209
		{ DoubiVAL.stmt = &ast.AssignStmt{DoubiS[Doubipt-2].expr_list, DoubiS[Doubipt-1].tok.Pos, token.ADD_ASSIGN, DoubiS[Doubipt-0].expr_list} }
	case 78:
		//line grammer.y:210
		{ DoubiVAL.stmt = &ast.AssignStmt{DoubiS[Doubipt-2].expr_list, DoubiS[Doubipt-1].tok.Pos, token.SUB_ASSIGN, DoubiS[Doubipt-0].expr_list} }
	case 79:
		//line grammer.y:211
		{ DoubiVAL.stmt = &ast.AssignStmt{DoubiS[Doubipt-2].expr_list, DoubiS[Doubipt-1].tok.Pos, token.MUL_ASSIGN, DoubiS[Doubipt-0].expr_list} }
	case 80:
		//line grammer.y:212
		{ DoubiVAL.stmt = &ast.AssignStmt{DoubiS[Doubipt-2].expr_list, DoubiS[Doubipt-1].tok.Pos, token.QUO_ASSIGN, DoubiS[Doubipt-0].expr_list} }
	case 81:
		//line grammer.y:213
		{ DoubiVAL.stmt = &ast.AssignStmt{DoubiS[Doubipt-2].expr_list, DoubiS[Doubipt-1].tok.Pos, token.REM_ASSIGN, DoubiS[Doubipt-0].expr_list} }
	case 82:
		//line grammer.y:214
		{ DoubiVAL.stmt = &ast.AssignStmt{DoubiS[Doubipt-2].expr_list, DoubiS[Doubipt-1].tok.Pos, token.AND_ASSIGN, DoubiS[Doubipt-0].expr_list} }
	case 83:
		//line grammer.y:215
		{ DoubiVAL.stmt = &ast.AssignStmt{DoubiS[Doubipt-2].expr_list, DoubiS[Doubipt-1].tok.Pos, token.OR_ASSIGN, DoubiS[Doubipt-0].expr_list} }
	case 84:
		//line grammer.y:216
		{ DoubiVAL.stmt = &ast.AssignStmt{DoubiS[Doubipt-2].expr_list, DoubiS[Doubipt-1].tok.Pos, token.XOR_ASSIGN, DoubiS[Doubipt-0].expr_list} }
	case 85:
		//line grammer.y:217
		{ DoubiVAL.stmt = &ast.AssignStmt{DoubiS[Doubipt-2].expr_list, DoubiS[Doubipt-1].tok.Pos, token.SHL_ASSIGN, DoubiS[Doubipt-0].expr_list} }
	case 86:
		//line grammer.y:218
		{ DoubiVAL.stmt = &ast.AssignStmt{DoubiS[Doubipt-2].expr_list, DoubiS[Doubipt-1].tok.Pos, token.SHR_ASSIGN, DoubiS[Doubipt-0].expr_list} }
	case 87:
		//line grammer.y:219
		{ DoubiVAL.stmt = &ast.AssignStmt{DoubiS[Doubipt-2].expr_list, DoubiS[Doubipt-1].tok.Pos, token.AND_NOT_ASSIGN, DoubiS[Doubipt-0].expr_list} }
	case 88:
		//line grammer.y:222
		{ DoubiVAL.stmt = &ast.GoStmt{DoubiS[Doubipt-1].tok.Pos, DoubiS[Doubipt-0].expr.(*ast.CallExpr)} }
	case 89:
		//line grammer.y:225
		{ DoubiVAL.stmt = &ast.ReturnStmt{DoubiS[Doubipt-1].tok.Pos, DoubiS[Doubipt-0].expr_list} }
	case 90:
		//line grammer.y:227
		{ DoubiVAL.stmt = &ast.BranchStmt{DoubiS[Doubipt-0].tok.Pos, token.BREAK} }
	case 91:
		//line grammer.y:228
		{ DoubiVAL.stmt = &ast.BranchStmt{DoubiS[Doubipt-0].tok.Pos, token.CONTINUE } }
	case 92:
		//line grammer.y:230
		{ DoubiVAL.stmt = &ast.BlockStmt{DoubiS[Doubipt-2].tok.Pos, DoubiS[Doubipt-1].stmt_list ,DoubiS[Doubipt-0].tok.Pos} }
	case 93:
		//line grammer.y:232
		{ DoubiVAL.stmt = &ast.IfStmt{DoubiS[Doubipt-2].tok.Pos, DoubiS[Doubipt-1].expr, DoubiS[Doubipt-0].stmt.(*ast.BlockStmt), nil} }
	case 94:
		//line grammer.y:233
		{ DoubiVAL.stmt = &ast.IfStmt{DoubiS[Doubipt-4].tok.Pos, DoubiS[Doubipt-3].expr, DoubiS[Doubipt-2].stmt.(*ast.BlockStmt), DoubiS[Doubipt-0].stmt} }
	case 95:
		//line grammer.y:235
		{ DoubiVAL.stmt = &ast.CaseClause{DoubiS[Doubipt-3].tok.Pos, DoubiS[Doubipt-2].expr_list, DoubiS[Doubipt-1].tok.Pos, DoubiS[Doubipt-0].stmt_list} }
	case 96:
		//line grammer.y:236
		{ DoubiVAL.stmt = &ast.CaseClause{DoubiS[Doubipt-2].tok.Pos, nil, DoubiS[Doubipt-1].tok.Pos, DoubiS[Doubipt-0].stmt_list} }
	case 97:
		//line grammer.y:238
		{ DoubiVAL.stmt_list = []ast.Stmt{} }
	case 98:
		//line grammer.y:239
		{ DoubiVAL.stmt_list = []ast.Stmt{DoubiS[Doubipt-0].stmt} }
	case 99:
		//line grammer.y:240
		{ DoubiVAL.stmt_list = append(DoubiS[Doubipt-1].stmt_list, DoubiS[Doubipt-0].stmt) }
	case 100:
		//line grammer.y:242
		{ DoubiVAL.stmt = &ast.BlockStmt{DoubiS[Doubipt-2].tok.Pos, DoubiS[Doubipt-1].stmt_list, DoubiS[Doubipt-0].tok.Pos} }
	case 101:
		//line grammer.y:244
		{ DoubiVAL.stmt = &ast.SwitchStmt{DoubiS[Doubipt-2].tok.Pos, DoubiS[Doubipt-1].stmt, DoubiS[Doubipt-0].stmt.(*ast.BlockStmt)} }
	case 102:
		//line grammer.y:246
		{ DoubiVAL.stmt = &ast.SelectStmt{DoubiS[Doubipt-1].tok.Pos, DoubiS[Doubipt-0].stmt.(*ast.BlockStmt)} }
	case 103:
		//line grammer.y:249
		{ DoubiVAL.stmt = &ast.ForStmt{DoubiS[Doubipt-6].tok.Pos, DoubiS[Doubipt-5].stmt, DoubiS[Doubipt-3].expr, DoubiS[Doubipt-1].stmt, DoubiS[Doubipt-0].stmt.(*ast.BlockStmt)} }
	case 104:
		//line grammer.y:251
		{ DoubiVAL.stmt = &ast.ForStmt{DoubiS[Doubipt-5].tok.Pos, nil, DoubiS[Doubipt-3].expr, DoubiS[Doubipt-1].stmt, DoubiS[Doubipt-0].stmt.(*ast.BlockStmt)} }
	case 105:
		//line grammer.y:253
		{ DoubiVAL.stmt = &ast.ForStmt{DoubiS[Doubipt-2].tok.Pos, nil, DoubiS[Doubipt-1].expr, nil, DoubiS[Doubipt-0].stmt.(*ast.BlockStmt)} }
	case 106:
		//line grammer.y:256
		{ DoubiVAL.stmt = &ast.RangeStmt{DoubiS[Doubipt-5].tok.Pos, DoubiS[Doubipt-4].expr_list, DoubiS[Doubipt-1].expr, DoubiS[Doubipt-0].stmt.(*ast.BlockStmt)} }
	case 107:
		//line grammer.y:259
		{ DoubiVAL.stmt = &ast.ImportStmt{DoubiS[Doubipt-1].tok.Pos, []string{DoubiS[Doubipt-0].tok.Lit}} }
	case 108:
		DoubiVAL.stmt = DoubiS[Doubipt-0].stmt
	case 109:
		DoubiVAL.stmt = DoubiS[Doubipt-0].stmt
	case 110:
		DoubiVAL.stmt = DoubiS[Doubipt-0].stmt
	case 111:
		DoubiVAL.stmt = DoubiS[Doubipt-0].stmt
	case 112:
		DoubiVAL.stmt = DoubiS[Doubipt-0].stmt
	case 113:
		DoubiVAL.stmt = DoubiS[Doubipt-0].stmt
	case 114:
		DoubiVAL.stmt = DoubiS[Doubipt-0].stmt
	case 115:
		DoubiVAL.stmt = DoubiS[Doubipt-0].stmt
	case 116:
		DoubiVAL.stmt = DoubiS[Doubipt-0].stmt
	case 117:
		DoubiVAL.stmt = DoubiS[Doubipt-0].stmt
	case 118:
		DoubiVAL.stmt = DoubiS[Doubipt-0].stmt
	case 119:
		DoubiVAL.stmt = DoubiS[Doubipt-0].stmt
	case 120:
		DoubiVAL.stmt = DoubiS[Doubipt-0].stmt
	case 121:
		DoubiVAL.stmt = DoubiS[Doubipt-0].stmt
	case 122:
		//line grammer.y:276
		{ DoubiVAL.stmt_list = []ast.Stmt{} }
	case 123:
		//line grammer.y:277
		{ DoubiVAL.stmt_list = []ast.Stmt{DoubiS[Doubipt-0].stmt} }
	case 124:
		//line grammer.y:278
		{ DoubiVAL.stmt_list = append(DoubiS[Doubipt-2].stmt_list, DoubiS[Doubipt-0].stmt) }
	case 125:
		//line grammer.y:279
		{ DoubiVAL.stmt_list = append(DoubiS[Doubipt-2].stmt_list, DoubiS[Doubipt-0].stmt) }
	case 126:
		//line grammer.y:280
		{ DoubiVAL.stmt_list = DoubiS[Doubipt-1].stmt_list }
	case 127:
		//line grammer.y:285
		{ ProgramAst = DoubiS[Doubipt-1].stmt_list }
	}
	goto Doubistack /* stack new state and value */
}
