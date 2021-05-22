package evaluator

import "go-monkey-shakyo/monkey/object"

var builtins = map[string]*object.Builtin{
	"len": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			return NULL
		}},
}
