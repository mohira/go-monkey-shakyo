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

	// 複数のテストケースを横断するような状態は持ち込むべきじゃないので、
	// testEval()のたびに「環境」を初期化
	env := object.NewEnvironment()

	return Eval(program, env)
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

		// Monkeyがサポートするのは、真偽値のオペランドに関しては等値演算子「==」と「!=」だけだ。
		// 真偽値の加算、減算、除算、乗算には対応しない。
		// trueがfalseより大きいかを「<」や「>」で比較するようなこともサポートしない。
		{"両方のオペランドが真偽値の場合の中置演算", "true == true", true},
		{"両方のオペランドが真偽値の場合の中置演算", "false == false", true},
		{"両方のオペランドが真偽値の場合の中置演算", "true == false", false},
		{"両方のオペランドが真偽値の場合の中置演算", "true != false", true},
		{"両方のオペランドが真偽値の場合の中置演算", "(1 < 2) == true", true},
		{"両方のオペランドが真偽値の場合の中置演算", "(1 < 2) == false", false},
		{"両方のオペランドが真偽値の場合の中置演算", "(1 > 2) == true", false},
		{"両方のオペランドが真偽値の場合の中置演算", "(1 > 2) == false", true},
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

// 条件分岐の文の評価
func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
	}{
		{"true は truthy", "if (true) { 10 }", 10},
		{"false は truthy ではない", "if (false) { 10 }", nil},

		{"整数 は truthy ", "if (1) { 10 }", 10},
		{"if文にマッチするやつ", "if (1 < 2) { 10 }", 10},

		{"条件分岐を評価した結果が何かの値にならなかった場合は NULL を返す", "if (1 > 2) { 10 }", nil},

		{"if-elseでif分岐にマッチ", "if (1 > 2) { 10 } else { 20 }", 20},
		{"if-elseでelse分岐にマッチ", "if (1 < 2) { 10 } else { 20 }", 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluated := testEval(tt.input)
			integer, ok := tt.expected.(int)
			if ok {
				testIntegerObject(t, evaluated, int64(integer))
			} else {
				testNullObject(t, evaluated)
			}
		})
	}
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}

	return true
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int64
	}{
		{"", "return 10;", 10},
		{"return文に続く文は評価に無関係", "return 10; 9;", 10},
		{"return <expression> の <expression>もちゃんと評価される", "return 2 * 5; 9;", 10},
		{"return文の前後の文は評価に無関係", "9; return 2 * 5; 9;;", 10},

		{"ネストしたブロック文を正しく評価できる",
			`
if (10 > 1) {
	if (10 > 1) {
		return 10;
	}

	return 1;
}
`,
			10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluated := testEval(tt.input)

			testIntegerObject(t, evaluated, tt.expected)
		})
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		expectedMessage string
	}{
		{
			"未定義の演算: 異なる型の演算はエラー: 整数 + 真偽値 はエラーである",
			"5 + true;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"複数の文において、エラーとなる演算が評価されたときに中断される",
			"5 + true; 5;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"未定義の演算: -boolean はエラーである",
			"-true",
			"unknown operator: -BOOLEAN",
		},
		{
			"未定義の演算: boolean + boolean はエラーである",
			"true + false;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"複数の文において、エラーとなる演算が評価されたときに中断される",
			"5; true + false; 5;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"未定義の演算: boolean + boolean はエラーである",
			"if (10 > 1) { true + false; }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"未定義の演算: boolean + boolean はエラーである",
			`
if (10 > 1) {
	if (10 > 1) {
		return true + false;
	}
	return 1;
}
`,
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"束縛されていない識別子を評価するとエラーになる",
			"foobar",
			"identifier not found: foobar",
		},
		{
			"未定義の演算: string - string はエラーである",
			`"Hello" - "World"`,
			"unknown operator: STRING - STRING",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluated := testEval(tt.input)

			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("no error object returned. got=%T(%+v)", evaluated, evaluated)
				return
			}

			if errObj.Message != tt.expectedMessage {
				t.Errorf("wrong error message. expected=%q, got=%q", tt.expectedMessage, errObj.Message)
			}

		})
	}
}

// let文において値を生成する式の評価と、名前に束縛された識別子の評価をしている
func TestLetStatements(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int64
	}{
		{"単一のlet文: 整数リテラルをそのまま代入", "let a = 5; a;", 5},
		{"単一のlet文: 値を生成する式を評価してからの代入", "let a = 5 * 5; a;", 25},
		{"複数のlet文: 束縛された識別子の評価", "let a = 5; let b = a; b;", 5},
		{"複数のlet文: 束縛された識別子を含む式の評価", "let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluated := testEval(tt.input)
			testIntegerObject(t, evaluated, tt.expected)
		})
	}
}

// 関数リテラルを評価したときに、正しいパラメータと本文を持った*object.Functionが返されるかのテスト
func TestFunctionObject(t *testing.T) {
	input := `fn(x) { x + 2; };`

	evaluated := testEval(input)

	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T(%+v)", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Paramters=%+v", fn.Parameters)
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Parameters[0])
	}

	expectedBody := "(x + 2)"
	if fn.Body.String() != expectedBody {
		t.Fatalf("body in not %q. got=%q", expectedBody, fn.Body.String())
	}
}

// 関数適用のテスト
func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int64
	}{
		{
			"return文のない関数の呼び出し(暗黙の戻り値)",
			"let identity = fn(x) { x; }; identity(5);",
			5,
		},
		{
			"return文のある関数の呼び出し",
			"let identity = fn(x) { return x; }; identity(5);",
			5,
		},
		{
			"関数本文の式の中でパラメータを利用する",
			"let double = fn(x) { x * 2; }; double(5)",
			10,
		},
		{
			"複数のパラメータを持つ関数の呼び出し",
			"let add = fn(x, y) { x + y; }; add(5, 5);",
			10,
		},
		{
			"関数に関数を渡せる",
			"let add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5));",
			20,
		},
		{
			"関数リテラルで即呼び出し",
			"fn(x) { x; }(5);",
			5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluated := testEval(tt.input)

			testIntegerObject(t, evaluated, tt.expected)
		})
	}
}

// 「なぜ、現在の環境ではなく、関数の環境を拡張するか？」の答えとしてのクロージャーのテスト
func TestClosures(t *testing.T) {
	input := `
let newAdder = fn(x) {
	fn(y) { x + y };
};

let addTwo = newAdder(2);
addTwo(2);
`

	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 4)
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello World!"`

	evaluated := testEval(input)

	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T", evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

// 文字列の結合
func TestStringConcatenation(t *testing.T) {
	input := `"Hello" + " " + "World!"`
	evaluated := testEval(input)

	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T(%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
	}{
		{
			"len(): 空文字は長さ0",
			`len("")`,
			0,
		},
		{
			"len(): 文字数を数える",
			`len("four")`,
			4,
		},
		{
			"len(): 空白も1文字と数える",
			`len("hello world")`,
			11,
		},
		{
			"len(): エラー: 整数を引数にはとれない",
			`len(1)`,
			"argument to `len` not supported, got INTEGER",
		},
		{
			"len(): エラー: 引数は1つでなければいけない",
			`len("one", "two")`,
			"wrong number of arguments. got=2, want=1",
		},
		{
			"len(): 配列の要素数を取得できる",
			`len([1, 2, 3])`,
			3,
		},
		{
			"len(): 空の配列の長さは0である",
			`len([])`,
			0,
		},
		{
			"first(): 配列の最初の要素を取得できる",
			`first([1, 2, 3])`,
			1,
		},
		{
			"first(): 空の配列の最初の要素はNULL",
			`first([])`,
			nil,
		},
		{
			"first(): エラー",
			`first(1)`,
			"argument to `first` must be ARRAY, got INTEGER",
		},
		{
			"last(): 配列の最後の要素を取得できる",
			`last([1, 2, 3])`,
			3,
		},
		{
			"last(): 空の配列の最後の要素はNULL",
			`last([])`,
			nil,
		},
		{
			"last(): エラー",
			`last(1)`,
			"argument to `last` must be ARRAY, got INTEGER",
		},
		{
			"rest(): cdrと同じ動き。与えられた配列の最初の1つを除いて残りを全て含む新しい配列を返す。",
			`rest([1, 2, 3])`,
			[]int{2, 3},
		},
		{
			"rest(): 空の配列のrestはNULL",
			`rest([])`,
			nil,
		},
		{
			"push(): 要素を追加した新しい配列を返す",
			`push([], 1)`,
			[]int{1},
		},
		{
			"push(): エラー",
			`push(1, 1)`,
			"argument to `push` must be ARRAY, got INTEGER",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluated := testEval(tt.input)

			switch expected := tt.expected.(type) {
			case nil:
				testNullObject(t, evaluated)
			case int:
				testIntegerObject(t, evaluated, int64(expected))
			case string:
				errObj, ok := evaluated.(*object.Error)
				if !ok {
					t.Errorf("object is not Error. got=%T(%+v)", evaluated, evaluated)
					return
				}

				if errObj.Message != expected {
					t.Errorf("wrong error message. expected=%q, got=%q", expected, errObj.Message)
				}
			case []int:
				array, ok := evaluated.(*object.Array)
				if !ok {
					t.Errorf("obj not Array. got=%T (%+v)", evaluated, evaluated)
					return
				}

				if len(array.Elements) != len(expected) {
					t.Errorf("wrong num of elements. want=%d, got=%d", len(expected), len(array.Elements))
					return
				}

				for i, expectedElem := range expected {
					testIntegerObject(t, array.Elements[i], int64(expectedElem))
				}

			}
		})
	}
}

// 配列リテラルのための評価器のテスト
func TestArray(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not Array. got=%T(%+v)", evaluated, evaluated)
	}

	if len(result.Elements) != 3 {
		t.Fatalf("array has wrong num of elements. got=%d", len(result.Elements))
	}

	testIntegerObject(t, result.Elements[0], 1)
	testIntegerObject(t, result.Elements[1], 4)
	testIntegerObject(t, result.Elements[2], 6)
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
	}{
		{
			"添字によるアクセス",
			"[1, 2, 3][0]",
			1,
		},
		{
			"添字によるアクセス",
			"[1, 2, 3][1]",
			2,
		},
		{
			"添字によるアクセス",
			"[1, 2, 3][2]",
			3,
		},
		{
			"添字に変数を使える",
			"let i = 0; [1][i];",
			1,
		},
		{
			"添字に式を使える",
			"[1, 2, 3][1 + 1];",
			3,
		},
		{
			"配列が変数でもOK",
			"let myArray = [1, 2, 3]; myArray[2];",
			3,
		},
		{
			"配列が変数でもOK",
			"let myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];",
			6,
		},
		{
			"添字演算式を変数の値として使える",
			"let myArray = [1, 2, 3]; let i = myArray[0]; myArray[i]",
			2,
		},
		{
			"範囲外の添字アクセスはNULL",
			"[1, 2, 3][3]",
			nil,
		},
		{
			"マイナスの添字アクセスはNULL",
			"[1, 2, 3][-1]",
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluated := testEval(tt.input)
			integer, ok := tt.expected.(int)
			if ok {
				testIntegerObject(t, evaluated, int64(integer))
			} else {
				testNullObject(t, evaluated)
			}
		})
	}
}

// *ast.HashLiteralに遭遇したときに何を返してほしいかを記述している
// また、Hashのキーとして任意の式が使えることも検証している
func TestHashLiterals(t *testing.T) {
	input := `let two = "two";
{
	"one": 10 - 9,
	two: 1 + 1,
	"thr" + "ee": 6 / 2,
	4: 4,
	true: 5,
	false: 6
}`
	evaluated := testEval(input)
	result, ok := evaluated.(*object.Hash)
	if !ok {
		t.Fatalf("Eval didn't return Hash. got=%T(%+v)", evaluated, evaluated)
	}

	expected := map[object.HashKey]int64{
		(&object.String{Value: "one"}).HashKey():   1,
		(&object.String{Value: "two"}).HashKey():   2,
		(&object.String{Value: "three"}).HashKey(): 3,
		(&object.Integer{Value: 4}).HashKey():      4,
		TRUE.HashKey():                             5,
		FALSE.HashKey():                            6,
	}

	if len(result.Pairs) != len(expected) {
		t.Fatalf("Hash has wrong num of pairs. got=%d", len(result.Pairs))
	}

	for expectedKey, expectedValue := range expected {
		pair, ok := result.Pairs[expectedKey]
		if !ok {
			t.Errorf("no pair for given key in Pairs")
		}

		testIntegerObject(t, pair.Value, expectedValue)
	}

}

// ハッシュリテラルでの添字演算式が正しい値を生成することを確認するテスト
func TestHashIndexExpressions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
	}{
		{
			"存在するハッシュキーでのアクセス",
			`{"foo": 5}["foo"]`,
			5,
		},
		{
			"存在しないハッシュキーでのアクセスはNULL",
			`{"foo": 5}["bar"]`,
			nil,
		},
		{
			"ハッシュキーでのアクセスに識別子が使える",
			`let key = "foo"; {"foo": 5}[key]`,
			5,
		},
		{
			"空のハッシュリテラルに何かしらのキーでアクセスしてもNULL",
			`{}["foo"]`,
			nil,
		},
		{
			"ハッシュキーに整数リテラルも使っても正しく評価できる",
			`{5: 5}[5]`,
			5,
		},
		{
			"ハッシュキーに真偽値リテラルも使っても正しく評価できる",
			`{true: 5}[true]`,
			5,
		},
		{
			"ハッシュキーに真偽値リテラルも使っても正しく評価できる",
			`{false: 5}[false]`,
			5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluated := testEval(tt.input)
			integer, ok := tt.expected.(int)
			if ok {
				testIntegerObject(t, evaluated, int64(integer))
			} else {
				testNullObject(t, evaluated)
			}
		})
	}
}
