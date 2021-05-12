package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL" // トークンが未知であることを示す
	EOF     = "EOF"     // ファイル終端

	// 識別子(Identifier) + リテラル
	IDENT = "IDENT" // add, foobar, x, y, ...
	INT   = "INT"   // 1343456

	// 演算子
	ASSIGN = "="
	PLUS   = "+"

	// デリミタ
	COMMA     = ","
	SEMICOLON = ";"

	// かっこ
	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// キーワード: 識別子のようにみえて識別子ではなく、実際の言語の一部である単語
	FUNCTION = "FUNCTION"
	LET      = "LET"
)
