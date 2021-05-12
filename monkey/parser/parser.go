package parser

import (
	"go-monkey-shakyo/monkey/ast"
	"go-monkey-shakyo/monkey/lexer"
	"go-monkey-shakyo/monkey/token"
)

// peekTokenが必要な理由
// p.35 例として、1つの行に 5; だけがある場合を考えてみよう。
//      ここで、curTokenはtoken.INTとなる。
//      このとき、行末にいるのか、算術式が始まったところなのかを判定するためにpeekTokenが必要だ。
type Parser struct {
	l *lexer.Lexer // パーサーは字句解析器を(のポインタ)もつ

	curToken  token.Token // 現在のトークンを指し示す(※文字じゃないよ！)
	peekToken token.Token // 次のトークンを指し示す(※文字じゃないよ！)
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}

	// 2つトークンを読み込む。
	// curToken と peekToken の両方がセットされる
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{} // ASTのルートノードの生成
	program.Statements = []ast.Statement{}

	// EOFになるまで「トークン」を読み続ける
	// 構文を解析しては、Statementとして溜め込んでいく
	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()

		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}
