package evaluator

import (
	"go-monkey-shakyo/monkey/ast"
	"go-monkey-shakyo/monkey/object"
)

var (
	NULL  = &object.NULL{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {

	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.BlockStatement:
		return evalStatements(node.Statements)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue)
		return &object.ReturnValue{Value: val}

	case *ast.IfExpression:
		return evalIfExpression(node)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.PrefixExpression:
		// (!true) → ! , true
		// (!(!true)) → ! , !true
		right := Eval(node.Right)

		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)

		return evalInfixExpression(node.Operator, left, right)

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

		// 直近の評価結果が object.ReturnValue かどうかを確認し、
		// もしそうならば評価を中断し、アンラップした値を返す。
		if returnValue, ok := result.(*object.ReturnValue); ok {
			return returnValue.Value
		}
	}
	return result
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		// エラーの処理の仕組みを実装していないので、今はNULLにしておく
		return NULL
	}
}

// !5 や !true みたいな式を評価する
func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		// 整数はここの条件に該当する
		// !Integer は false
		return FALSE
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return NULL
	}
	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

// 中置式を評価する
func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	default:
		return NULL
	}
}

// 整数の中置演算式を評価する
func evalIntegerInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return NULL
	}
}

func evalIfExpression(ie *ast.IfExpression) object.Object {
	// if (<condition>) { <consequence> } else { <alternative> }
	condition := Eval(ie.Condition)
	if isTruthy(condition) {
		return Eval(ie.Consequence)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative)
	} else {
		// 条件分岐を評価した結果が何かの値にならなかった場合は NULL を返す
		return NULL
	}
}

// nullでもfalseでもなければ、それはtruthy
func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}
