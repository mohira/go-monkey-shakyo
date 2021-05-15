package parser

import (
	"fmt"
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

	errors []string

	curToken  token.Token // 現在のトークンを指し示す(※文字じゃないよ！)
	peekToken token.Token // 次のトークンを指し示す(※文字じゃないよ！)

	// トークンタイプごとに適切な構文解析関数を持てるようにする
	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

type (
	// 前置構文を解析するための関数
	// 関連付けられたトークンタイプが前置で出現した場合に呼ばれる
	prefixParseFn func() ast.Expression

	// 中置構文を解析するための関数
	// 引数は、中置演算子の「左側」
	// 関連付けられたトークンタイプが中置で出現した場合に呼ばれる
	infixParseFn func(ast.Expression) ast.Expression
)

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// 2つトークンを読み込む。
	// curToken と peekToken の両方がセットされる
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
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

// プログラムが読んでいるトークンにあわせて構文を解析していく
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return nil
	}
}

func (p *Parser) parseLetStatement() ast.Statement {
	// let文は
	// 		let <identifier> = <expression>;
	// という構造なので、 let → Identifier → ASSIGN → Expression → SEMICOLON と期待していく
	stmt := &ast.LetStatement{Token: p.curToken}

	//
	if !p.expectPeek(token.IDENT) {
		return nil
	}

	// p.expectPeek() の次のトークンに進めているので、ここの curToken は 識別子 が入っている
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO: セミコロンに遭遇するまで式を読み飛ばしてしまっている
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// 構文解析器のよくある「アサーション関数」らしい
// 次のトークンの型が期待されるものだったときだけトークンを進める
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		// トークンが期待されるものでなかったら、デバッグ情報をもたせる
		p.peekError(t)
		return false
	}
}

func (p *Parser) parseReturnStatement() ast.Statement {
	// return文
	// 	return <expression>;
	// という構造なので、 RETURN → EXPRESSION → SEMICOLON と期待していく感じ
	stmt := &ast.ReturnStatement{Token: p.curToken}

	// RETURN の次のトークンを読む
	p.nextToken()

	// TODO: セミコロンに遭遇するまで式を読み飛ばしてしまっている
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// 前置構文解析関数をマップに追加するためのヘルパーメソッド
func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

// 中置構文解析関数をマップに追加するためのヘルパーメソッド
func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}
