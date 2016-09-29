package xlate

import (
	"fmt"
	"testing"
)

func capturePanics(fn func()) (e error) {
	defer func() {
		if rcvr := recover(); rcvr != nil {
			e = rcvr.(error)
		}
	}()
	fn()
	return
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
