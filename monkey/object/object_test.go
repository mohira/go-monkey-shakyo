package object

import (
	"testing"
)

// 簡単に比較可能で、object.Hashのハッシュキーとして使えるような、オブジェクトのハッシュ値を生成する方法
// ある*object.Stringに対するハッシュキーは、
// 別の*object.Stringインスタンスのハッシュキーとも比較可能であり、
// 同じ.Valueを持つのであれば、その結果は等しくなければならない。
func TestStringHashKey(t *testing.T) {
	hello1 := &String{Value: "Hello World"}
	hello2 := &String{Value: "Hello World"}

	diff1 := &String{Value: "My name is johnny"}
	diff2 := &String{Value: "My name is johnny"}

	// 同じ.Valueを持つなら、ハッシュ値も等しくなければならない
	if hello1.HashKey() != hello2.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}

	// 同じ.Valueを持つなら、ハッシュ値も等しくなければならない
	if diff1.HashKey() != diff2.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}

	// 異なる.Valueを持つなら、ハッシュ値も異ならないといけない
	if hello1.HashKey() == diff1.HashKey() {
		t.Errorf("strings with different content have save hash keys")
	}
}
