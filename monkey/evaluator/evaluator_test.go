package evaluator

import (
	"go-monkey-shakyo/monkey/lexer"
	"go-monkey-shakyo/monkey/object"
	"go-monkey-shakyo/monkey/parser"
	"testing"
)

// 整数リテラルを含む式文が与えられたときの評価
func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int64
	}{
		{"5 を評価したら 5 になる", "5", 5},
		{"10 を評価したら 10 になる", "10", 10},
		{"「-」前置演算子", "-5", -5},
		{"「-」前置演算子", "-10", -10},

		{"中置式", "5 + 5 + 5 + 5 - 10", 10},
		{"中置式", "2 * 2 * 2 * 2 * 2", 32},
		{"中置式", "-50 + 100 + -50", 0},
		{"中置式", "5 * 2 + 10", 20},
		{"中置式", "5 + 2 * 10", 25},
		{"中置式", "20 + 2 * -10", 0},
		{"中置式", "50 / 2 * 2 + 10", 60},
		{"中置式", "2 * (5 + 10)", 30},
		{"中置式", "3 * 3 * 3 + 10", 37},
		{"中置式", "3 * (3 * 3) + 10", 37},
		{"中置式", "(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluated := testEval(tt.input)

			testIntegerObject(t, evaluated, tt.expected)
		})
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	return Eval(program)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
		return false
	}

	return true
}

// Booleanに関連する式を評価するテスト
func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"true を 評価すると true", "true", true},
		{"false を 評価すると false", "false", false},

		{"真偽値を返す整数同士の中置演算", "1 < 2", true},
		{"真偽値を返す整数同士の中置演算", "1 > 2", false},
		{"真偽値を返す整数同士の中置演算", "1 < 1", false},
		{"真偽値を返す整数同士の中置演算", "1 > 1", false},
		{"真偽値を返す整数同士の中置演算", "1 == 1", true},
		{"真偽値を返す整数同士の中置演算", "1 != 1", false},
		{"真偽値を返す整数同士の中置演算", "1 == 2", false},
		{"真偽値を返す整数同士の中置演算", "1 != 2", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluated := testEval(tt.input)

			testBooleanObject(t, evaluated, tt.expected)
		})
	}
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Fatalf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Fatalf("object has wrong value. got=%t, want=%t", result.Value, expected)
		return false
	}

	return true
}

// 「!」演算子は、オペランドを真偽値に変換して、その否定を返す
func TestBangOperator(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"", "!true", false},
		{"", "!false", true},
		{"!5 は false (5はtruthy)", "!5", false},

		{"2回適用する場合", "!!true", true},
		{"2回適用する場合", "!!false", false},
		{"2回適用する場合", "!!5", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluated := testEval(tt.input)

			testBooleanObject(t, evaluated, tt.expected)
		})
	}
}
