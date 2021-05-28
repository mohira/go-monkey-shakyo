package main

import (
	"fmt"
	"go-monkey-shakyo/monkey/object"
)

// p.219 object.Hashの実装でうまくいかない問題を示すサンプル
func main() {
	name1 := &object.String{Value: "name"}
	monkey := &object.String{Value: "Monkey"}

	pairs := map[object.Object]object.Object{}
	pairs[name1] = monkey

	fmt.Printf("pairs[name1]=%+v\n", pairs[name1]) // pairs[name1]=&{Value:Monkey}

	// name1とname2は 同じ .Value を持っているが...
	name2 := &object.String{Value: "name"}
	fmt.Printf("pairs[name2]=%+v\n", pairs[name2]) // pairs[name2]=<nil>

	fmt.Printf("(name1 == name2)=%t\n", name1 == name2) // (name1 == name2)=false
}
