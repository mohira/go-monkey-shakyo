package evaluator

import (
	"go-monkey-shakyo/monkey/ast"
	"go-monkey-shakyo/monkey/object"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {

	// 文の場合
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)

	// 式の場合
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	}
	return nil
}

func nativeBoolToBooleanObject(input bool) object.Object {
	if input == true {
		return TRUE
	} else {
		return FALSE
	}
}

func evalStatements(statements []ast.Statement) object.Object {
	var result object.Object

	for _, statement := range statements {
		result = Eval(statement)
	}
	return result
}
