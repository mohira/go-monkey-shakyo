package lexer

import "go-monkey-shakyo/monkey/token"

type Lexer struct {
	input        string
	position     int  // 入力における現在の位置(現在の文字を指し示す)
	readPosition int  // これから読み込む位置(現在の文字の次)
	ch           byte // 現在検査中の文字
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// 次の1文字を読んでinput文字列の現在位置をすすめる
// 注意: ASCII文字対応のみでUnicodeには非対応。バイト列の解析が必要になるからね(詳しくはp.7参照))。
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		// 入力の終端チェック(読み切った場合)

		// byte(0) は ASCIIコードの "NUL"文字
		// https://play.golang.org/p/6NnGcUgwNBt
		l.ch = 0
	} else {
		// 次の文字を読み込む
		l.ch = l.input[l.readPosition]
	}

	// 検査対象の文字の位置を進める
	// lo.positionは常に最後に読んだ場所を指し示す
	l.position = l.readPosition
	l.readPosition += 1
}

// 現在検査中の文字l.chを見て、その文字が何であるかに応じてトークンを返す
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	switch l.ch {
	case '=':
		tok = newToken(token.ASSIGN, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	}

	// トークンを返す前に次の文字に進める
	l.readChar()

	return tok
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}
