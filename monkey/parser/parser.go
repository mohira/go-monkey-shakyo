package parser

import (
	"fmt"
	"go-monkey-shakyo/monkey/ast"
	"go-monkey-shakyo/monkey/lexer"
	"go-monkey-shakyo/monkey/token"
	"strconv"
)

// Monkey言語における優先順位の定義
const (
	// iotaは0の値を取り、続く定数には 1 から 7 の値が割り振られる
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > または <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X または !X
	CALL        // myFunction(X)
)

// 演算子の優先順位テーブル
var precedences = map[token.TokenType]int{
	token.EQ:     EQUALS,
	token.NOT_EQ: EQUALS,

	token.LT: LESSGREATER,
	token.GT: LESSGREATER,

	token.PLUS:  SUM,
	token.MINUS: SUM,

	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
}

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

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)

	// 前置演算子の解析用関数の登録
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)

	// 中置演算子の解析用関数の登録
	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)

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
		// Monkeyにおける純粋な文は2種類で、let文とreturn文しか存在しない。
		// もしそれ以外のものが出現したら式文の構文解析を試みることにしよう
		return p.parseExpressionStatement()
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

func (p *Parser) parseExpressionStatement() ast.Statement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	// p.58
	// 省略可能なセミコロンをチェックする。
	// もしセミコロンがなかったとしても問題はない。
	// 構文解析器にエラーを追加するようなことはしない。
	// なぜかというと、式文のセミコロンを省略できるようにしたいからだ
	// （こうしておけば後ほど5 + 5のようなものをREPLに入力しやすくなる）。
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// p.curToken.Typeの前置に関連付けられた構文解析関数を調べて、存在するなら呼び出す
func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}

	leftExp := prefix()

	// `precedence < p.peekPrecedence()` は"結合力"のチェックをしている
	//  つまり、2つの演算子(op1, op2 とする)の優先度を比べて、i)かii)の判断をする感じ
	// 		 i) op1を採用して、  Left-op1-Right               をExpressionとする
	// 		ii) op1は無視して、           Right-op2-NextRight をExpressionとする
	//
	//   i) 左結合力が強い → 左に吸い込んで1つの式にする
	//          Left op1 Right op2 NextRight
	//	   式文: 1     *    2    +      3
	//			    ←←←←←        →
	//	   結果: (1 * 2)         +      3
	//
	//  ii) 右結合力が強い → 右に吸い込んで1つの式にする
	//          Left op1 Right op2 NextRight
	//	   式文: 1     +    2    *      3
	//			      ←       →→→→→
	//	   結果: 1     +           (2 * 3)
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	// 整数リテラルの文字列をint64に変換する
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value

	return lit
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}

	return LOWEST
}

// left という Expression を受け取っているのが、 parsePrefixExpression との重要な違い
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()

	if expression.Operator == "+" {
		expression.Right = p.parseExpression(precedence - 1)
	} else {
		expression.Right = p.parseExpression(precedence)

	}

	return expression
}
