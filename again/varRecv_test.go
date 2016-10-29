package again

import (
	"fmt"
	"testing"
)

func TestWow(t *testing.T) {
	var v int
	vr := NewVarReceiver(&v)
	tok := Token(4)
	vr.Step(&tok)

	assert(t, "simple literal test", 4, v)
}

func stringyEquality(x, y interface{}) bool {
	return fmt.Sprintf("%#v", x) == fmt.Sprintf("%#v", y)
}

func assert(t *testing.T, title string, expect, actual interface{}) {
	if !stringyEquality(expect, actual) {
		t.Errorf("test %q FAILED:\n\texpected  %#v\n\tactual    %#v",
			title, expect, actual)
	}
}
