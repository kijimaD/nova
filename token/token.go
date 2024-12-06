package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// 識別子 + リテラル。数値や変数名など、予約語ではないもの。
	STRING = "STRING"
	IDENT  = "IDENT"
	TEXT   = "TEXT"

	LBRACKET = "["
	RBRACKET = "]"
	COMMA    = ","
	EQUAL    = "="
	ASTERISK = "*"
	NEWLINE  = "\n"

	CMD_FLUSH         = "p"
	CMD_LINE_END_WAIT = "l"
	CMD_NEWLINE       = "r"
	CMD_IMAGE         = "image"
	CMD_WAIT          = "wait"
	CMD_JUMP          = "jump"
)

// 予約語
var keywords = map[string]TokenType{}

// 予約語の場合はその種類を、それ以外の場合はIDENTを返す
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
