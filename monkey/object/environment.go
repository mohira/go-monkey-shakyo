package object

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s}
}

type Environment struct {
	store map[string]Object
	outer *Environment
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]

	// 現在の「環境」には存在しないが、「外の環境」にはあるかもしれないので探しに行く
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}

	return obj, ok
}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

// "enclosed" の意味を理解するとイメージできる
// http://ejbridge2.blog.fc2.com/blog-entry-146.html
// https://res.cloudinary.com/dyd911kmh/image/upload/f_auto,q_auto:best/v1588956604/Scope_fbrzcw.png
func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}
