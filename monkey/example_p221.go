package main

import (
	"fmt"
	"go-monkey-shakyo/monkey/object"
)

// p.221 ハッシュ値を生成できるようになったので、object.Hashの実装でうまくいく
func main() {
	name1 := &object.String{Value: "name"}
	monkey := &object.String{Value: "Monkey"}

	pairs := map[object.HashKey]object.Object{}
	pairs[name1.HashKey()] = monkey

	fmt.Printf("pairs[name1.HashKey()]=%+v\n", pairs[name1.HashKey()]) // pairs[name1.HashKey()]=&{Value:Monkey}

	// name1とname2は 同じ .Value を持っているのでハッシュ値も一致する
	name2 := &object.String{Value: "name"}
	fmt.Printf("pairs[name2.HashKey()]=%+v\n", pairs[name2.HashKey()]) // pairs[name2.HashKey()]=&{Value:Monkey}

	// name1とname2はポインタなので、直接比較するとfalse
	fmt.Printf("(name1 == name2)=%t\n", name1 == name2) // (name1 == name2)=false

	// name1とname2をハッシュ値で比較するとtrue
	fmt.Printf("(name1.HashKey() == name2.HashKey())=%t\n", name1.HashKey() == name2.HashKey()) // (name1 == name2)=true
}
