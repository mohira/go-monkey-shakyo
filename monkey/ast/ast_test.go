package ast

import (
	"go-monkey-shakyo/monkey/token"
	"testing"
)

func TestString(t *testing.T) {
	// 次のようなソースコードを解析していたときに、
	// 		let myVar = anotherVar;
	// そのままの文字列を復元できる
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}

	if program.String() != "let myVar = anotherVar " {
		t.Errorf("program.String() wrong. got=%q", program.String())
	}
}
