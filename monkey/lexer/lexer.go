package lexer

import (
	"go-monkey-shakyo/monkey/token"
)

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

	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.EQ, Literal: literal}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case ':':
		tok = newToken(token.COLON, l.ch)
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
	case '[':
		tok = newToken(token.LBRACEKT, l.ch)
	case ']':
		tok = newToken(token.RBRACEKT, l.ch)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.NOT_EQ, Literal: literal}
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF

	default:
		// l.chが認識された文字でないときに「識別子」かどうかを点検する

		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)

			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()

			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	// トークンを返す前に次の文字に進める
	l.readChar()

	return tok
}

// 先読み用のヘルパー関数
// 覗き見(peek)するだけで、字句解析の文字位置は進めない
// MEMO: 言語におけるパースの難易度の違いは、ソースコードを解釈する際に、
//       どの程度先まで読む（もしくは戻って読む！）必要があるかによるところが大きい。
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// 整数しか読めない
// 浮動小数点数も16進数も扱えない

func (l *Lexer) readNumber() string {
	position := l.position

	// 数字である限り読みすすめる
	for isDigit(l.ch) {
		l.readChar()
	}

	return l.input[position:l.position]
}

func (l *Lexer) readIdentifier() string {
	position := l.position

	// 英字を区切りまで読み進める
	for isLetter(l.ch) {
		l.readChar()
	}

	return l.input[position:l.position]
}

// どういう文字を読み飛ばすかを決める
// 後の構文解析ステップを少し簡単にするため、改行文字を単に読み飛ばすことにする。
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// 英字かどうかの判定をする
// 重要: '_' も英字として扱う ⇔ 識別子とキーワードに '_' が含まれることを許容する！
// この関数によって、何が許されるかを決められるので重要(例えば、 ? を許可することもできるわけで)
func isLetter(ch byte) bool {
	isAlphabet := ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z')

	return ch == '_' || isAlphabet
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

// 閉じ二重引用符か入力の最後に至るまで readChar を呼ぶ。
func (l *Lexer) readString() string {
	position := l.position + 1

	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}

	return l.input[position:l.position]
}
