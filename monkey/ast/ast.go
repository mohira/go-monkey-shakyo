package ast

import "go-monkey-shakyo/monkey/token"

type Node interface {
	// そのノードが関連づけられているトークンのリテラル値を返す
	// デバッグのみに用いる
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

// let <identifier> = <expression>
type LetStatement struct {
	Token token.Token // token.LET トークン
	Name  *Identifier // 束縛した識別子を保持するため
	Value Expression  // 値を生成する式を保持するため
}

// LetStatement は Statement インタフェース を実装している
// (同時に Node インタフェース を実装している)
func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

type Identifier struct {
	Token token.Token // token.IDENT トークン
	Value string
}

// Identifier は Expression インタフェース を実装している
// (同時に Node インタフェース を実装している)
func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

type ReturnStatement struct {
	Token       token.Token // 'return' トークン
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}
func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es ExpressionStatement) statementNode() {}

func (es ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}
