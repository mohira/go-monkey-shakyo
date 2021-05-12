package lexer

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
